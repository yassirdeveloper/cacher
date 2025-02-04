package main

import (
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Delete(K)
	Clear()
}

type cache[K comparable, V any] struct {
	data   map[K]V
	locker sync.RWMutex
}

func NewCache[K comparable, V any]() *cache[K, V] {
	return &cache[K, V]{
		data: make(map[K]V),
	}
}

func (c *cache[K, V]) Get(key K) (V, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	value, exists := c.data[key]
	return value, exists
}

func (c *cache[K, V]) Set(key K, value V) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.data[key] = value
}

func (c *cache[K, V]) Delete(key K) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.data, key)
}

func (c *cache[K, V]) Clear() {
	clear(c.data)
}

type syncCache[K comparable, V any] struct {
	data sync.Map
}

func NewSyncCache[K comparable, V any]() *syncCache[K, V] {
	return &syncCache[K, V]{}
}

func (c *syncCache[K, V]) Get(key K) (V, bool) {
	value, exists := c.data.Load(key)
	if !exists {
		var zero V // Return zero value of V if key doesn't exist
		return zero, false
	}
	return value.(V), true // Type assertion
}

func (c *syncCache[K, V]) Set(key K, value V) {
	c.data.Store(key, value)
}

func (c *syncCache[K, V]) Delete(key K) {
	c.data.Delete(key)
}

func (c *syncCache[K, V]) Clear() {
	c.data.Clear()
}

type CacheManager[K comparable, V any] interface {
	Get(bool) Cache[K, V]
	Clear()
}

type cacheManager[K comparable, V any] struct {
	cache     Cache[K, V]
	syncCache Cache[K, V]
}

func NewCacheManager[K comparable, V any]() CacheManager[K, V] {
	return &cacheManager[K, V]{
		cache:     NewCache[K, V](),
		syncCache: NewSyncCache[K, V](),
	}
}

func (cm *cacheManager[K, V]) Get(frequentAccess bool) Cache[K, V] {
	if frequentAccess {
		return cm.syncCache
	}
	return cm.cache
}

func (cm *cacheManager[K, V]) Clear() {
	cm.cache.Clear()
	cm.syncCache.Clear()
}
