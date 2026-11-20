package main

import "sync"

type LinkCache struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewLinkCache() *LinkCache {
	return &LinkCache{
		items: make(map[string]string),
	}
}

func (cache *LinkCache) Set(slug, originalLink string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if len(cache.items) >= 10000 {
		clear(cache.items)
	}
	cache.items[slug] = originalLink
}

func (cache *LinkCache) Get(slug string) (string, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	url, ok := cache.items[slug]
	return url, ok
}
