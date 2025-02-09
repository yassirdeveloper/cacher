package main

import (
	"fmt"
	"sync"
	"time"
)

type CacheValue[V any] interface {
	Value() V
	ExpiresAt() time.Time
}

type cacheValue[V any] struct {
	value     *V
	expiresAt time.Time
}

func (c *cacheValue[V]) Value() V {
	return *c.value
}

func (c *cacheValue[V]) ExpiresAt() time.Time {
	return c.expiresAt
}

type Cache[K comparable, V any] interface {
	get(K) (CacheValue[V], bool)
	Get(K) (*V, bool)
	set(K, CacheValue[V])
	Set(K, V, time.Time)
	Delete(K)
	ClearExpired() int
	Clear()
	String() string
}

type cache[K comparable, V any] struct {
	data    map[K]CacheValue[V]
	records Records[K]
	locker  sync.RWMutex
	logger  Logger
}

func NewCache[K comparable, V any](precision time.Duration, logger Logger) (*cache[K, V], Error) {
	records, err := NewRecords[K](precision, logger)
	if err != nil {
		return nil, err
	}
	return &cache[K, V]{
		data:    make(map[K]CacheValue[V]),
		records: records,
		logger:  logger,
	}, nil
}

func (c *cache[K, V]) get(key K) (CacheValue[V], bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	cacheValue, exists := c.data[key]
	if !exists {
		return nil, false
	}
	return cacheValue, true
}

func (c *cache[K, V]) Get(key K) (*V, bool) {
	cacheValue, ok := c.get(key)
	if ok {
		if cacheValue.ExpiresAt().After(time.Now()) {
			c.logger.Warning(fmt.Sprintf("[MAIN_CACHE_EVENT] value with key %v exists after expired!", key))
			return nil, false
		}
		value := cacheValue.Value()
		return &value, true
	}
	c.logger.Info(fmt.Sprintf("[MAIN_CACHE_EVENT] Key not found: %v", key))
	return nil, false
}

func (c *cache[K, V]) set(key K, value CacheValue[V]) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if oldVal, exists := c.get(key); exists {
		c.records.Delete(key, oldVal.ExpiresAt())
	}
	c.data[key] = value
	c.records.Add(key, value.ExpiresAt())
}

func (c *cache[K, V]) Set(key K, value V, expiresAt time.Time) {
	if expiresAt.Before(time.Now()) {
		c.logger.Warning(fmt.Sprintf("[MAIN_CACHE_EVENT] Attempted to set key %v with past expiration time", key))
		return
	}
	c.set(key, &cacheValue[V]{value: &value, expiresAt: expiresAt})
}

func (c *cache[K, V]) Delete(key K) {
	c.locker.Lock()
	defer c.locker.Unlock()
	cacheValue, ok := c.get(key)
	if !ok {
		c.logger.Info(fmt.Sprintf("[MAIN_CACHE_EVENT] tried to delete inexistant key: %v", key))
		return
	}
	expiresAt := cacheValue.ExpiresAt()
	delete(c.data, key)
	c.records.Delete(key, expiresAt)
}

func (c *cache[K, V]) DeleteMany(keys []K) {
	for _, key := range keys {
		c.Delete(key)
	}
}

func (c *cache[K, V]) ClearExpired() int {
	now := time.Now()
	c.logger.Info("[MAIN_CACHE_EVENT] clearing expired keys...")
	nbrKeys := c.records.DeleteBefore(now, func(keys []K) {
		c.DeleteMany(keys)
	})
	c.logger.Info(fmt.Sprintf("[MAIN_CACHE_EVENT] clearing expired keys done: cleared %d keys", nbrKeys))
	return nbrKeys
}

func (c *cache[K, V]) Clear() {
	c.logger.Info("[MAIN_CACHE_EVENT] clearing all...")
	clear(c.data)
	c.records.Clear()
	c.logger.Info("[MAIN_CACHE_EVENT] clearing all done")
}

func (c *cache[K, V]) String() string {
	return "Main Cache"
}

type syncCache[K comparable, V any] struct {
	data    sync.Map
	records Records[K]
	logger  Logger
}

func NewSyncCache[K comparable, V any](precision time.Duration, logger Logger) (*syncCache[K, V], Error) {
	records, err := NewRecords[K](precision, logger)
	if err != nil {
		return nil, err
	}

	return &syncCache[K, V]{
		data:    sync.Map{},
		records: records,
		logger:  logger,
	}, nil
}

func (c *syncCache[K, V]) get(key K) (CacheValue[V], bool) {
	value, exists := c.data.Load(key)
	if exists {
		return value.(CacheValue[V]), true
	}
	return nil, false
}

func (c *syncCache[K, V]) Get(key K) (*V, bool) {
	cacheValue, ok := c.get(key)
	if ok {
		if cacheValue.ExpiresAt().After(time.Now()) {
			c.logger.Warning(fmt.Sprintf("[SYNC_CACHE_EVENT] value with key %v exists after expired!", key))
			return nil, false
		}
		value := cacheValue.Value()
		return &value, true
	}
	c.logger.Info(fmt.Sprintf("[SYNC_CACHE_EVENT] Key not found: %v", key))
	return nil, false
}

func (c *syncCache[K, V]) set(key K, cacheValue CacheValue[V]) {
	oldValue, ok := c.get(key)
	if ok {
		c.records.Delete(key, oldValue.ExpiresAt())
	}
	c.data.Store(key, cacheValue)
}

func (c *syncCache[K, V]) Set(key K, value V, expiresAt time.Time) {
	if expiresAt.Before(time.Now()) {
		c.logger.Warning(fmt.Sprintf("[SYNC_CACHE_EVENT] Attempted to set key %v with past expiration time", key))
		return
	}
	c.set(key, &cacheValue[V]{value: &value, expiresAt: expiresAt})
}

func (c *syncCache[K, V]) Delete(key K) {
	cacheValue, ok := c.get(key)
	if !ok {
		c.logger.Info(fmt.Sprintf("[SYNC_CACHE_EVENT] tried to delete inexistant key: %v", key))
		return
	}
	expiresAt := cacheValue.ExpiresAt()
	c.data.Delete(key)
	c.records.Delete(key, expiresAt)
}

func (c *syncCache[K, V]) DeleteMany(keys []K) {
	for _, key := range keys {
		c.Delete(key)
	}
}

func (c *syncCache[K, V]) ClearExpired() int {
	now := time.Now()
	c.logger.Info("[SYNC_CACHE_EVENT] clearing expired keys...")
	nbrKeys := c.records.DeleteBefore(now, func(keys []K) {
		c.DeleteMany(keys)
	})
	c.logger.Info(fmt.Sprintf("[SYNC_CACHE_EVENT] clearing expired keys done: cleared %d keys", nbrKeys))
	return nbrKeys
}

func (c *syncCache[K, V]) Clear() {
	c.logger.Info("[SYNC_CACHE_EVENT] clearing all...")
	c.data.Clear()
	c.logger.Info("[SYNC_CACHE_EVENT] clearing all...")
}

func (c *syncCache[K, V]) String() string {
	return "Sync Cache"
}
