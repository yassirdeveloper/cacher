package main

import (
	"fmt"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

const MinPrecision = time.Second

type Records[T any] interface {
	truncateTime(time.Time) int64
	Get(time.Time) []T
	Add(T, time.Time)
	AddMany([]T, time.Time)
	Delete(T, time.Time)
	DeleteBefore(time.Time, func([]T)) int
	DeleteAfter(time.Time, func([]T)) int
	Clear()
}

type records[K comparable] struct {
	data      *orderedmap.OrderedMap[int64, *Set[K]]
	precision time.Duration
	logger    Logger
}

func NewRecords[K comparable](precision time.Duration, logger Logger) (Records[K], Error) {
	if precision <= MinPrecision {
		err := &SetupError{message: fmt.Sprintf("Invalid precision value: must be greater than %s", MinPrecision.String())}
		logger.Error(fmt.Sprintf("[RECORDS_EVENT] %s", err.Error()))
		return nil, err
	}

	logger.Info("[RECORDS_EVENT] Initializing records with precision: " + precision.String())
	return &records[K]{
		data:      orderedmap.New[int64, *Set[K]](),
		precision: precision,
		logger:    logger,
	}, nil
}

func (r *records[K]) truncateTime(t time.Time) int64 {
	return t.Truncate(r.precision).Unix()
}

func (r *records[K]) Get(t time.Time) []K {
	recordKey := r.truncateTime(t)
	set, ok := r.data.Get(recordKey)
	if !ok {
		r.logger.Warning(fmt.Sprintf("[RECORDS_EVENT] No record found for key %v", recordKey))
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
		r.logger.Warning(fmt.Sprintf("[RECORDS_EVENT] Attempted to delete key %v with invalid record", key))
		return
	}
	s, ok := r.data.Get(recordKey)
	if !ok {
		r.logger.Warning(fmt.Sprintf("[RECORDS_EVENT] Attempted to delete key %v from non-existent record %v", key, recordKey))
		return
	}
	s.Remove(key)
	if s.Len() == 0 {
		r.data.Delete(recordKey)
	}
}

func (r *records[K]) DeleteBefore(t time.Time, deleteCallback func([]K)) int {
	nbrKeys := 0
	recordKey := r.truncateTime(t)
	r.logger.Info(fmt.Sprintf("[RECORDS_EVENT] Deleting records before %v", t))
	for pair := r.data.GetPair(recordKey); pair != nil; pair = pair.Prev() {
		keys := pair.Value.Items()
		deleteCallback(keys)
		r.data.Delete(pair.Key)
		nbrKeys++
	}
	r.logger.Info(fmt.Sprintf("[RECORDS_EVENT] Deleted %d records before %v", nbrKeys, t))
	return nbrKeys
}

func (r *records[K]) DeleteAfter(t time.Time, deleteCallback func([]K)) int {
	nbrKeys := 0
	recordKey := r.truncateTime(t)
	r.logger.Info(fmt.Sprintf("[RECORDS_EVENT] Deleting records after %v", t))
	for pair := r.data.GetPair(recordKey); pair != nil; pair = pair.Next() {
		keys := pair.Value.Items()
		deleteCallback(keys)
		r.data.Delete(pair.Key)
		nbrKeys++
	}
	r.logger.Info(fmt.Sprintf("[RECORDS_EVENT] Deleted %d records after %v", nbrKeys, t))
	return nbrKeys
}

func (r *records[K]) Clear() {
	r.logger.Info("[RECORDS_EVENT] Clearing all records...")
	r.data = orderedmap.New[int64, *Set[K]]()
	r.logger.Info("[RECORDS_EVENT] All records cleared")
}
