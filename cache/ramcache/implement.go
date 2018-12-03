package ramcache

import (
	"github.com/finalist736/gokit/cache"
	rCache "github.com/pmylund/go-cache"
	"time"
)

type RamCache struct {
	cache *rCache.Cache
}

func New(defaultExpiration, cleanupInterval time.Duration) cache.CacheService {
	return &RamCache{cache: rCache.New(defaultExpiration, cleanupInterval)}
}

func (s *RamCache) Get(key string) (interface{}, bool) {
	return s.cache.Get(key)
}

func (s *RamCache) Set(key string, value interface{}) {
	s.cache.Set(key, value, rCache.DefaultExpiration)
}

func (s *RamCache) Delete(key string) {
	s.cache.Delete(key)
}

func (s *RamCache) Flush() {
	s.cache.Flush()
}
