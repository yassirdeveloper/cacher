package main

import (
	"time"
)

type CacheManager[K comparable, V any] interface {
	Get(bool) Cache[K, V]
	SetupMainCache(time.Duration) Error
	SetupMainCacheJanitor(time.Duration) Error
	SetupSyncCache(time.Duration) Error
	SetupSyncCacheJanitor(time.Duration) Error
	StartJanitors()
	StopJanitors()
	ClearCaches()
}

type cacheManager[K comparable, V any] struct {
	cache            Cache[K, V]
	syncCache        Cache[K, V]
	cacheJanitor     Janitor
	syncCacheJanitor Janitor
	logger           Logger
}

func NewCacheManager[K comparable, V any](logger Logger) CacheManager[K, V] {
	return &cacheManager[K, V]{logger: logger}
}

func (cm *cacheManager[K, V]) Get(frequentAccess bool) Cache[K, V] {
	if frequentAccess {
		return cm.syncCache
	}
	return cm.cache
}

func (cm *cacheManager[K, V]) SetupMainCache(precision time.Duration) Error {
	cm.logger.Info("Setting up main cache...")
	cache, err := NewCache[K, V](precision, cm.logger)
	if err != nil {
		return err
	}
	cm.cache = cache
	cm.logger.Info("Main cache setup done")
	return nil
}

func (cm *cacheManager[K, V]) SetupMainCacheJanitor(interval time.Duration) Error {
	cm.logger.Info("Setting up main cache janitor...")
	janitor, err := NewJanitor(cm.cache, interval, cm.logger)
	if err != nil {
		return err
	}
	cm.cacheJanitor = janitor
	cm.logger.Info("Main cache janitor setup done")
	return nil
}

func (cm *cacheManager[K, V]) SetupSyncCache(precision time.Duration) Error {
	cm.logger.Info("Setting up sync cache...")
	syncCache, err := NewSyncCache[K, V](precision, cm.logger)
	if err != nil {
		return err
	}
	cm.syncCache = syncCache
	cm.logger.Info("Sync cache setup done")
	return nil
}

func (cm *cacheManager[K, V]) SetupSyncCacheJanitor(interval time.Duration) Error {
	cm.logger.Info("Setting up sync cache janitor...")
	janitor, err := NewJanitor(cm.cache, interval, cm.logger)
	if err != nil {
		return err
	}
	cm.syncCacheJanitor = janitor
	cm.logger.Info("Sync cache janitor setup done")
	return nil
}

func (cm *cacheManager[K, V]) StartJanitors() {
	cm.logger.Info("Starting janitors...")
	cm.cacheJanitor.Start()
	cm.syncCacheJanitor.Start()
	cm.logger.Info("Janitors started successfully")
}

func (cm *cacheManager[K, V]) StopJanitors() {
	cm.logger.Info("Stopping janitors...")
	cm.cacheJanitor.Stop()
	cm.syncCacheJanitor.Stop()
	cm.logger.Info("Janitors stopped successfully")
}

func (cm *cacheManager[K, V]) ClearCaches() {
	cm.logger.Info("Clearing caches...")
	cm.cache.Clear()
	cm.syncCache.Clear()
	cm.logger.Info("Caches cleared successfully")
}
