package memcache

import (
	"bytes"
	"testing"
	"time"
)

type testCacheData struct {
	id   string
	data []byte
}

func (tc *testCacheData) GetNodeID() string {
	return tc.id
}

func (tc *testCacheData) GetNodeBytes() []byte {
	return tc.data
}

func TestAddAndGet(t *testing.T) {
	ttl := 3000
	cache := NewMemCacheMap(1, 1000, ttl)
	data := new(testCacheData)
	data.id = "1"
	data.data = []byte("1")
	err := cache.Add(data)
	if err != nil {
		t.Error("Error add data")
		return
	}
	err = cache.Add(data)
	if err == nil {
		t.Error("Error test duplicate add")
		return
	}
	data2 := &testCacheData{
		id:   "2",
		data: []byte("2"),
	}

	err = cache.Add(data2)
	if err != nil {
		t.Error("Error add data")
		return
	}
	d, err := cache.Get("1")
	if err != nil {
		t.Error("Error get data")
		return
	}
	if !bytes.Equal(d.GetNodeBytes(), data.data) {
		t.Error("Error get data with wrong content")
		return
	}
	d, err = cache.Get("3")
	if err == nil {
		t.Error("Error get not exist id")
		return
	}
	time.Sleep(time.Duration(ttl+500) * time.Millisecond)
	d, err = cache.Get("1")
	if err == nil {
		t.Error("Error test ttl ")
		return
	}
}
