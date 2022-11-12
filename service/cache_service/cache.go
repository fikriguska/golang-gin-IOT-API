package cache_service

import (
	"github.com/dgraph-io/ristretto"
)

// var Cache *
var Cache *ristretto.Cache

func Init() {
	// Cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	Cache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	// if err != nil {
	// 	panic(err)
	// }

	// // set a value with a cost of 1
	// cache.Set("key", "value", 1)

	// // wait for value to pass through buffers
	// time.Sleep(10 * time.Millisecond)

	// value, found := cache.Get("key")
	// if !found {
	// 	panic("missing value")
	// }
	// fmt.Println(value)
	// cache.Del("key")
}
