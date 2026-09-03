package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxTrackedClients   = 5000
	clientInactivityTTL = 10 * time.Minute
)

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
}

// IPRateLimiter tracks per-IP token bucket rate limits.
type IPRateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientBucket
	rate     float64 // tokens per second
	capacity float64 // burst capacity
}

// NewIPRateLimiter creates an in-memory token bucket rate limiter.
// ratePerMinute specifies how many tokens are replenished per minute, and burst defines capacity.
func NewIPRateLimiter(ratePerMinute int, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		clients:  make(map[string]*clientBucket),
		rate:     float64(ratePerMinute) / 60.0,
		capacity: float64(burst),
	}
}

func (lim *IPRateLimiter) cleanupExpiredLocked(now time.Time) {
	if len(lim.clients) > maxTrackedClients {
		for ip, bucket := range lim.clients {
			if now.Sub(bucket.lastRefill) >= clientInactivityTTL {
				delete(lim.clients, ip)
			}
		}
	}
}

// Allow checks and consumes 1 token for the given IP. Returns allowed status and retry-after duration.
func (lim *IPRateLimiter) Allow(ip string) (bool, time.Duration) {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	now := time.Now()
	bucket, exists := lim.clients[ip]
	if !exists {
		lim.cleanupExpiredLocked(now)
		lim.clients[ip] = &clientBucket{
			tokens:     lim.capacity - 1,
			lastRefill: now,
		}
		return true, 0
	}

	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens = math.Min(lim.capacity, bucket.tokens+(elapsed*lim.rate))
	bucket.lastRefill = now

	if bucket.tokens < 1.0 {
		missing := 1.0 - bucket.tokens
		retryAfter := time.Duration((missing / lim.rate) * float64(time.Second))
		return false, retryAfter
	}

	bucket.tokens -= 1.0
	return true, 0
}

// RateLimit returns a Gin middleware enforcing the token bucket rate limit.
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		allowed, retryAfter := limiter.Allow(ip)
		if !allowed {
			retrySeconds := max(int(math.Ceil(retryAfter.Seconds())), 1)
			c.Header("Retry-After", strconv.Itoa(retrySeconds))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
