package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HttPPool struct {
	self        string
	basePath    string
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
	mu          sync.Mutex
}

func NewHttPPool(self string) *HttPPool {
	return &HttPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (h *HttPPool) Log(format string, v ...interface{}) {
	log.Printf("[Http Server %s] %s", h.self, fmt.Sprintf(format, v...))
}

func (h *HttPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断访问路径的前缀是否为basePath,不是则返回错误
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("excepted r.URL.Path should be" + h.basePath)
	}

	h.Log("%s %s", r.Method, r.URL.Path)

	// 约定访问路径为/<basepath>/<gourpname>/<key>
	sp := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)

	groupName := sp[0]
	key := sp[1]

	// 通过group name 获取group实例
	group, err := GetGroup(groupName)
	if err != nil {
		http.Error(w, "no such group", 404)
		return
	}

	// 使用group.Get获取缓存
	bytes, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(bytes.ByteSlice())
}

// 实现HttPPool的客户端

// HttpGetter是客户端HttPPool发请求的主体
type httpGetter struct {
	baseURL string // 表示将要访问的节点地址,http://example.com/_geecache/
}

// 实现客户端的Get方法
func (h *httpGetter) Get(group, key string) ([]byte, error) {
	// http://example.com/_geecache/<gourpname>/<key>
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("excepted httpCode 200, but get %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// HttPPool 设置节点
func (h *HttPPool) Set(peers ...string) {
	mu.Lock()
	defer mu.Unlock()
	// 根据默认的虚拟节点倍数初始化一致性哈希算法Map
	h.peers = consistenthash.New(defaultReplicas, nil)
	// 将所有的真实节点添加到环上
	h.peers.Add(peers...)
	// 为每一个真实节点初始化一个httpGetter的方法，baseURL = peer + p.basePath
	h.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

// 实现PickerPeer,通过Get方法,根据key返回对应的Http客户端
func (h *HttPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.peers != nil { // 已经设置了节点
		// peer不为空且不等于本地节点
		if peer := h.peers.Get(key); peer != "" && peer != h.self {
			log.Printf("get peer success %s", peer)
			return h.httpGetters[peer], true
		}
	}
	return nil, false
}

var _ PeerGetter = (*httpGetter)(nil)

var _ PeerPicker = (*HttPPool)(nil)