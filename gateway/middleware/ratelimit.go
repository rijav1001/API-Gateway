package middleware

import (
	"net/http"
	"sync"
	"time"
)

type client struct {
	tokens   float64
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	rate     float64 // tokens per second
	maxTokens float64
}

func NewRateLimiter(rate float64, max float64) *RateLimiter {
	rl := &RateLimiter{
		clients:   make(map[string]*client),
		rate:      rate,
		maxTokens: max,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &client{tokens: rl.maxTokens, lastSeen: time.Now()}
		return true
	}

	elapsed := time.Since(c.lastSeen).Seconds()
	c.tokens += elapsed * rl.rate
	if c.tokens > rl.maxTokens {
		c.tokens = rl.maxTokens
	}
	c.lastSeen = time.Now()

	if c.tokens < 1 {
		return false
	}
	c.tokens--
	return true
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rl.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > 10*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}