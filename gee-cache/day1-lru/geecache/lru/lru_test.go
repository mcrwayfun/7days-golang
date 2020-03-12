package lru

import "testing"

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	// lru := New(int64(0), nil)
	// lru.
}
