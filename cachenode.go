package memcache

import "time"

type cacheNode struct {
	createTime time.Time
	data       CacheData
	// ttl        time.Duration
	// dataError  error
}

type CacheData interface {
	GetID() string
	GetBytes() ([]byte, error)
	FillData() error
	GetError() error
}
