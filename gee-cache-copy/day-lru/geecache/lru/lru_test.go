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
