package geecache

import (
	"fmt"
	"log"
	"net/http"
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
