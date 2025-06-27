// Package eventloop implements a simple event loop store that uses a single goroutine to handle all commands.
package eventloop

import (
	"context"
	"time"

	"github.com/PlayerNeo42/gvalkey/resp"
	"github.com/PlayerNeo42/gvalkey/store"
)

var _ store.Store = (*EventloopStore)(nil)

type EventloopStore struct {
	m          map[string]any
	expiration map[string]time.Time

	cmdCh chan cmd
}

func NewEventloopStore() *EventloopStore {
	s := &EventloopStore{
		m:          make(map[string]any),
		expiration: make(map[string]time.Time),
		cmdCh:      make(chan cmd, 1),
	}

	go s.Run(context.Background())
	return s
}

func (s *EventloopStore) Get(key string) (any, bool) {
	result := executeCommand[operationResult](s, CmdGet, key)
	return result.Value, result.OK
}

func (s *EventloopStore) Del(key string) bool {
	return executeCommand[bool](s, CmdDel, key)
}

func (s *EventloopStore) Set(args resp.SetArgs) (any, bool) {
	result := executeCommand[operationResult](s, CmdSet, args)
	return result.Value, result.OK
}

// Close closes the event loop and stops the cleanup goroutine.
func (s *EventloopStore) Close() {
	close(s.cmdCh)
}

// Run starts the event loop and cleans up expired keys every second.
func (s *EventloopStore) Run(ctx context.Context) {
	// clean up expired keys every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.expireKeys()
		case cmd := <-s.cmdCh:
			// handle commands, all data operations are executed in this single goroutine
			s.handleCommand(cmd)
		}
	}
}

func (s *EventloopStore) handleCommand(cmd cmd) {
	switch cmd.typ {
	case CmdGet:
		if respCh, ok := cmd.resp.(chan operationResult); ok {
			if key, ok := cmd.payload.(string); ok {
				respCh <- s.handleGet(key)
			}
		}

	case CmdSet:
		if respCh, ok := cmd.resp.(chan operationResult); ok {
			if args, ok := cmd.payload.(resp.SetArgs); ok {
				respCh <- s.handleSet(args)
			}
		}

	case CmdDel:
		if respCh, ok := cmd.resp.(chan bool); ok {
			if key, ok := cmd.payload.(string); ok {
				respCh <- s.handleDel(key)
			}
		}
	}
}

func (s *EventloopStore) handleGet(key string) operationResult {
	if s.isExpired(key) {
		delete(s.m, key)
		delete(s.expiration, key)
		return operationResult{Value: nil, OK: false}
	}

	value, exists := s.m[key]
	return operationResult{Value: value, OK: exists}
}

func (s *EventloopStore) handleSet(args resp.SetArgs) operationResult {
	key := string(args.Key.MarshalBinary())

	oldValue, exists := s.m[key]

	// treat expired keys as not existing for the purpose of nx/xx logic.
	if exists && s.isExpired(key) {
		exists = false
	}

	shouldNotSet := (args.NX && exists) || (args.XX && !exists)
	if shouldNotSet {
		if args.Get && exists {
			return operationResult{Value: oldValue, OK: true}
		}
		return operationResult{Value: nil, OK: false}
	}

	// set value
	s.m[key] = args.Value

	// set expiration time
	if !args.ExpireAt.IsZero() {
		s.expiration[key] = args.ExpireAt
	} else {
		// remove expiration time if previously set
		delete(s.expiration, key)
	}

	// return result
	if args.Get && exists {
		return operationResult{Value: oldValue, OK: true}
	}
	// if key didn't exist or was expired, there's no old value to return, but the set was successful.
	return operationResult{Value: nil, OK: true}
}

func (s *EventloopStore) handleDel(key string) bool {
	_, exists := s.m[key]
	if exists {
		delete(s.m, key)
		delete(s.expiration, key)
	}

	return exists
}

func (s *EventloopStore) expireKeys() {
	now := time.Now()
	for key, expireAt := range s.expiration {
		if now.After(expireAt) {
			delete(s.m, key)
			delete(s.expiration, key)
		}
	}
}

func (s *EventloopStore) isExpired(key string) bool {
	expireTime, exists := s.expiration[key]
	return exists && time.Now().After(expireTime)
}

// executeCommand executes a command and waits for the result
func executeCommand[T any](s *EventloopStore, cmdType int, payload any) T {
	respCh := make(chan T)
	s.cmdCh <- cmd{typ: cmdType, payload: payload, resp: respCh}
	return <-respCh
}
