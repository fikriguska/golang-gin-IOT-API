package cache_service

import (
	"fmt"

	"github.com/dgraph-io/ristretto"
)

// var Cache *ristretto.Cache

var cache, _ = ristretto.NewCache(&ristretto.Config{
	NumCounters: 1e7,     // number of keys to track frequency of (10M).
	MaxCost:     1 << 30, // maximum cost of cache (1GB).
	BufferItems: 64,      // number of keys per Get buffer.
})

func Get(prefix string, id int) (interface{}, bool) {
	key := fmt.Sprintf("%d-node", id)
	return cache.Get(key)
}

func Set(prefix string, id int, val interface{}) {
	key := fmt.Sprintf("%d-node", id)
	cache.Set(key, val, 0)
}

func Del(prefix string, id int) {
	key := fmt.Sprintf("%d-node", id)
	cache.Del(key)
}
