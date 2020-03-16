package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // hash环
	hashMap  map[int]string // 虚拟节点与真实节点的映射,键是虚拟节点的哈希值,值是真实节点
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 计算虚拟节点的hash名称(编号+key)
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 将节点加入到环中
			m.keys = append(m.keys, hash)
			// 维护虚拟节点与真实节点的映射
			m.hashMap[hash] = key
		}
	}
	// 环上的节点排序
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 { // 环上无节点
		return ""
	}

	// 使用hash函数计算key的hash值
	hash := int(m.hash([]byte(key)))
	// 使用二分搜索法查找对应的虚拟节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 返回真实的节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
