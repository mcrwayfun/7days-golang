package geecache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var basePath = "/_geecache/"

type HttPPool struct {
	self     string
	basePath string
}

func NewHttPPool(self string) *HttPPool {
	return &HttPPool{
		self:     self,
		basePath: basePath,
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
		http.Error(w, "no such key", 404)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(bytes.ByteSlice())
}

// 实现HttPPool的客户端

// HttpGetter是客户端HttPPool发请求的主体
type HttpGetter struct {
	baseURL string // 表示将要访问的节点地址,http://example.com/_geecache/
}

// 实现客户端的Get方法
func (h *HttpGetter) Get(group, key string) ([]byte, error) {
	// http://example.com/_geecache/<gourpname>/<key>
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("excepted httpCode 200, but get %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

