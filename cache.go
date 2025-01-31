package main

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	data   map[K]V
	locker sync.RWMutex
}

func NewCahe[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]V),
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.data[key] = value
}

func (c *Cache[K, V]) Delete(key K) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.data, key)
}
