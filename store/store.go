package store

import (
	"sync"
	"time"
)

type Store struct {
	store       sync.Map
	stopCleanup chan struct{} // channel for stopping the cleanup goroutine
}

func NewStore() *Store {
	ms := &Store{
		stopCleanup: make(chan struct{}),
	}

	go ms.cleanupExpiredKeys()

	return ms
}

// Background task to clean up expired keys
func (s *Store) cleanupExpiredKeys() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.store.Range(func(key, value any) bool {
				if item, ok := value.(*storeItem); ok && item.isExpired() {
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
func (s *Store) Close() {
	close(s.stopCleanup)
}

func (s *Store) Set(key string, value any, ex, px int64, nx, xx, get bool) (any, bool) {
	var expiration int64
	if ex > 0 {
		expiration = time.Now().UnixMilli() + ex*1000
	} else if px > 0 {
		expiration = time.Now().UnixMilli() + px
	}

	newItem := &storeItem{
		value:      value,
		expiration: expiration,
	}

	// If we need to return the old value, get it first
	var oldValue any
	var hasOldValue bool
	if get {
		if existing, exists := s.store.Load(key); exists {
			if item, ok := existing.(*storeItem); ok && !item.isExpired() {
				oldValue = item.value
				hasOldValue = true
			}
		}
	}

	// Handle NX (Not eXists) option: only set if key doesn't exist
	if nx {
		if existing, exists := s.store.Load(key); exists {
			if item, ok := existing.(*storeItem); ok && !item.isExpired() {
				// Key exists and is not expired, cannot set
				if get {
					return oldValue, hasOldValue
				}
				return nil, false
			}
		}
		// Key doesn't exist or is expired, can set
		s.store.Store(key, newItem)
		if get {
			return oldValue, hasOldValue
		}
		return nil, true
	}

	// Handle XX (eXists) option: only set if key already exists
	if xx {
		if existing, exists := s.store.Load(key); exists {
			if item, ok := existing.(*storeItem); ok && !item.isExpired() {
				// Key exists and is not expired, can set
				s.store.Store(key, newItem)
				if get {
					return oldValue, hasOldValue
				}
				return nil, true
			}
		}
		// Key doesn't exist or is expired, cannot set
		if get {
			return oldValue, hasOldValue
		}
		return nil, false
	}

	// Normal set: unconditional set
	s.store.Store(key, newItem)
	if get {
		return oldValue, hasOldValue
	}
	return nil, true
}

func (s *Store) Get(key string) (any, bool) {
	value, exists := s.store.Load(key)
	if !exists {
		return nil, false
	}

	item, ok := value.(*storeItem)
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

func (s *Store) Del(key string) bool {
	// Check if key exists and is not expired
	if existing, exists := s.store.Load(key); exists {
		if item, ok := existing.(*storeItem); ok {
			if item.isExpired() {
				// Expired key is deleted directly, but return false (logically doesn't exist)
				s.store.Delete(key)
				return false
			}
		}
	}

	_, existed := s.store.LoadAndDelete(key)
	return existed
}
