// Package store provides multiple thread-safe in-memory key-value store implementations.
package store

type Store interface {
	Get(key string) (any, bool)
	Set(key string, value any, ex, px *int64, nx, xx, get bool) (any, bool)
	Del(key string) bool
}
