package memcache

import (
	"errors"
	"sync"
	"time"
)

type MemCacheMap struct {
	cache     map[string]cacheNode
	ttl       time.Duration
	lock      sync.RWMutex
	clearRate time.Duration
}

var (
	ExistError    = errors.New("Cache is already exist")
	NotExistError = errors.New("Cache is not exist")
)

const (
	// DefaultCacheSize : default cache size is 10
	DefaultCacheSize = 10
	// DefaultClearRate : default clear rate is 60 second
	DefaultClearRate = 60 * 1000
	// DefalutTTL : default TTL is 10 minutes
	DefalutTTL = 60 * 1000 * 10
)

func NewMemCacheMap(cacheSize int, clearRate int, ttl int) *MemCacheMap {
	c := new(MemCacheMap)
	if cacheSize <= 0 {
		cacheSize = DefaultCacheSize
	}
	if clearRate <= 0 {
		clearRate = DefaultClearRate
	}
	if ttl <= 0 {
		ttl = DefalutTTL
	}
	c.ttl = time.Millisecond * time.Duration(ttl)
	c.cache = make(map[string]cacheNode, cacheSize)
	c.clearRate = time.Millisecond * time.Duration(clearRate)

	return c
}

func (c *MemCacheMap) Add(data CacheData) error {
	id := data.GetNodeID()
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.cache[id]; ok {
		return ExistError
	}
	c.cache[id] = cacheNode{
		data:       data,
		createTime: time.Now(),
		ttl:        c.ttl,
	}
	return nil
}

func (c *MemCacheMap) Get(id string) (CacheData, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c, ok := c.cache[id]; ok {
		return c.data, nil
	}
	return nil, NotExistError
}

func (c *MemCacheMap) clearLoop() {
	// var now time.Time
	for {
		// now = time.Now()
		c.lock.Lock()
		// for k, v := range c.cache {
		// 	if v.createTime.Add(v.ttl).Sub(now) <= 0 {

		// 	}
		// }
		c.lock.Unlock()
		time.Sleep(time.Millisecond * c.clearRate)
	}
}
