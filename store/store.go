// Package store provides multiple thread-safe in-memory key-value store implementations.
package store

import "github.com/PlayerNeo42/gvalkey/resp"

type Store interface {
	Get(key string) (any, bool)
	Set(args resp.SetArgs) (any, bool)
	Del(key string) bool
}
