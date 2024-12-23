// Problem: Implement a simple rate limiter that allows only n requests per second. 
// Use a goroutine to periodically reset the rate limiter and a mutex to synchronize 
// access to the counter that tracks the number of requests in the current second.

package main 

import (
	"sync"
	"fmt"
	"time"
)

type RateLimiter struct {
	mu sync.Mutex
	requests int 
	limit int 
}

// NewRateLimiter creates a new rate limiter with the given requests per second limit.
func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		limit: limit,
	}
	go rl.resetCounter()
	return rl
}

// Allow checks if a request is allowed.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.requests < rl.limit {
		rl.requests++
		return true
	}
	return false
}

// resetCounter resets the request counter every second.
func (rl *RateLimiter) resetCounter() {
	for {
		time.Sleep(1*time.Second)
		rl.mu.Lock()
		rl.requests = 0
		rl.mu.Unlock()
	}
}

func main() {
	n := 5
	rl := NewRateLimiter(n)

	for i := 1; i <= 10; i++ {
		if rl.Allow() {
			fmt.Printf("Request %d: Allowed\n", i)
		} else {
			fmt.Printf("Request %d: Denied\n", i)
		}
		time.Sleep(200 * time.Millisecond) // Some delay between requests.
	}
}