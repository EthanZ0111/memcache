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

func (tc *testCacheData) GetID() string {
	return tc.id
}

func (tc *testCacheData) GetBytes() ([]byte, error) {
	return tc.data, nil
}

func (tc *testCacheData) FillData() error {
	tc.data = []byte(tc.id)
	return nil
}

func (tc *testCacheData) GetError() error {
	return nil
}

func TestAddAndGet(t *testing.T) {
	ttl := 3000
	cache := NewMemCacheMap(1, 1000, ttl)
	data := new(testCacheData)
	data.id = "1"
	err := cache.AsyncAdd(data)
	if err != nil {
		t.Error("Error add data")
		return
	}
	err = cache.AsyncAdd(data)
	if err == nil {
		t.Error("Error test duplicate add")
		return
	}
	data2 := &testCacheData{
		id:   "2",
		data: []byte("2"),
	}

	err = cache.AsyncAdd(data2)
	if err != nil {
		t.Error("Error add data")
		return
	}
	d, err := cache.Get("1")
	if err != nil {
		t.Error("Error get data")
		return
	}
	getByte, err := d.GetBytes()
	oriByte, err := data.GetBytes()
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes.Equal(getByte, oriByte) {
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
