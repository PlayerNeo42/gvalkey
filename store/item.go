package store

import "time"

type storeItem struct {
	value      any
	expiration int64 // Expiration timestamp, 0 means never expire
}

func (item *storeItem) isExpired() bool {
	if item.expiration == 0 {
		return false
	}
	return time.Now().Unix() > item.expiration
}
