package main

import (
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Records[T any] interface {
	truncateTime(time.Time) int64
	Get(time.Time) []T
	Add(T, time.Time)
	AddMany([]T, time.Time)
	Delete(T, time.Time)
	DeleteBefore(time.Time, func([]T))
	DeleteAfter(time.Time, func([]T))
	Clear()
}

type records[K comparable] struct {
	data      *orderedmap.OrderedMap[int64, *Set[K]]
	precision time.Duration
}

func NewRecords[K comparable](precision time.Duration) Records[K] {
	return &records[K]{
		data:      orderedmap.New[int64, *Set[K]](),
		precision: precision,
	}
}

func (r *records[K]) truncateTime(t time.Time) int64 {
	return t.Truncate(r.precision).Unix()
}

func (r *records[K]) Get(t time.Time) []K {
	recordKey := r.truncateTime(t)
	set, ok := r.data.Get(recordKey)
	if !ok {
		return nil
	}
	return set.Items()
}

func (r *records[K]) Add(key K, t time.Time) {
	recordKey := r.truncateTime(t)
	s, ok := r.data.Get(recordKey)
	if ok {
		s.Add(key)
	} else {
		s := NewSet[K]()
		s.Add(key)
		r.data.Set(recordKey, s)
	}
}

func (r *records[K]) AddMany(keys []K, t time.Time) {
	recordKey := r.truncateTime(t)
	s, ok := r.data.Get(recordKey)
	if ok {
		s.AddMany(keys)
	} else {
		s := NewSet[K]()
		s.AddMany(keys)
		r.data.Set(recordKey, s)
	}
}

func (r *records[K]) Delete(key K, t time.Time) {
	recordKey := r.truncateTime(t)
	if recordKey == -1 {
		return
	}
	s, ok := r.data.Get(recordKey)
	if ok {
		s.Remove(key)
		if s.Len() == 0 {
			r.data.Delete(recordKey)
		}
	}
}

func (r *records[K]) DeleteBefore(t time.Time, deleteCallback func([]K)) {
	recordKey := r.truncateTime(t)
	for pair := r.data.GetPair(recordKey); pair != nil; pair = pair.Prev() {
		keys := pair.Value.Items()
		deleteCallback(keys)
		r.data.Delete(pair.Key)
	}
}

func (r *records[K]) DeleteAfter(t time.Time, deleteCallback func([]K)) {
	recordKey := r.truncateTime(t)
	for pair := r.data.GetPair(recordKey); pair != nil; pair = pair.Next() {
		keys := pair.Value.Items()
		deleteCallback(keys)
		r.data.Delete(pair.Key)
	}
}

func (r *records[K]) Clear() {
	r.data = orderedmap.New[int64, *Set[K]]()
}
