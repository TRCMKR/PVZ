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
	currentCapacity int
	capacity        int
	cache           map[K]*list.Element
	dataList        *list.List
	mu              sync.Mutex
}

// NewCache ...
func NewCache[K cacheKey, T any](capacity int) *Cache[K, T] {
	return &Cache[K, T]{
		currentCapacity: 0,
		capacity:        capacity,
		cache:           make(map[K]*list.Element, capacity),
		dataList:        list.New(),
	}
}

// Put ...
func (c *Cache[K, T]) Put(key K, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.cache[key]; ok {
		item.Value.(*cacheItem[K, T]).value = value
		c.dataList.MoveToFront(item)

		return
	}

	item := &cacheItem[K, T]{key, value}
	c.cache[key] = c.dataList.PushFront(item)
	c.currentCapacity++
	if c.currentCapacity > c.capacity {
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

	c.dataList.MoveToFront(item)

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

	c.dataList.Remove(item)
	delete(c.cache, key)
	c.currentCapacity--
}

func (c *Cache[K, T]) removeOldest() {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldest := c.dataList.Back()
	if oldest != nil {
		c.dataList.Remove(oldest)
		key := oldest.Value.(*cacheItem[K, T]).key
		delete(c.cache, key)
		c.currentCapacity--
	}
}

// GetAll ...
func (c *Cache[K, T]) GetAll() []T {
	c.mu.Lock()
	defer c.mu.Unlock()

	items := make([]T, 0, c.currentCapacity)
	for _, v := range c.cache {
		item := v.Value.(*cacheItem[K, T]).value
		items = append(items, item)
	}

	return items
}

// GetAllBy ...
func (c *Cache[K, T]) GetAllBy(op func(T) (bool, error)) []T {
	c.mu.Lock()
	defer c.mu.Unlock()

	items := make([]T, 0, c.currentCapacity)
	for _, v := range c.cache {
		item := v.Value.(*cacheItem[K, T]).value
		ok, _ := op(item)
		if ok {
			items = append(items, item)
		}
	}

	return items
}
