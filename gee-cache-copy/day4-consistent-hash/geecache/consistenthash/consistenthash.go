package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对于每一个真实节点，创建对应个replicas的虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 通过添加编号来区分不同的虚拟节点
			node := strconv.Itoa(i) + key
			// 使用hash算法计算虚拟节点的hash值
			hash := int(m.hash([]byte(node)))
			// 添加到环上
			m.keys = append(m.keys, hash)
			// 在hashMap中增加虚拟节点和真实节点的映射
			m.hashMap[hash] = key
		}
	}

	// 环上的hash重新排序
	sort.Ints(m.keys)
}

// Get 选择节点
func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}

	// 计算key的hash值
	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点下标idx
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 因为m.keys是环状的，所以需要通过取余的方式获取下标。通过下标idx找到真正的index
	realIndex := idx % len(m.keys)
	// 从m.keys获取到对应的hash值
	realHash := m.keys[realIndex]
	// 通过hashMap获取到真实的节点
	return m.hashMap[realHash]
}
