package pokecache

import (
	"sync"
	"time"
)

// inSTRUCTional media
type Cache struct {
	entry    map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

//weeeeeee need the func! Gotta have that func!

// create a new cache and start reaploop goroutines
func NewCache(interval time.Duration) *Cache {

	cache := &Cache{
		entry:    make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.ReapLoop()

	return cache
}

// add entries to cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// read entries from cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entry[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// remove cache entries older than interval specified in NewCache(interval time.Duration) *Cache
func (c *Cache) ReapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entry {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}
