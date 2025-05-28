package server

import (
	"net/http"
	"sync"
	"time"
)

// tracks requests per IP and enforces rate limits.
type RateLimiter struct {
	mu		sync.Mutex
	clients	map[string][]time.Time // maps ip timestamps of requests
	limit	int
	interval time.Duration
}

// initializes a RateLimiter with a given request limit and interval.
func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter {
		clients:	make(map[string][]time.Time),
		limit:		limit,
		interval:	interval,
	}
}

// allow checks whether the IP is allowed to proceed based on the rate limit
func (r1 *RateLimiter) allow(ip string) bool {
	r1.mu.Lock()
	defer r1.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r1.interval)

	times := r1.clients[ip]
	var updatedTimes []time.Time

	// filter out old request timestamps outside the current interval window
	for _, t := range times {
		if t.After(windowStart) {
			updatedTimes = append(updatedTimes, t)
		}
	}
	// if number of recent requests exceeds limit, deny
	if len(updatedTimes) >= r1.limit {
		r1.clients[ip] = updatedTimes
		return false
	}

	// record the current request time and allow
	r1.clients[ip] = append(updatedTimes, now)
	return true
}

// MiddleWare wraps an http.Handler and applies rate limiting to requests.
func (r1 *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // identifies client by their IP address.

		// Check if this request is allowed
		if !r1.allow(ip) {
			if r.Method == http.MethodPost {
				// shows error msg through html template via POST
				data := CombinedPageData{
					Section:        "art",
					DecodeInput:    "",
					EncodeInput:    "",
					StatusCode:     http.StatusTooManyRequests,
					StatusType:     "error",
					StatusMessage:  "429 too many requests: Please wait a few seconds and try again.",
					LineCount:      4,
				}
				tmpl.Execute(w, data) // render error via frontend
			} else {
				// If not POST return plain text HTTP error
				http.Error(w, "429 too many requests", http.StatusTooManyRequests)
			}
			return
		}

		// If allowed, pass request to the next handler.
		next.ServeHTTP(w, r)
	})
}