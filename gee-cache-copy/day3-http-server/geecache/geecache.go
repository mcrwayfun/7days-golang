package geecache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g.Get(key)
}

type Group struct {
	name      string // 每个group的名字,不应该重复
	mainCache *Cache // key为name,value为name对应的缓存数据
	getter    Getter // 回调函数
}

var (
	mu     sync.RWMutex // 读多于写,使用读写锁即可
	groups map[string]*Group
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	mu.Lock()
	defer mu.Unlock()

	group := &Group{
		name:      name,
		mainCache: &Cache{maxBytes: cacheBytes},
		getter:    getter,
	}

	groups[name] = group
	return group
}

func GetGroup(name string) (*Group, error) {
	mu.RLock() // 使用读锁,因为name不会重复
	defer mu.RUnlock()
	if g, ok := groups[name]; ok {
		return g, nil
	}
	return nil, fmt.Errorf("fail to find the group")
}

func (g *Group) Get(key string) ([]byte, error) {
	if key == "" { // 外部调用时，需要对key作校验
		return nil, fmt.Errorf("key is required")
	}

	if view, ok := g.mainCache.get(key); ok {
		log.Printf("GeeCache hit [%s]\n", key)
		return view.ByteSlice(), nil
	}

	// 缓存中不存在,则去加载
	return g.load(key)
}

func (g *Group) load(key string) ([]byte, error) {
	return g.getLocally(key) // 还可以调用远程的方法
}

func (g *Group) getLocally(key string) ([]byte, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return nil, err
	}

	// 缓存到本地
	v := ByteView{bytes}
	g.populateCache(key, v)
	return bytes, nil
}

func (g *Group) populateCache(key string, v ByteView) {
	g.mainCache.add(key, v)
}