package cache

import (
	"sync"
	"time"
)

// CacheManager manages caches.
type CacheManager[K comparable, V any] struct {
	data     map[K]*Item[V]
	mu       sync.RWMutex
	Duration time.Duration
}

// New Creates new cache manager.
func New[K comparable, V any](duration time.Duration) *CacheManager[K, V] {
	cacheManager := CacheManager[K, V]{
		data: make(map[K]*Item[V]),
	}

	cacheManager.startCleanupExpiredCaches()

	return &cacheManager
}

// Set creates new cache.
func (m *CacheManager[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; ok {
		return
	}

	m.data[key] = &Item[V]{value, time.Now().Add(m.Duration)}
}

// Get returns cache.
func (m *CacheManager[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if item, ok := m.data[key]; ok {
		return item.Item, true
	}

	var zero V
	return zero, false
}

// All returns all caches.
func (m *CacheManager[K, V]) All() map[K]V {
	cacheList := make(map[K]V, len(m.data))

	for key, cache := range m.data {
		cacheList[key] = cache.Item
	}

	return cacheList
}

// Delete deletes cache.
func (m *CacheManager[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *CacheManager[K, V]) startCleanupExpiredCaches() {
	ticker := time.NewTicker(time.Minute * 10)
	go func() {
		for range ticker.C {
			for key, item := range m.data {
				if item.Expired() {
					m.Delete(key)
				}
			}
		}
	}()
}
