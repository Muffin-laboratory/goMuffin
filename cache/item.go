package cache

import "time"

// Item contains cache and cache's expiration.
type Item[T any] struct {
	Item       T
	Expiration time.Time
}

// Expired returns is expired the cache.
func (i *Item[T]) Expired() bool {
	return time.Now().After(i.Expiration)
}
