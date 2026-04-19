// Package idempotency deduplicates Task submissions by idempotency_key.
// The in-memory implementation is fine for a single-process server; swap
// for a durable store (Redis, Postgres) when calls may span processes.
package idempotency

import (
	"sync"

	"github.com/example/project/mcp-server/internal/primitives"
)

type Cache interface {
	Get(key string) (*primitives.Result, bool)
	Put(key string, result *primitives.Result)
}

// MemoryCache is a naive map + RWMutex. No eviction; rotate the process
// or swap for a real cache when memory pressure matters.
type MemoryCache struct {
	mu sync.RWMutex
	m  map[string]*primitives.Result
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{m: make(map[string]*primitives.Result)}
}

func (c *MemoryCache) Get(key string) (*primitives.Result, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.m[key]
	return r, ok
}

func (c *MemoryCache) Put(key string, result *primitives.Result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = result
}
