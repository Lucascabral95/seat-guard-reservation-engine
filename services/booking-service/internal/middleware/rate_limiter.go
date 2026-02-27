package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
	rate    rate.Limit
	burst   int

	cleanupInterval time.Duration
	inactiveTTL     time.Duration
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return newRateLimiterWithConfig(r, b, 5*time.Minute, 10*time.Minute)
}

func newRateLimiterWithConfig(r rate.Limit, b int, cleanupInterval, inactiveTTL time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*client),
		rate:            r,
		burst:           b,
		cleanupInterval: cleanupInterval,
		inactiveTTL:     inactiveTTL,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) getClient(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if c, exists := rl.clients[ip]; exists {
		c.lastSeen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.clients[ip] = &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > rl.inactiveTTL {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getClient(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Demasiadas solicitudes, intenta de nuevo mas tarde",
			})
			return
		}
		c.Next()
	}
}
