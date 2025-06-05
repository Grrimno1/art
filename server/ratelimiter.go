package server

import (
	"log"
	"net/http"
	"sync"
	"time"
	"net"
)

// RateLimiter tracks incoming requests per IP and enforces rate limiting to prevent abuse.
type RateLimiter struct {
	mu		sync.Mutex
	clients	map[string][]time.Time // maps ip timestamps of requests
	limit	int
	interval time.Duration
}

// NewRateLimiter initializes a RateLimiter with a given request limit and interval.
func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter {
		clients:	make(map[string][]time.Time),
		limit:		limit,
		interval:	interval,
	}
}

// allow checks whether the given IP is permitted to make a new request.
// It filters timestamps to keep only recent ones within the allowed time window,
// and denies access if the request count exceeds the limit.
func (limiter *RateLimiter) allow(ip string) bool {
	limiter.mu.Lock() // protect clients map access and modification
	defer limiter.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-limiter.interval)

	times := limiter.clients[ip]
	var updatedTimes []time.Time

	// filter out old request timestamps outside the current interval window
	for _, t := range times {
		if t.After(windowStart) {
			updatedTimes = append(updatedTimes, t)
		}
	}
	// if number of recent requests exceeds limit, deny
	if len(updatedTimes) >= limiter.limit {
		limiter.clients[ip] = updatedTimes
		return false
	}

	// record the current request time and allow
	limiter.clients[ip] = append(updatedTimes, now)
	return true
}

// Middleware wraps an existing http.Handler and applies IP-based rate limiting
// before allowing the request to proceed. If the rate limit is exceeded, it returns
// an appropriate error response depending on the request method.
func (limiter *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr) // identify user by IP and remove :port
		if err != nil {
			ip = r.RemoteAddr
		}

		if !limiter.allow(ip) {
			if r.Method == http.MethodPost {
				// Render error page via template for POST requests.
				data := CombinedPageData{
					Section:        "art",
					DecodeInput:    "",
					EncodeInput:    "",
					StatusCode:     http.StatusTooManyRequests,
					StatusType:     "error",
					StatusMessage:  "429 too many requests: Please wait a few seconds and try again.",
					LineCount:      4,
				}
				if err := tmpl.Execute(w, data); err != nil {
   	 				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    				log.Printf("template execute error: %v", err)
				}
			} else {
				// Return plain HTTP error for non-POST requests.
				http.Error(w, "429 too many requests", http.StatusTooManyRequests)
			}
			return
		}

		// Proceed to the next handler if rate limit not exceeded.
		next.ServeHTTP(w, r)
	})
}