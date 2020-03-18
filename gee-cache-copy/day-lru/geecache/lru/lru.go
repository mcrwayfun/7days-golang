package lru

import "container/list"

type Cache struct {
	maxBytes int64 // 内存中允许最大存储的字节长度
	nBytes   int64 // 当前内存中存储的字节长度
	ll       *list.List
	cache    map[string]*list.Element // 存储key与节点的映射关系,查询时间复杂度降低至O(1)
}

// 计算长度时，key和value都需要计算
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New Cache
func NewCache(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		nBytes:   0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

// Add Cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 已经存在key
		// 移动到队头（假设队头存放的是最近使用的元素）
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 重新计算nByte（此刻key没变，计算不用考虑key）
		c.nBytes += int64(kv.value.Len()) - int64(value.Len())
		// 存在则替换value
		kv.value = value
	} else { // 不存在
		// 在队头新增一个元素
		ele := c.ll.PushFront(&entry{key, value})
		// 重新计算nByte
		c.nBytes += int64(len(key)) + int64(value.Len())
		// 设置到cache中
		c.cache[key] = ele
	}

	// 因为计算了nByte,所以需要判断
	// 循环的方式,移除所有不符合要求的节点
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.removeOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	return
}

func (c *Cache) removeOldest() {
	// 获取队尾元素
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		// 将队尾元素从链表中移除
		c.ll.Remove(ele)
		// 将元素从map中删除
		delete(c.cache, kv.key)
		// 重新计算nBytes
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
	}
}

// len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
