// Package naive implements a thread-safe in-memory key-value store using sync.Map.
package naive

import (
	"sync"
	"time"

	"github.com/PlayerNeo42/gvalkey/resp"
	"github.com/PlayerNeo42/gvalkey/store"
)

var _ store.Store = (*NaiveStore)(nil)

type naiveStoreItem struct {
	value      any
	expiration time.Time // expiration timestamp, 0 means never expire
}

func (item *naiveStoreItem) isExpired() bool {
	if item.expiration.IsZero() {
		return false
	}
	return time.Now().After(item.expiration)
}

// NaiveStore is a thread-safe in-memory key-value store implementation using Go's sync.Map.
type NaiveStore struct {
	store       sync.Map
	stopCleanup chan struct{} // channel for stopping the cleanup goroutine
}

func NewNaiveStore() *NaiveStore {
	ms := &NaiveStore{
		stopCleanup: make(chan struct{}),
	}

	go ms.cleanupExpiredKeys()

	return ms
}

func (s *NaiveStore) Set(args resp.SetArgs) (any, bool) {
	key := args.Key.String()

	// load the existing item first to check its status.
	existing, exists := s.store.Load(key)

	var oldItem *naiveStoreItem
	if exists {
		var ok bool
		oldItem, ok = existing.(*naiveStoreItem)
		if !ok {
			// this should not happen in normal operation, but as a safeguard, treat it as not exists.
			exists = false
		} else if oldItem.isExpired() {
			// treat expired keys as not existing for the purpose of nx/xx logic.
			exists = false
		}
	}

	// handle conditional set flags (NX, XX).
	// we should not set the key if:
	// 1. the key exists and NX is true, or
	// 2. the key doesn't exist and XX is true.
	shouldNotSet := (args.NX && exists) || (args.XX && !exists)
	if shouldNotSet {
		if args.Get && exists {
			// for NX, if key exists, we return the old value.
			return oldItem.value, true
		}
		// for XX, if key doesn't exist, there is no old value.
		// for NX, if key exists, we don't set and return success=false.
		return nil, false
	}

	newItem := &naiveStoreItem{
		value:      args.Value,
		expiration: args.ExpireAt,
	}
	s.store.Store(key, newItem)

	if args.Get && exists {
		// 'exists' is true only if the key existed and was not expired.
		return oldItem.value, true
	}

	// if key didn't exist or was expired, there's no old value to return, but the set was successful.
	return nil, true
}

func (s *NaiveStore) Get(key string) (any, bool) {
	value, exists := s.store.Load(key)
	if !exists {
		return nil, false
	}

	item, ok := value.(*naiveStoreItem)
	if !ok {
		return nil, false
	}

	// check if expired
	if item.isExpired() {
		s.store.Delete(key)
		return nil, false
	}

	return item.value, true
}

func (s *NaiveStore) Del(key string) bool {
	existing, existed := s.store.LoadAndDelete(key)
	if !existed {
		return false
	}

	item, ok := existing.(*naiveStoreItem)
	if !ok {
		// the key existed but was not a naiveStoreItem.
		// this is unexpected, but it was deleted, so we return true.
		return true
	}

	// return false if the key was expired (logically didn't exist), true otherwise.
	return !item.isExpired()
}

func (s *NaiveStore) cleanupExpiredKeys() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.store.Range(func(key, value any) bool {
				if item, ok := value.(*naiveStoreItem); ok && item.isExpired() {
					s.store.Delete(key)
				}
				return true
			})
		case <-s.stopCleanup:
			return
		}
	}
}

// Close stops the cleanup goroutine
func (s *NaiveStore) Close() {
	close(s.stopCleanup)
}
