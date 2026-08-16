package main

import (
	"sync"
	"time"
)

type cacheEntry struct {
	body      []byte
	expiresAt time.Time
}

type responseCache struct {
	entries    map[string]cacheEntry
	maxEntries int
	mutex      sync.RWMutex
}

func newResponseCache(maxEntries int) *responseCache {
	return &responseCache{
		entries:    make(map[string]cacheEntry),
		maxEntries: maxEntries,
	}
}

func (cache *responseCache) Get(key string) ([]byte, bool) {
	cache.mutex.RLock()
	entry, exists := cache.entries[key]
	cache.mutex.RUnlock()
	if !exists {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		cache.mutex.Lock()
		if current, ok := cache.entries[key]; ok && current.expiresAt.Equal(entry.expiresAt) {
			delete(cache.entries, key)
		}
		cache.mutex.Unlock()
		return nil, false
	}
	return append([]byte(nil), entry.body...), true
}

func (cache *responseCache) Set(key string, body []byte, lifetime time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	if len(cache.entries) >= cache.maxEntries {
		for cacheKey, entry := range cache.entries {
			if now.After(entry.expiresAt) {
				delete(cache.entries, cacheKey)
			}
		}
	}
	if len(cache.entries) >= cache.maxEntries {
		cache.entries = make(map[string]cacheEntry)
	}

	cache.entries[key] = cacheEntry{
		body:      append([]byte(nil), body...),
		expiresAt: now.Add(lifetime),
	}
}
