package memcache

import "time"

type cacheNode struct {
	createTime time.Time
	ttl        time.Duration
	data       CacheData
	dataError  error
}

type CacheData interface {
	GetID() string
	GetBytes() ([]byte, error)
	FillData() error
}
