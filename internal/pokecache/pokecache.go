package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mx      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) PrintCache() {
	for key, value := range c.entries {
		fmt.Printf("%s at %s", string(value.val), key)
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.mx.Lock()
	defer c.mx.Unlock()
	var ce cacheEntry
	ce.createdAt = time.Now()
	ce.val = val

	c.entries[key] = ce
}

func (c *Cache) Get(key string) (val []byte, recieved bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	ce, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return ce.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.mx.Lock()
		for key, ce := range c.entries {
			if time.Since(ce.createdAt) >= interval {
				delete(c.entries, key)
			}
		}
		c.mx.Unlock()
	}
}
