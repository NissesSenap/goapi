package cache

import (
	"net/http"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type response struct {
	header http.Header
	code   int
	body   []byte
}

type memCache struct {
	lock sync.RWMutex
	data map[string]response
}

var (
	cache          = memCache{data: map[string]response{}}
	cacheProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nr_served_cache_hits_total",
		Help: "The total number of times cached repsonse was successfull served",
	})
)

func set(resource string, response *response) {
	cache.lock.Lock()
	if response == nil {
		delete(cache.data, resource)
	} else {
		cache.data[resource] = *response
	}
	cache.lock.Unlock()
}

func get(resource string) *response {
	cache.lock.Lock()
	resp, ok := cache.data[resource]
	cache.lock.Unlock()
	if ok {
		return &resp
	}
	return nil
}

func copyHeader(src, dst http.Header) {
	for key, list := range src {
		for _, value := range list {
			dst.Add(key, value)
		}
	}
}

// MakeResource returns a valid cache resource
func MakeResource(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.TrimSuffix(r.URL.RequestURI(), "/")
}

// Clean removes all cache enteries
func Clean() {
	cache.lock.Lock()
	cache.data = map[string]response{}
	cache.lock.Unlock()
}

// Drop removes a cache entry
func Drop(res string) {
	set(res, nil)
}

// Serve functions returns true if a cached repsonse was successfull served
func Serve(w http.ResponseWriter, r *http.Request) bool {
	if w == nil || r == nil {
		return false
	}
	if r.Header.Get("Cache-Control") == "no-cache" {
		return false
	}
	resp := get(MakeResource(r))
	if resp == nil {
		return false
	}

	copyHeader(resp.header, w.Header())
	w.WriteHeader(resp.code)
	if r.Method != http.MethodHead {
		w.Write(resp.body)
	}
	cacheProcessed.Inc()
	return true
}
