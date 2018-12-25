package memcache

import "time"

type cacheNode struct {
	createTime time.Time
	ttl        time.Duration
	data       CacheData
}

type CacheData interface {
	GetNodeID() string
	GetNodeBytes() []byte
}
