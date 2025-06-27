package store_test

import (
	"sync"
	"testing"
	"time"

	"github.com/PlayerNeo42/gvalkey/resp"
	"github.com/PlayerNeo42/gvalkey/store"
	"github.com/PlayerNeo42/gvalkey/store/eventloop"
	"github.com/PlayerNeo42/gvalkey/store/naive"
	"github.com/stretchr/testify/suite"
)

// TestNaiveStore tests the naive store implementation
func TestNaiveStore(t *testing.T) {
	naiveStoreFactory := func() store.Store {
		return naive.NewNaiveStore()
	}

	var naiveStore *naive.NaiveStore
	cleanup := func() {
		if naiveStore != nil {
			naiveStore.Close()
		}
	}

	suite.Run(t, &StoreTestSuite{
		storeFactory: naiveStoreFactory,
		cleanup:      cleanup,
	})
}

// TestEventloopStore tests the eventloop store implementation
func TestEventloopStore(t *testing.T) {
	eventloopStoreFactory := func() store.Store {
		return eventloop.NewEventloopStore()
	}

	cleanup := func() {}

	suite.Run(t, &StoreTestSuite{
		storeFactory: eventloopStoreFactory,
		cleanup:      cleanup,
	})
}

// MockBinaryMarshaler implements BinaryMarshaler interface for testing
type MockBinaryMarshaler struct {
	data string
}

func (m MockBinaryMarshaler) MarshalBinary() []byte {
	return []byte(m.data)
}

// StoreTestSuite defines a common test suite that can test any type that implements the Store interface
type StoreTestSuite struct {
	suite.Suite
	storeFactory func() store.Store
	cleanup      func()
	store        store.Store
}

// SetupTest initializes the store before each test
func (s *StoreTestSuite) SetupTest() {
	s.store = s.storeFactory()

	// For EventloopStore, wait for event loop to start
	s.Require().Eventually(func() bool {
		// Try a simple Set operation to test if event loop has started
		testArgs := resp.SetArgs{
			Key:   MockBinaryMarshaler{data: "__startup_test__"},
			Value: "test",
		}
		_, ok := s.store.Set(testArgs)
		if ok {
			// Clean up test key
			s.store.Del("__startup_test__")
			return true
		}
		return false
	}, 1*time.Second, 10*time.Millisecond, "store should start within specified time")
}

// TearDownTest cleans up after each test
func (s *StoreTestSuite) TearDownTest() {
	// Clean up the store itself (if needed)
	if s.cleanup != nil {
		s.cleanup()
	}
}

// TestBasicOperations tests basic Set/Get/Del operations
func (s *StoreTestSuite) TestBasicOperations() {
	// Test Set and Get
	setArgs := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "testkey"},
		Value: "testvalue",
	}
	_, ok := s.store.Set(setArgs)
	s.Require().True(ok, "Set operation should succeed")

	value, exists := s.store.Get("testkey")
	s.Require().True(exists, "Key should exist")
	s.Require().Equal("testvalue", value, "Value should match the set value")

	// Test Del
	deleted := s.store.Del("testkey")
	s.Require().True(deleted, "Delete operation should return true")

	// Verify key has been deleted
	_, exists = s.store.Get("testkey")
	s.Require().False(exists, "Key should not exist after deletion")

	// Test deleting non-existent key
	deleted = s.store.Del("nonexistent")
	s.Require().False(deleted, "Deleting non-existent key should return false")
}

// TestExpiration tests key expiration functionality
func (s *StoreTestSuite) TestExpiration() {
	// Set a key that expires after 1 second
	setArgs := resp.SetArgs{
		Key:      MockBinaryMarshaler{data: "expirekey"},
		Value:    "expirevalue",
		ExpireAt: time.Now().Add(1 * time.Second),
	}
	_, ok := s.store.Set(setArgs)
	s.Require().True(ok, "Setting expiring key should succeed")

	// Should exist immediately
	value, exists := s.store.Get("expirekey")
	s.Require().True(exists, "Key should exist immediately after setting")
	s.Require().Equal("expirevalue", value, "Value should match the set value")

	// Wait for expiration
	s.Require().Eventually(func() bool {
		_, exists := s.store.Get("expirekey")
		return !exists
	}, 3*time.Second, 100*time.Millisecond, "Key should expire after specified time")
}

