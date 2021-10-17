package main

import (
	"sync"

	"golang.org/x/time/rate"
)

// IPAddress represents an IP address string.
type IPAddress string

// IPRateLimiter represents a rate limiter based on an IP address.
type IPRateLimiter struct {
	sync.Mutex
	limiters        map[IPAddress]*rate.Limiter
	tokensPerSecond rate.Limit
	tokenBucketSize int
}

// NewIPRateLimiter returns a new IPRateLimiter.
func NewIPRateLimiter(tps rate.Limit, size int) *IPRateLimiter {
	ipLimiter := &IPRateLimiter{
		limiters:        make(map[IPAddress]*rate.Limiter),
		tokensPerSecond: tps,
		tokenBucketSize: size,
	}

	return ipLimiter
}

// AddLimiter creates a new rate limiter and adds it to the limiters map,
// using the IP address as the key.
func (ipLimiter *IPRateLimiter) AddLimiter(ipAddr string) *rate.Limiter {
	ipLimiter.Lock()
	defer ipLimiter.Unlock()

	limiter := rate.NewLimiter(ipLimiter.tokensPerSecond, ipLimiter.tokenBucketSize)

	ipLimiter.limiters[IPAddress(ipAddr)] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for the provided IP address if it exists.
// Otherwise calls AddLimiter to add a new limiter to the map.
func (ipLimiter *IPRateLimiter) GetLimiter(ipAddr string) *rate.Limiter {
	ipLimiter.Lock()
	limiter, exists := ipLimiter.limiters[IPAddress(ipAddr)]

	if !exists {
		ipLimiter.Unlock()
		return ipLimiter.AddLimiter(ipAddr)
	}

	ipLimiter.Unlock()

	return limiter
}
