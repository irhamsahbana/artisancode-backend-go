package ratelimit

import (
	"sync"
	"time"
)

type AttemptLimiter interface {
	Allow(key string, limit int, window time.Duration) (bool, time.Duration)
}

type entry struct {
	Count   int
	ResetAt time.Time
}

type CacheLimiter struct {
	mu      sync.Mutex
	entries map[string]entry
}

func NewCacheLimiter() *CacheLimiter {
	return &CacheLimiter{
		entries: make(map[string]entry),
	}
}

func (l *CacheLimiter) Allow(key string, limit int, window time.Duration) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UTC()
	current, found := l.entries[key]
	if !found || now.After(current.ResetAt) {
		l.entries[key] = entry{
			Count:   1,
			ResetAt: now.Add(window),
		}
		return true, 0
	}

	if current.Count >= limit {
		return false, time.Until(current.ResetAt)
	}

	current.Count++
	l.entries[key] = current
	return true, 0
}
