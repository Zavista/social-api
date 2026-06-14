package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	mu      sync.Mutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	count, exists := rl.clients[ip]

	// starts a new goroutine to reset the count after duration
	// basically a TTL for the rate
	if !exists {
		go rl.resetCount(ip)
	}

	if count >= rl.limit {
		return false, rl.window
	}

	rl.clients[ip]++
	return true, 0
}

func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.mu.Lock()
	delete(rl.clients, ip)
	rl.mu.Unlock()
}
