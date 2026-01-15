package services

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/time/rate"
)

// Global Rate Limiter
// Limit: 60 requests per minute (1 request every 1s) to be safe (max 70/min)
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 1)

type RateLimitedClient struct {
	Client *http.Client
}

func NewRateLimitedClient() *RateLimitedClient {
	client := &http.Client{
		Timeout: 20 * time.Second, // Increased timeout for proxy calls
	}

	// Check if HTTP_PROXY is set
	if proxyEnv := os.Getenv("HTTP_PROXY"); proxyEnv != "" {
		proxyURL, err := url.Parse(proxyEnv)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	return &RateLimitedClient{
		Client: client,
	}
}

func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
	// Wait for rate limiter permission
	err := limiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *RateLimitedClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Add User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	return c.Do(req)
}
