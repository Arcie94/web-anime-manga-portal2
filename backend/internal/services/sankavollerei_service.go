package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// Rate limiter to respect API limits (70 req/min max)
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	mu         sync.Mutex
	lastRefill time.Time
}

func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on time passed
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)

	if tokensToAdd > 0 {
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	// Check if we have tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// Cache entry
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Simple in-memory cache
type Cache struct {
	store map[string]*CacheEntry
	mu    sync.RWMutex
}

func NewCache() *Cache {
	cache := &Cache{
		store: make(map[string]*CacheEntry),
	}
	// Start cleanup goroutine
	go cache.cleanup()
	return cache
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.store[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.store {
			if now.After(entry.ExpiresAt) {
				delete(c.store, key)
			}
		}
		c.mu.Unlock()
	}
}

// Sankavollerei Service
type SankavollereiService struct {
	BaseURL     string
	Client      *http.Client
	Prefix      string
	RateLimiter *RateLimiter
	Cache       *Cache
}

func NewSankavollereiService(prefix string) *SankavollereiService {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Explicitly set proxy if available
	if proxyEnv := os.Getenv("HTTP_PROXY"); proxyEnv != "" {
		proxyURL, err := url.Parse(proxyEnv)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	return &SankavollereiService{
		BaseURL: "https://www.sankavollerei.com",
		Client:  client,
		Prefix:  prefix,
		// 70 requests per minute = 1 token every ~857ms
		RateLimiter: NewRateLimiter(70, 857*time.Millisecond),
		Cache:       NewCache(),
	}
}

func (s *SankavollereiService) makeRequest(endpoint string, result interface{}) error {
	// Check rate limit
	if !s.RateLimiter.Allow() {
		return fmt.Errorf("rate limit exceeded, please wait")
	}

	// Build URL - check if endpoint starts with "comic/"
	var url string
	if len(endpoint) >= 6 && endpoint[:6] == "comic/" {
		// For manga/comic endpoints, use the endpoint directly from BaseURL
		// Example: "comic/home" -> "https://www.sankavollerei.com/comic/home"
		url = fmt.Sprintf("%s/%s", s.BaseURL, endpoint)
	} else {
		// For anime endpoints, add /anime/ prefix and service prefix if exists
		// Example: "home" -> "https://www.sankavollerei.com/anime/home"
		url = fmt.Sprintf("%s/anime/%s%s", s.BaseURL, s.Prefix, endpoint)
	}

	// Make request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("fake request failed: %w", err)
	}

	// Add User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

func (s *SankavollereiService) makeRequestWithCache(endpoint string, result interface{}, cacheTTL time.Duration) error {
	// Check cache first
	cacheKey := s.Prefix + endpoint
	if cached, found := s.Cache.Get(cacheKey); found {
		// Use type assertion to copy cached data to result
		cachedJSON, _ := json.Marshal(cached)
		json.Unmarshal(cachedJSON, result)
		return nil
	}

	// Make actual API request
	if err := s.makeRequest(endpoint, result); err != nil {
		return err
	}

	// Store in cache
	s.Cache.Set(cacheKey, result, cacheTTL)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
