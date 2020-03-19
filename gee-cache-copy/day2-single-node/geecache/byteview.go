package geecache

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneByteView(v.b) // 只返回一个副本，防止外部数据的修改
}

func cloneByteView(b []byte) []byte {
	cp := make([]byte, len(b))
	copy(cp, b)
	return cp
}
