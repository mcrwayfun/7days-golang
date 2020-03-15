package geecache

// A ByteView holds an immutable view of bytes
type ByteView struct {
	b []byte
}

// len returns the view's length
func (v ByteView) Len() int {
	return len(v.b)
}

// BytesSlice returns a copy of the data as a byte slice
func (v ByteView) ByteSlice() []byte {
	// b是只读的,使用cloneBytes返回一个拷贝,防止缓存值被外部程序修改
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
