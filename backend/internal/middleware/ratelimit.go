package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count    int
	lastSeen time.Time
}

var (
	rateLimits   = make(map[string]*rateLimitEntry)
	rateLimitsMu sync.Mutex
)

const (
	loginWindow    = time.Minute
	loginMaxReqs   = 10
	apiWindow      = time.Minute
	apiMaxReqs     = 60
	cleanupEvery   = 5 * time.Minute
)

func init() {
	go cleanupLoop()
}

func cleanupLoop() {
	for {
		time.Sleep(cleanupEvery)
		rateLimitsMu.Lock()
		now := time.Now()
		for ip, entry := range rateLimits {
			if now.Sub(entry.lastSeen) > loginWindow*2 {
				delete(rateLimits, ip)
			}
		}
		rateLimitsMu.Unlock()
	}
}

// LoginRateLimit limits login attempts per IP
func LoginRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rateLimitsMu.Lock()
		entry, exists := rateLimits[ip]
		now := time.Now()
		if !exists || now.Sub(entry.lastSeen) > loginWindow {
			rateLimits[ip] = &rateLimitEntry{count: 1, lastSeen: now}
			rateLimitsMu.Unlock()
			next(w, r)
			return
		}
		entry.count++
		entry.lastSeen = now
		currentCount := entry.count
		rateLimitsMu.Unlock()

		if currentCount > loginMaxReqs {
			http.Error(w, `{"error":"Too many login attempts. Please try again later."}`, http.StatusTooManyRequests)
			w.Header().Set("Content-Type", "application/json")
			return
		}
		next(w, r)
	}
}
