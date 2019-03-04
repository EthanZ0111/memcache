package memcache

import (
	"errors"
	"sync"
	"time"
)

type MemCacheMap struct {
	cachePool map[string]*cacheNode
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
	c.cachePool = make(map[string]*cacheNode, cacheSize)
	c.clearRate = time.Millisecond * time.Duration(clearRate)
	go c.clearLoop()
	return c
}

// AsyncAdd : will download or build data bytes call FillData() with async mode
func (c *MemCacheMap) AsyncAdd(data CacheData) error {
	id := data.GetID()
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.cachePool[id]; ok {
		return ExistError
	}
	c.cachePool[id] = &cacheNode{
		data:       data,
		createTime: time.Now(),
		// ttl:        c.ttl,
		// dataError:  nil,
	}
	go func() {
		data.FillData()
		// c.lock.Lock()
		// c.cache[id].dataError = err
		// c.lock.Unlock()
	}()
	return nil
}

// SyncAdd : will download or build data bytes call FillData() with sync mode,
// so user must make sure that FillData() will not take to long time.
func (c *MemCacheMap) SyncAdd(data CacheData) error {
	id := data.GetID()
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.cachePool[id]; ok {
		return ExistError
	}
	c.cachePool[id] = &cacheNode{
		data:       data,
		createTime: time.Now(),
		// ttl:        c.ttl,
		// dataError:  nil,
	}
	data.FillData()
	// c.cache[id].dataError = err
	return nil
}

func (c *MemCacheMap) Delete(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.cachePool[id]; ok {
		delete(c.cachePool, id)
	}
}

func (c *MemCacheMap) Get(id string) (CacheData, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if cd, ok := c.cachePool[id]; ok {
		if cd.data.GetError() != nil {
			delete(c.cachePool, id)
			return nil, cd.data.GetError()
		}
		return cd.data, nil
	}
	return nil, NotExistError
}

func (c *MemCacheMap) clearLoop() {
	var now time.Time
	for {
		now = time.Now()
		c.lock.Lock()
		for k, v := range c.cachePool {
			if v.createTime.Add(c.ttl).Sub(now) <= 0 {
				delete(c.cachePool, k)
			}
		}
		c.lock.Unlock()
		time.Sleep(c.clearRate)
	}
}
