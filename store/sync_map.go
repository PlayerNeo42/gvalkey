package store

import (
	"sync"
	"time"

	"github.com/PlayerNeo42/gvalkey/resp"
)

type syncMapItem struct {
	value      any
	expiration int64 // Expiration timestamp, 0 means never expire
}

func (item *syncMapItem) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().UnixMilli() > item.expiration
}

// SyncMap is a thread-safe in-memory key-value store implementation using Go's sync.Map.
type SyncMap struct {
	store       sync.Map
	stopCleanup chan struct{} // channel for stopping the cleanup goroutine
}

func NewSyncMap() *SyncMap {
	ms := &SyncMap{
		stopCleanup: make(chan struct{}),
	}

	go ms.cleanupExpiredKeys()

	return ms
}

func (s *SyncMap) Set(args resp.SetArgs) (any, bool) {
	key := args.Key.MarshalBinary()

	// load the existing item first to check its status.
	existing, loaded := s.store.Load(string(key))

	var oldItem *syncMapItem
	if loaded {
		var ok bool
		oldItem, ok = existing.(*syncMapItem)
		if !ok {
			// this should not happen in normal operation, but as a safeguard, treat it as not loaded.
			loaded = false
		} else if oldItem.isExpired() {
			// treat expired keys as not existing for the purpose of nx/xx logic.
			loaded = false
		}
	}

	// handle conditional set flags (NX, XX).
	// we should not set the key if:
	// 1. the key exists and NX is true, or
	// 2. the key doesn't exist and XX is true.
	shouldNotSet := (args.NX && loaded) || (args.XX && !loaded)
	if shouldNotSet {
		if args.Get && loaded {
			// for NX, if key exists, we return the old value.
			return oldItem.value, true
		}
		// for XX, if key doesn't exist, there is no old value.
		// for NX, if key exists, we don't set and return success=false.
		return nil, false
	}

	// If we proceed, it means we will perform a set operation.
	newItem := &syncMapItem{
		value:      args.Value,
		expiration: args.Expire,
	}
	s.store.Store(string(key), newItem)

	if args.Get {
		if loaded { // 'loaded' is true only if the key existed and was not expired.
			return oldItem.value, true
		}
		// If key didn't exist or was expired, there's no old value to return, but the set was successful.
		return nil, true
	}

	return nil, true
}

func (s *SyncMap) Get(key string) (any, bool) {
	value, exists := s.store.Load(key)
	if !exists {
		return nil, false
	}

	item, ok := value.(*syncMapItem)
	if !ok {
		return nil, false
	}

	// Check if expired
	if item.isExpired() {
		s.store.Delete(key)
		return nil, false
	}

	return item.value, true
}

func (s *SyncMap) Del(key string) bool {
	existing, existed := s.store.LoadAndDelete(key)
	if !existed {
		return false
	}

	item, ok := existing.(*syncMapItem)
	if !ok {
		// The key existed but was not a syncMapItem.
		// This is unexpected, but it was deleted, so we return true.
		return true
	}

	// Return false if the key was expired (logically didn't exist), true otherwise.
	return !item.isExpired()
}

func (s *SyncMap) cleanupExpiredKeys() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.store.Range(func(key, value any) bool {
				if item, ok := value.(*syncMapItem); ok && item.isExpired() {
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
func (s *SyncMap) Close() {
	close(s.stopCleanup)
}
