package cache

import (
	"context"
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	data          map[K]Item[V]
	onExpire      func(key K, value V)
	ttl           time.Duration
	checkInterval time.Duration

	mu sync.Mutex
}

type Item[V any] struct {
	Value  V
	Expiry time.Time
}

func (i Item[V]) isExpired() bool {
	return time.Now().After(i.Expiry)
}

func NewCache[K comparable, V any](ttl time.Duration, checkInterval time.Duration, onExpire func(key K, value V)) *Cache[K, V] {
	return &Cache[K, V]{
		data:          make(map[K]Item[V]),
		onExpire: onExpire,
		ttl:           ttl,
		checkInterval: checkInterval,
	}
}

func (c *Cache[K, V]) Run(ctx context.Context) error {
	timeCh := time.Tick(c.checkInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeCh:
			c.removeExpiredItems()
		}
	}
}

func (c *Cache[K, V]) removeExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.data {
		if item.isExpired() {
			delete(c.data, key)
			go c.onExpire(key, item.Value)
		}
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiry := time.Now().Add(c.ttl)
	item := Item[V]{
		Value:  value,
		Expiry: expiry,
	}

	c.data[key] = item
}

func (c *Cache[K, V]) Get(key K) (Item[V], bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.data[key]
	if !ok {
		return item, false
	}

	if item.isExpired() {
		delete(c.data, key)
		return item, false
	}

	return item, true
}

func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
