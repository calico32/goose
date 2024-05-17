package main

import "sync"

type mutexMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]V
}

func (m *mutexMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok = m.m[key]
	return
}

func (m *mutexMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = value
}

func (m *mutexMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

func NewMutexMap[K comparable, V any]() *mutexMap[K, V] {
	return &mutexMap[K, V]{
		m: make(map[K]V),
	}
}
