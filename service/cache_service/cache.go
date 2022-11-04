package cache_service

import (
	"time"

	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache

func Init() {
	Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
}
