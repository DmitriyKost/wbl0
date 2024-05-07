package database

import (
	"log"
	"sync"
	"time"
)
var OrderCache *Cache = NewCache()

type Cache struct {
	entries map[string]*cacheEntry
	mutex   sync.RWMutex
	ticker  *time.Ticker
}

type cacheEntry struct {
	order      []byte
	expiration time.Time
}

func NewCache() *Cache {
	cache := &Cache{
		entries: make(map[string]*cacheEntry),
		ticker:  time.NewTicker(time.Minute), // Check for expired entries every minute
	}

	// Start a background goroutine to periodically check for expired entries
	go func() {
		for range cache.ticker.C {
			cache.cleanup()
		}
	}()

	return cache
}

// Set adds or updates a value in the cache with the specified UID and expiration time.
func (c *Cache) Set(orderUID string, order []byte, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[orderUID] = &cacheEntry{
		order:      order,
		expiration: time.Now().Add(expiration),
	}
}

// Retrieves the value associated with the specified key from the cache.
//
// If key does not exist in cache it attempts to get it from database.
//
// If the key does not exist in database returns nil.
func (c *Cache) Get(orderUID string) []byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

    var value []byte
	entry, ok := c.entries[orderUID]
	if !ok {
        log.Printf("%s - UID not found in cache\n", orderUID)
        var err error
        err, value = getOrder(orderUID)
        if err != nil {
            log.Printf("%s - UID not found in database\n", orderUID)
            return nil
        }
        go func() {
            c.Set(orderUID, value, time.Hour)
		}()
	} else {
        value = entry.order
    }
	return value
}

// removes expired entries from the cache.
func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, entry := range c.entries {
		if entry.expired() {
			delete(c.entries, key)
		}
	}
}

// returns true if the cache entry has expired.
func (e *cacheEntry) expired() bool {
	return !e.expiration.IsZero() && time.Now().After(e.expiration)
}
