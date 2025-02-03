package main

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
