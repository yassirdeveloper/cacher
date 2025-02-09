package main

import (
	"fmt"
	"sync"
	"time"
)

const MinJanitorInterval = time.Second

type Janitor interface {
	Start()
	Stop()
	IsRunning() bool
	AdjustInterval(time.Duration)
	String() string
}

type janitor[K comparable, V any] struct {
	cache     Cache[K, V]
	interval  time.Duration
	isRunning bool
	stop      chan struct{}
	wg        sync.WaitGroup
	logger    Logger
}

func NewJanitor[K comparable, V any](cache Cache[K, V], interval time.Duration, logger Logger) (Janitor, Error) {
	if interval <= 0 {
		err := &SetupError{message: fmt.Sprintf("[JANITOR_EVENT] Invalid interval: must be greater than %s", MinJanitorInterval.String())}
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info(fmt.Sprintf("[JANITOR_EVENT] Initializing janitor for %s with interval: %v", cache.String(), interval))
	return &janitor[K, V]{
		cache:    cache,
		interval: interval,
		stop:     make(chan struct{}),
		logger:   logger,
	}, nil
}

func (j *janitor[K, V]) Start() {
	if j.IsRunning() {
		j.logger.Warning("[JANITOR_EVENT] Janitor is already running!")
		return
	}

	j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] Starting %s...", j.String()))
	j.wg.Add(1)
	j.isRunning = true
	go func() {
		defer j.wg.Done()
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] %s running...", j.String()))
				nbrExpired := j.cache.ClearExpired()
				j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] %s done, cleared %d keys", j.String(), nbrExpired))
			case <-j.stop:
				j.isRunning = false
				return
			}
		}
	}()
	j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] %s started", j.String()))
}

func (j *janitor[K, V]) Stop() {
	if !j.IsRunning() {
		j.logger.Warning("[JANITOR_EVENT] Janitor is not running!")
		return
	}

	j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] Stopping %s...", j.String()))
	close(j.stop)
	j.wg.Wait()
	j.isRunning = false
	j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] %s stopped", j.String()))
}

func (j *janitor[K, V]) IsRunning() bool {
	return j.isRunning
}

func (j *janitor[K, V]) AdjustInterval(newInterval time.Duration) {
	j.logger.Info(fmt.Sprintf("[JANITOR_EVENT] Adjusting interval for %s from %v to %v", j.String(), j.interval, newInterval))
	j.Stop()
	j.interval = newInterval
	j.Start()
}

func (j *janitor[K, V]) String() string {
	return fmt.Sprintf("Janitor for %s", j.cache.String())
}
