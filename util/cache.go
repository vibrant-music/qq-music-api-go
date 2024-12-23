package util

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache = cache.New(5*time.Minute, 10*time.Minute)

// GetCache TODO check
func GetCache(key string) (interface{}, bool) {
	return Cache.Get(key)
}

func SetCache(key string, value interface{}, duration time.Duration) {
	Cache.Set(key, value, duration)
}
