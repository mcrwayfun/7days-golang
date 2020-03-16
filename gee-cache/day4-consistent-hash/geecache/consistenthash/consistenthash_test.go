package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 每个真实节点可有3个虚拟节点
	// 使用自定义的hash算法,返回对应的int
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 创建的虚拟节点为:02/12/22,04/14/24,06/16/26
	// 排序的虚拟节点为:02,04,06,12,14,16,22,24,26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	// 获取真实节点
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// Adds 8/18/28
	hash.Add("8")

	// 排序的虚拟节点为:02,04,06,08,12,14,16,18,22,24,26,28
	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
