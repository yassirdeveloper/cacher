package main

import "sync"

type Set[K comparable] struct {
	mu    sync.Mutex
	items map[K]struct{}
}

func NewSet[K comparable]() *Set[K] {
	return &Set[K]{items: make(map[K]struct{})}
}

func (s *Set[K]) Add(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = struct{}{}
}

func (s *Set[K]) Remove(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

func (s *Set[K]) Contains(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.items[key]
	return exists
}

func (s *Set[K]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}

func (s *Set[K]) Items() []K {
	s.mu.Lock()
	defer s.mu.Unlock()
	var keys []K
	for k := range s.items {
		keys = append(keys, k)
	}
	return keys
}

func (s *Set[K]) AddMany(keys []K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range keys {
		s.items[k] = struct{}{}
	}
}
