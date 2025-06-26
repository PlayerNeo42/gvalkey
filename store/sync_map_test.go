package store

import (
	"sync"
	"testing"
	"time"

	"github.com/PlayerNeo42/gvalkey/resp"
	"github.com/stretchr/testify/suite"
)

type NaiveStoreSuite struct {
	suite.Suite
	s *NaiveStore
}

func (s *NaiveStoreSuite) SetupTest() {
	s.s = NewNaiveStore()
}

func (s *NaiveStoreSuite) TearDownTest() {
	s.s.Close()
}

func TestNaiveStoreSuite(t *testing.T) {
	suite.Run(t, new(NaiveStoreSuite))
}

func (s *NaiveStoreSuite) TestSet_Get_Del() {
	// Test basic Set and Get
	_, ok := s.s.Set(resp.SetArgs{Key: resp.BulkString("key1"), Value: "value1"})
	s.Require().True(ok, "Set should succeed")

	val, ok := s.s.Get("key1")
	s.Require().True(ok, "Get should find the key")
	s.Require().Equal("value1", val, "Get should return the correct value")

	// Test Del
	deleted := s.s.Del("key1")
	s.Require().True(deleted, "Del should succeed for existing key")

	_, ok = s.s.Get("key1")
	s.Require().False(ok, "Get should not find the key after Del")

	// Test Del on non-existent key
	deleted = s.s.Del("non-existent")
	s.Require().False(deleted, "Del should fail for non-existent key")
}

func (s *NaiveStoreSuite) TestExpiration() {
	// Test EX
	expire := time.Now().UnixMilli() + 1000 // 1 second
	_, ok := s.s.Set(resp.SetArgs{Key: resp.BulkString("key_ex"), Value: "value_ex", Expire: expire})
	s.Require().True(ok)
	s.Require().Eventually(
		func() bool {
			_, ok = s.s.Get("key_ex")
			return !ok
		},
		2*time.Second,
		100*time.Millisecond,
		"EX expiration failed: key should have expired",
	)

	// Test PX
	expire = time.Now().UnixMilli() + 100 // 100 milliseconds
	_, ok = s.s.Set(resp.SetArgs{Key: resp.BulkString("key_px"), Value: "value_px", Expire: expire})
	s.Require().True(ok)
	s.Require().Eventually(
		func() bool {
			_, ok = s.s.Get("key_px")
			return !ok
		},
		200*time.Millisecond,
		10*time.Millisecond,
		"PX expiration failed: key should have expired",
	)
}

func (s *NaiveStoreSuite) TestSet_NX() {
	// Set with NX when key doesn't exist
	_, ok := s.s.Set(resp.SetArgs{Key: resp.BulkString("key_nx"), Value: "value_nx", NX: true})
	s.Require().True(ok, "Set with NX should succeed when key doesn't exist")

	val, ok := s.s.Get("key_nx")
	s.Require().True(ok)
	s.Require().Equal("value_nx", val)

	// Set with NX when key exists
	_, ok = s.s.Set(resp.SetArgs{Key: resp.BulkString("key_nx"), Value: "new_value", NX: true})
	s.Require().False(ok, "Set with NX should fail when key exists")

	val, _ = s.s.Get("key_nx")
	s.Require().Equal("value_nx", val, "Value should not have changed")
}

func (s *NaiveStoreSuite) TestSet_XX() {
	// Set with XX when key doesn't exist
	_, ok := s.s.Set(resp.SetArgs{Key: resp.BulkString("key_xx"), Value: "value_xx", XX: true})
	s.Require().False(ok, "Set with XX should fail when key doesn't exist")

	_, ok = s.s.Get("key_xx")
	s.Require().False(ok, "Key should not exist")

	// Set with XX when key exists
	_, ok = s.s.Set(resp.SetArgs{Key: resp.BulkString("key_xx"), Value: "value_xx"})
	s.Require().True(ok)

	_, ok = s.s.Set(resp.SetArgs{Key: resp.BulkString("key_xx"), Value: "new_value", XX: true})
	s.Require().True(ok, "Set with XX should succeed when key exists")

	val, _ := s.s.Get("key_xx")
	s.Require().Equal("new_value", val, "Value should have been updated")
}

func (s *NaiveStoreSuite) TestSet_GET() {
	// Set with GET when key doesn't exist
	oldVal, ok := s.s.Set(resp.SetArgs{Key: resp.BulkString("key_get"), Value: "value_get", Get: true})
	s.Require().True(ok, "Set with GET should succeed")
	s.Require().Nil(oldVal, "Old value should be nil for a new key")

	// Set with GET when key exists
	oldVal, ok = s.s.Set(resp.SetArgs{Key: resp.BulkString("key_get"), Value: "new_value", Get: true})
	s.Require().True(ok, "Set with GET should succeed for an existing key")
	s.Require().Equal("value_get", oldVal, "Should return the old value")
}

func (s *NaiveStoreSuite) TestConcurrency() {
	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent Set
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := resp.BulkString("key")
			value := i
			_, ok := s.s.Set(resp.SetArgs{Key: key, Value: value})
			s.True(ok)
		}(i)
	}
	wg.Wait()

	// Concurrent Get
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.s.Get("key")
		}()
	}
	wg.Wait()

	// Concurrent Del
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.s.Del("key")
		}()
	}
	wg.Wait()
}
