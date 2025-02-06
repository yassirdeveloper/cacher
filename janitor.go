package main

import (
	"sync"
	"time"
)

type Janitor[K comparable, V any] struct {
	cache    Cache[K, V]
	interval time.Duration
	stop     chan struct{}
	wg       sync.WaitGroup
}

func NewJanitor[K comparable, V any](cache Cache[K, V], interval time.Duration) *Janitor[K, V] {
	return &Janitor[K, V]{
		cache:    cache,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

func (j *Janitor[K, V]) Start() {
	j.wg.Add(1)
	go func() {
		defer j.wg.Done()
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.cache.ClearExpired()
			case <-j.stop:
				return
			}
		}
	}()
}

func (j *Janitor[K, V]) Stop() {
	close(j.stop)
	j.wg.Wait()
}