// TestSetNX tests the NX flag (only set if key does not exist)
func (s *StoreTestSuite) TestSetNX() {
	// First set a key
	setArgs := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "nxkey"},
		Value: "original",
	}
	_, ok := s.store.Set(setArgs)
	s.Require().True(ok, "Initial set should succeed")

	// Try to reset with NX, should fail
	setArgsNX := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "nxkey"},
		Value: "new",
		NX:    true,
	}
	_, ok = s.store.Set(setArgsNX)
	s.Require().False(ok, "NX set should fail when key exists")

	// Verify value hasn't changed
	value, exists := s.store.Get("nxkey")
	s.Require().True(exists, "Key should still exist")
	s.Require().Equal("original", value, "Value should remain unchanged when NX fails")

	// Using NX on non-existent key should succeed
	setArgsNXNew := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "newkey"},
		Value: "newvalue",
		NX:    true,
	}
	_, ok = s.store.Set(setArgsNXNew)
	s.Require().True(ok, "NX set should succeed when key does not exist")

	value, exists = s.store.Get("newkey")
	s.Require().True(exists, "New key should exist")
	s.Require().Equal("newvalue", value, "New key value should be correct")
}

// TestSetXX tests the XX flag (only set if key exists)
func (s *StoreTestSuite) TestSetXX() {
	// Using XX on non-existent key should fail
	setArgsXX := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "xxkey"},
		Value: "value",
		XX:    true,
	}
	_, ok := s.store.Set(setArgsXX)
	s.Require().False(ok, "XX set should fail when key does not exist")

	_, exists := s.store.Get("xxkey")
	s.Require().False(exists, "Key should not be created")

	// First set the key
	setArgs := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "xxkey"},
		Value: "original",
	}
	_, ok = s.store.Set(setArgs)
	s.Require().True(ok, "Initial set should succeed")

	// Now updating with XX should succeed
	setArgsXXUpdate := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "xxkey"},
		Value: "updated",
		XX:    true,
	}
	_, ok = s.store.Set(setArgsXXUpdate)
	s.Require().True(ok, "XX set should succeed when key exists")

	value, exists := s.store.Get("xxkey")
	s.Require().True(exists, "Key should exist")
	s.Require().Equal("updated", value, "Value should be updated")
}

// TestSetGET tests the GET flag (returns old value)
func (s *StoreTestSuite) TestSetGET() {
	// Using GET on non-existent key
	setArgsGET := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "getkey"},
		Value: "newvalue",
		Get:   true,
	}
	oldValue, ok := s.store.Set(setArgsGET)
	s.Require().True(ok, "Set with GET should succeed")
	s.Require().Nil(oldValue, "Old value for new key should be nil")

	// Using GET on existing key
	setArgsGETUpdate := resp.SetArgs{
		Key:   MockBinaryMarshaler{data: "getkey"},
		Value: "updatedvalue",
		Get:   true,
	}
	oldValue, ok = s.store.Set(setArgsGETUpdate)
	s.Require().True(ok, "Update with GET should succeed")
	s.Require().Equal("newvalue", oldValue, "Should return old value")

	// Verify new value is set
	value, exists := s.store.Get("getkey")
	s.Require().True(exists, "Key should exist")
	s.Require().Equal("updatedvalue", value, "Should have new value")
}

// TestConcurrency tests concurrent operations
func (s *StoreTestSuite) TestConcurrency() {
	var wg sync.WaitGroup
	numGoroutines := 50

	// Concurrent setting of different keys
	for i := range numGoroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := MockBinaryMarshaler{data: "concurrentkey" + string(rune(i))}
			setArgs := resp.SetArgs{Key: key, Value: i}
			_, ok := s.store.Set(setArgs)
			s.True(ok)
		}(i)
	}
	wg.Wait()

	// Concurrent reading
	for i := range numGoroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "concurrentkey" + string(rune(i))
			value, exists := s.store.Get(key)
			s.True(exists)
			s.Equal(i, value)
		}(i)
	}
	wg.Wait()

	// Concurrent deletion
	for i := range numGoroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "concurrentkey" + string(rune(i))
			deleted := s.store.Del(key)
			s.True(deleted)
		}(i)
	}
	wg.Wait()

	// Verify all keys have been deleted
	for i := range numGoroutines {
		key := "concurrentkey" + string(rune(i))
		_, exists := s.store.Get(key)
		s.Require().False(exists, "Key should be deleted")
	}
}
