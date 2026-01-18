package controllers

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ProxyController struct {
	Client *http.Client
}

func NewProxyController() *ProxyController {
	// Create a client that respects HTTP_PROXY environment variable
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	// Helper to check if a proxy env is set (optional, for logging)
	if os.Getenv("HTTP_PROXY") != "" {
		// fmt.Println("ProxyController: Using HTTP_PROXY")
	}

	return &ProxyController{
		Client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// GetImage proxies an image URL through the backend
func (c *ProxyController) GetImage(ctx *fiber.Ctx) error {
	imageURL := ctx.Query("url")
	if imageURL == "" {
		return ctx.Status(400).SendString("Missing url parameter")
	}

	// Decode URL if it's double encoded or just to be safe
	if decoded, err := url.QueryUnescape(imageURL); err == nil {
		imageURL = decoded
	}

	// Create request
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return ctx.Status(500).SendString("Failed to create request: " + err.Error())
	}

	// Set headers to mimic a real browser to bypass hotlink protection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	// Some sites check Referer. We can try setting it to the domain of the image or generic.
	// For now, let's parse the host from the URL and use it as referer
	if parsedURL, err := url.Parse(imageURL); err == nil {
		req.Header.Set("Referer", parsedURL.Scheme+"://"+parsedURL.Host+"/")
	}

	// Execute request
	resp, err := c.Client.Do(req)
	if err != nil {
		return ctx.Status(502).SendString("Failed to fetch image: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ctx.Status(resp.StatusCode).SendString("Upstream server returned error")
	}

	// Copy Content-Type header
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// Fallback detection
		if strings.HasSuffix(imageURL, ".jpg") || strings.HasSuffix(imageURL, ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(imageURL, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(imageURL, ".webp") {
			contentType = "image/webp"
		}
	}
	ctx.Set("Content-Type", contentType)

	// Stream the body directly to the client
	// Fiber's Context.Response().BodyWriter() might be needed for streaming,
	// but SetBodyStreamWriter expects func(*bufio.Writer)

	// Using io.Copy to Fiber's writer
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		io.Copy(w, resp.Body)
		w.Flush()
	})

	return nil
}
