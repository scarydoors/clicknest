package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	data map[K]Item[V]
	ttl time.Duration

	mu sync.Mutex
}

type Item[V any] struct {
	value V
	expiry time.Time
}

func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]Item[V]),
		ttl: ttl,
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Unlock()	
	defer c.mu.Lock()

	expiry := time.Now().Add(c.ttl)
	item := Item[V]{
		value: value,
		expiry: expiry,
	}

	c.data[key] = item	
}
