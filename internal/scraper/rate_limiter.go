package scraper

import (
	"context"
	"sync"
	"time"
)

// RateLimiter controls the rate of requests
type RateLimiter struct {
	rate     int
	interval time.Duration
	tokens   int
	lastTime time.Time
	mutex    sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		interval: interval,
		tokens:   rate,
		lastTime: time.Now(),
	}
}

// Wait waits for the next available token
func (rl *RateLimiter) Wait(ctx context.Context) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime)

	// Add tokens based on elapsed time
	tokensToAdd := int(elapsed / rl.interval * time.Duration(rl.rate))
	rl.tokens += tokensToAdd
	if rl.tokens > rl.rate {
		rl.tokens = rl.rate
	}
	rl.lastTime = now

	if rl.tokens > 0 {
		rl.tokens--
		return nil
	}

	// Calculate wait time
	waitTime := rl.interval / time.Duration(rl.rate)

	select {
	case <-time.After(waitTime):
		rl.tokens = rl.rate - 1
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
