package lru

import "testing"

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	cache := NewCache(int64(0))
	cache.Add("key1", String("123"))

	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "123" {
		t.Fatalf("Fail to hit key1=123")
	}

	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("Excepted miss, but success get")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"

	cap := len(k2 + k3 + v2 + v3)
	lru := NewCache(int64(cap))

	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get(k1); ok || lru.ll.Len() != 2 {
		t.Fatalf("Remove k1 failed")
	}
}

func TestAdd(t *testing.T) {
	lru := NewCache(int64(0))
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.nBytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nBytes)
	}
}