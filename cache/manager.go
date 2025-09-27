package cache

import (
	"sync"
	"time"
)

// CacheManager manages caches.
type CacheManager[T any] struct {
	data     map[string]*Item[T]
	mu       sync.RWMutex
	Duration time.Duration
}

// New Creates new cache manager.
func New[T any](duration time.Duration) *CacheManager[T] {
	cacheManager := CacheManager[T]{
		data: make(map[string]*Item[T]),
	}

	cacheManager.startCleanupExpiredCaches()

	return &cacheManager
}

// Set creates new cache.
func (m *CacheManager[T]) Set(key string, value T) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; ok {
		return
	}

	m.data[key] = &Item[T]{value, time.Now().Add(m.Duration)}
}

// Get returns cache.
func (m *CacheManager[T]) Get(key string) (T, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		return item.Item, true
	}

	var zero T
	return zero, false
}

// Delete deletes cache.
func (m *CacheManager[T]) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *CacheManager[T]) startCleanupExpiredCaches() {
	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		for range ticker.C {
			m.mu.Lock()

			for key, item := range m.data {
				if item.Expired() {
					m.Delete(key)
				}
			}

			m.mu.Unlock()
		}
	}()
}
