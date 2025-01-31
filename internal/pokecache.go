package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entry    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
	ticker    <-chan time.Time
}

func NewCache(interval time.Duration) *Cache {

	cache := &Cache{
		entry:    make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.ReapLoop()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	return
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entry[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) ReapLoop() {
	for {
		c.mu.Lock()
		for key, entry := range c.entry {
			if int(time.Since(entry.createdAt)) >= int(c.interval) {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}
