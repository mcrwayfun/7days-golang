package geecache

import (
	"geecache/lru"
	"sync"
)

type Cache struct {
	mu  sync.Mutex
	lru *lru.Cache
}

func (c *Cache) add(key string, maxBytes int64, v ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil { // 延时加载
		c.lru = lru.NewCache(maxBytes)
	}

	c.lru.Add(key, v)
}

func (c *Cache) get(key string) (ByteView, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), nil
	}

	return ByteView{}, nil
}


