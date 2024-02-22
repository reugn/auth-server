package http

import (
	"log/slog"
	"net/netip"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

// IPWhiteList contains white list information for rate limiting.
type IPWhiteList struct {
	addresses map[string]*netip.Addr
	networks  []*netip.Prefix
	allowAny  bool
}

// NewIPWhiteList builds a new IPWhiteList from the list of IPs.
func NewIPWhiteList(ipList []string) (*IPWhiteList, error) {
	addresses := make(map[string]*netip.Addr)
	networks := make([]*netip.Prefix, 0)
	var allowAny bool
	for _, ip := range ipList {
		ip := strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		if strings.HasPrefix(ip, "0.0.0.0") {
			allowAny = true
		}
		network, err := netip.ParsePrefix(ip)
		if err != nil {
			ipAddr, err := netip.ParseAddr(ip)
			if err != nil {
				return nil, err
			}
			addresses[ip] = &ipAddr
		} else {
			networks = append(networks, &network)
		}
	}
	return &IPWhiteList{
		addresses: addresses,
		networks:  networks,
		allowAny:  allowAny,
	}, nil
}

func (wl *IPWhiteList) isAllowed(ip string) bool {
	if wl.allowAny {
		return true
	}
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		slog.Warn("Invalid client ip", "ip", ip, "err", err)
		return false
	}
	_, ok := wl.addresses[ip]
	if ok {
		return true
	}

	for _, network := range wl.networks {
		if network.Contains(ipAddr) {
			return true
		}
	}

	return false
}

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
