package main

import (
	"sync"
	"time"
)

type CacheValue[V any] interface {
	ExpiresAt() time.Time
	Value() V
}

type cacheValue[V any] struct {
	value     *V
	expiresAt time.Time
}

type Cache[K comparable, V any] interface {
	Get(K) (CacheValue[V], bool)
	Set(K, CacheValue[V])
	Delete(K)
	ClearExpired()
	Clear()
}

type cache[K comparable, V any] struct {
	data    map[K]CacheValue[V]
	records Records[K]
	locker  sync.RWMutex
}

func NewCache[K comparable, V any](precision time.Duration) *cache[K, V] {
	return &cache[K, V]{
		data:    make(map[K]CacheValue[V]),
		records: NewRecords[K](precision),
	}
}

func (c *cache[K, V]) Get(key K) (CacheValue[V], bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	cacheValue, exists := c.data[key]
	if !exists {
		return nil, false
	}
	return cacheValue, true
}

func (c *cache[K, V]) Set(key K, value CacheValue[V]) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if oldVal, exists := c.Get(key); exists {
		c.records.Delete(key, oldVal.ExpiresAt())
	}
	c.data[key] = value
	c.records.Add(key, value.ExpiresAt())
}

func (c *cache[K, V]) Delete(key K) {
	c.locker.Lock()
	defer c.locker.Unlock()
	cacheValue, ok := c.Get(key)
	if ok {
		expiresAt := cacheValue.ExpiresAt()
		delete(c.data, key)
		c.records.Delete(key, expiresAt)
	}
}

func (c *cache[K, V]) DeleteMany(keys []K) {
	for _, key := range keys {
		c.Delete(key)
	}
}

func (c *cache[K, V]) ClearExpired() {
	now := time.Now()
	c.records.DeleteBefore(now, func(keys []K) {
		c.DeleteMany(keys)
	})
}

func (c *cache[K, V]) Clear() {
	clear(c.data)
	c.records.Clear()
}

type syncCache[K comparable, V any] struct {
	data    sync.Map
	records Records[K]
}

func NewSyncCache[K comparable, V any](precision time.Duration) *syncCache[K, V] {
	return &syncCache[K, V]{
		data:    sync.Map{},
		records: NewRecords[K](precision),
	}
}

func (c *syncCache[K, V]) Get(key K) (CacheValue[V], bool) {
	value, exists := c.data.Load(key)
	if exists {
		return value.(CacheValue[V]), true
	}
	return nil, false
}

func (c *syncCache[K, V]) Set(key K, cacheValue CacheValue[V]) {
	oldValue, ok := c.Get(key)
	if ok {
		c.records.Delete(key, oldValue.ExpiresAt())
	}
	c.data.Store(key, cacheValue)
}

func (c *syncCache[K, V]) Delete(key K) {
	cacheValue, ok := c.data.Load(key)
	if ok {
		expiresAt := cacheValue.(CacheValue[V]).ExpiresAt()
		c.data.Delete(key)
		c.records.Delete(key, expiresAt)
	}
}

func (c *syncCache[K, V]) DeleteMany(keys []K) {
	for _, key := range keys {
		c.Delete(key)
	}
}

func (c *syncCache[K, V]) ClearExpired() {
	now := time.Now()
	c.records.DeleteBefore(now, func(keys []K) {
		c.DeleteMany(keys)
	})
}

func (c *syncCache[K, V]) Clear() {
	c.data.Clear()
}

type CacheManager[K comparable, V any] interface {
	Get(bool) Cache[K, V]
	SetupMainCache(time.Duration)
	SetupSyncCache(time.Duration)
	Clear()
}

type cacheManager[K comparable, V any] struct {
	cache     Cache[K, V]
	syncCache Cache[K, V]
}

func NewCacheManager[K comparable, V any]() CacheManager[K, V] {
	return &cacheManager[K, V]{}
}

func (cm *cacheManager[K, V]) Get(frequentAccess bool) Cache[K, V] {
	if frequentAccess {
		return cm.syncCache
	}
	return cm.cache
}

func (cm *cacheManager[K, V]) SetupMainCache(precision time.Duration) {
	cm.cache = NewCache[K, V](precision)
}

func (cm *cacheManager[K, V]) SetupSyncCache(precision time.Duration) {
	cm.cache = NewSyncCache[K, V](precision)
}

func (cm *cacheManager[K, V]) Clear() {
	cm.cache.Clear()
	cm.syncCache.Clear()
}
