package lru

import (
	"container/list"
	"sync"
)

type cacheKey interface {
	string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type cacheItem[K cacheKey, T any] struct {
	key   K
	value T
}

// Cache ...
type Cache[K cacheKey, T any] struct {
	capacity int
	cache    map[K]*list.Element
	list     *list.List
	mu       sync.Mutex
}

// NewCache ...
func NewCache[K cacheKey, T any](capacity int) *Cache[K, T] {
	return &Cache[K, T]{
		capacity: capacity,
		cache:    make(map[K]*list.Element, capacity),
		list:     list.New(),
	}
}

// Put ...
func (c *Cache[K, T]) Put(key K, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.cache[key]; ok {
		item.Value.(*cacheItem[K, T]).value = value
		c.list.MoveToFront(item)

		return
	}

	item := &cacheItem[K, T]{key, value}
	c.cache[key] = c.list.PushFront(item)
	if len(c.cache) > c.capacity {
		c.removeOldest()
	}
}

// Get ...
func (c *Cache[K, T]) Get(key K) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T
	item, ok := c.cache[key]
	if !ok {
		return zero, false
	}

	c.list.MoveToFront(item)

	return item.Value.(*cacheItem[K, T]).value, true
}

// Remove ...
func (c *Cache[K, T]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.cache[key]
	if !ok {
		return
	}

	c.list.Remove(item)
	delete(c.cache, key)
}

func (c *Cache[K, T]) removeOldest() {
	oldest := c.list.Back()
	if oldest != nil {
		c.list.Remove(oldest)
		key := oldest.Value.(*cacheItem[K, T]).key
		delete(c.cache, key)
	}
}

// GetAll ...
func (c *Cache[K, T]) GetAll() []T {
	c.mu.Lock()
	defer c.mu.Unlock()

	items := make([]T, 0, len(c.cache))
	for _, v := range c.cache {
		item := v.Value.(*cacheItem[K, T]).value
		items = append(items, item)
	}

	return items
}

// GetAllBy ...
func (c *Cache[K, T]) GetAllBy(op func(T) (bool, error)) ([]T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	items := make([]T, 0, len(c.cache))
	for _, v := range c.cache {
		item := v.Value.(*cacheItem[K, T]).value
		ok, err := op(item)
		if err != nil {
			return nil, err
		}
		if ok {
			items = append(items, item)
		}
	}

	return items, nil
}
