package geecache

import (
	"geecache/lru"
	"sync"
)

type Cache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	maxBytes int64
}

func (c *Cache) add(key string, v ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil { // 延时加载
		c.lru = lru.NewCache(c.maxBytes)
	}

	c.lru.Add(key, v)
}

func (c *Cache) get(key string) (v ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {// 懒加载 需要加判断
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return ByteView{}, false
}


