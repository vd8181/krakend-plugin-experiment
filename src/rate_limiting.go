package krakend_ratelimiter

import (

	"context"

	"errors"

	"net/http"

	"sync"

	"time"

	"github.com/devopsfaith/krakend/transport/http/server"

)

var (

	globalLimiter *RateLimiter

	endpointLimiters = make(map[string]*RateLimiter)

	mu sync.Mutex

)

// RateLimiter implements a token bucket algorithm

type RateLimiter struct {

	capacity   int

	tokens     int

	refillRate int // tokens per second

	lastRefill time.Time

	mu         sync.Mutex

}

func NewRateLimiter(capacity, refillRate int) *RateLimiter {

	return &RateLimiter{

		capacity:   capacity,

		tokens:     capacity,

		refillRate: refillRate,

		lastRefill: time.Now(),

	}

}

func (r *RateLimiter) Allow() bool {

	r.mu.Lock()

	defer r.mu.Unlock()

	now := time.Now()

	elapsed := now.Sub(r.lastRefill).Seconds()

	newTokens := int(elapsed * float64(r.refillRate))

	if newTokens > 0 {

		r.tokens = min(r.capacity, r.tokens+newTokens)

		r.lastRefill = now

	}

	if r.tokens > 0 {

		r.tokens--

		return true

	}

	return false

}

func min(a, b int) int {

	if a < b {

		return a

	}

	return b

}

// Register plugin

func init() {

	server.RegisterMiddleware("custom-ratelimit", rateLimitMiddleware)

}

func rateLimitMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		endpoint := r.URL.Path

		limiter := getLimiterForEndpoint(endpoint)

		if !limiter.Allow() {

			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)

			return

		}

		next.ServeHTTP(w, r)

	})

}

func getLimiterForEndpoint(endpoint string) *RateLimiter {

	mu.Lock()

	defer mu.Unlock()

	if limiter, ok := endpointLimiters[endpoint]; ok {

		return limiter

	}

	return globalLimiter

}

// Plugin Configuration Interface

type Config struct {

	Global struct {

		Capacity   int `json:"capacity"`

		RefillRate int `json:"refill_rate"`

	} `json:"global"`

	Endpoints map[string]struct {

		Capacity   int `json:"capacity"`

		RefillRate int `json:"refill_rate"`

	} `json:"endpoints"`

}

// Called by KrakenD at plugin initialization

func New(cfg map[string]interface{}) (func(http.Handler) http.Handler, error) {

	parsed := Config{}

	// Decode configuration into struct

	if err := decode(cfg, &parsed); err != nil {

		return nil, err

	}

	// Set global limiter

	globalLimiter = NewRateLimiter(parsed.Global.Capacity, parsed.Global.RefillRate)

	// Set per-endpoint limiters

	for ep, conf := range parsed.Endpoints {

		endpointLimiters[ep] = NewRateLimiter(conf.Capacity, conf.RefillRate)

	}

	return rateLimitMiddleware, nil

}

// Decode helper (use your preferred decoder)

func decode(src map[string]interface{}, dst *Config) error {

	// Use mapstructure, json.Unmarshal with conversion, or a manual decode like below

	// For simplicity, you can assume the structure is always valid

	return errors.New("use a proper decoding method like mapstructure or json marshal")

}
