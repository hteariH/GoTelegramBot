package main

import "sync"

type Cache struct {
	store map[string][]byte
	mu    sync.RWMutex
}

// NewCache creates a new cache struct.
func NewCache() *Cache {
	return &Cache{
		store: make(map[string][]byte),
	}
}

// Get retrieves an object from cache.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, found := c.store[key]
	return val, found
}

// Set adds an object to cache.
func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}
