package geecache

// 根据传入的key选择对应节点的PeerGetter
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// 用于从对应的group中获取缓存值
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}


