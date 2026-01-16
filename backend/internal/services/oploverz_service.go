package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// OploverzService handles communication with Oploverz API (sankavollerei.com)
type OploverzService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// OploverzDownload represents a single download option
type OploverzDownload struct {
	Name       string `json:"name"`       // e.g., "GD", "FILEDON", "AKIRA"
	Resolution string `json:"resolution"` // e.g., "360p", "480p", "720p", "1080p"
	URL        string `json:"url"`        // Direct download/stream URL
}

// OploverzStream represents a streaming option
type OploverzStream struct {
	Name string `json:"name"` // e.g., "Main Stream", "google v2"
	URL  string `json:"url"`  // Embed URL (usually Blogger iframe)
}

// OploverzEpisodeResponse represents the ACTUAL API response from Oploverz
// Fields are at ROOT level, not nested in "data"!
type OploverzEpisodeResponse struct {
	Status       string             `json:"status"`        // "success" or "error" (STRING, not boolean!)
	Creator      string             `json:"creator"`       // "Sanka Vollerei"
	Source       string             `json:"source"`        // "Oploverz"
	EpisodeTitle string             `json:"episode_title"` // Full episode title
	Streams      []OploverzStream   `json:"streams"`       // Streaming embeds
	Downloads    []OploverzDownload `json:"downloads"`     // Quality/download options
}

// NewOploverzService creates a new Oploverz service client
func NewOploverzService() *OploverzService {
	return &OploverzService{
		BaseURL: "https://www.sankavollerei.com",
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetEpisodeStream fetches streaming URLs for an anime episode from Oploverz API
func (s *OploverzService) GetEpisodeStream(otakudesuSlug string) (map[string]string, error) {
	// Convert Otakudesu slug to Oploverz format
	oploverzSlug := s.convertSlugToOploverz(otakudesuSlug)

	// Build API URL
	apiURL := fmt.Sprintf("%s/anime/oploverz/episode/%s", s.BaseURL, oploverzSlug)
	fmt.Printf("[Oploverz] Fetching: %s\n", apiURL)

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Make HTTP request
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Oploverz episode: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Oploverz API returned status %d", resp.StatusCode)
	}

	// Parse JSON response
	var apiResp OploverzEpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Oploverz response: %w", err)
	}

	// Check if API call was successful
	if apiResp.Status != "success" {
		// Log what we got
		return nil, fmt.Errorf("Oploverz API error: status='%s' (struct=%+v)", apiResp.Status, apiResp)
	}

	// Build quality map from downloads
	// Group by resolution and use first URL for each quality
	qualityMap := make(map[string]string)

	// Process downloads - group by resolution
	for _, download := range apiResp.Downloads {
		resolution := strings.ToLower(strings.TrimSpace(download.Resolution))

		// Only store first URL for each resolution (avoid duplicates)
		if _, exists := qualityMap[resolution]; !exists {
			// Convert download URL to player/embed URL for better playback
			embedURL := s.convertToEmbedURL(download.URL)
			qualityMap[resolution] = embedURL
		}
	}

	// Add streams as additional options
	// Add streams as additional options - REMOVED AS REQUESTED
	// for i, stream := range apiResp.Streams {
	// 	// Use Name if available, otherwise fallback to generic
	// 	streamKey := stream.Name
	// 	if streamKey == "" {
	// 		streamKey = fmt.Sprintf("Stream %d", i+1)
	// 	}
	// 	qualityMap[streamKey] = stream.URL
	// }

	// Use first stream as default if available (still useful for internal default, but won't be in map)
	// Actually, if we don't add them to map, they won't appear.
	// For default key:
	if len(apiResp.Streams) > 0 {
		qualityMap["default"] = apiResp.Streams[0].URL
	} else if len(apiResp.Downloads) > 0 {
		// Fallback to first download if no streams
		qualityMap["default"] = apiResp.Downloads[0].URL
	}

	fmt.Printf("[Oploverz] ✅ Found %d quality options for %s\n", len(qualityMap), oploverzSlug)
	return qualityMap, nil
}

// convertToEmbedURL converts download URLs to embed/player URLs for direct playback
// Example: https://acefile.co/f/110881619/... -> https://acefile.co/player/110881619
func (s *OploverzService) convertToEmbedURL(downloadURL string) string {
	// Handle acefile.co URLs - convert from download to player
	if strings.Contains(downloadURL, "acefile.co/f/") {
		// Extract file ID using regex: acefile.co/f/110881619/filename
		re := regexp.MustCompile(`acefile\.co/f/(\d+)`)
		matches := re.FindStringSubmatch(downloadURL)
		if len(matches) > 1 {
			fileID := matches[1]
			playerURL := fmt.Sprintf("https://acefile.co/player/%s", fileID)
			fmt.Printf("[Oploverz] 🔄 Converted acefile URL: %s -> %s\n", downloadURL, playerURL)
			return playerURL
		}
	}

	// Handle filedon.co URLs - already player-ready, just ensure proper format
	if strings.Contains(downloadURL, "filedon.co/view/") {
		return downloadURL // These are already embed URLs
	}

	// Handle akirabox.to URLs - convert to embed format if needed
	if strings.Contains(downloadURL, "akirabox.to/") {
		// Extract file ID from URL
		re := regexp.MustCompile(`akirabox\.to/([^/]+)`)
		matches := re.FindStringSubmatch(downloadURL)
		if len(matches) > 1 && !strings.Contains(downloadURL, "/embed/") {
			fileID := matches[1]
			embedURL := fmt.Sprintf("https://akirabox.to/embed/%s", fileID)
			fmt.Printf("[Oploverz] 🔄 Converted akirabox URL: %s -> %s\n", downloadURL, embedURL)
			return embedURL
		}
	}

	// For buzzheavier and other URLs, return as-is (may need iframe handling)
	return downloadURL
}

// convertSlugToOploverz converts Otakudesu slug format to Oploverz format
func (s *OploverzService) convertSlugToOploverz(otakudesuSlug string) string {
	// Extract anime title and episode number using regex
	re := regexp.MustCompile(`^(.+?)-episode-(\d+)`)
	matches := re.FindStringSubmatch(otakudesuSlug)

	if len(matches) != 3 {
		// Fallback: return as-is if parsing fails
		return otakudesuSlug
	}

	animeSlug := matches[1]  // e.g., "wpoiec"
	episodeNum := matches[2] // e.g., "1155"

	// Map common Otakudesu abbreviations to full titles
	fullTitle := s.mapSlugToFullTitle(animeSlug)

	// Build Oploverz slug format
	oploverzSlug := fmt.Sprintf("%s-episode-%s-subtitle-indonesia", fullTitle, episodeNum)
	return oploverzSlug
}

// mapSlugToFullTitle maps common Otakudesu anime slugs to full readable titles for Oploverz
func (s *OploverzService) mapSlugToFullTitle(slug string) string {
	// Common anime title mappings
	titleMap := map[string]string{
		"wpoiec":          "one-piece",
		"bkunhro":         "boku-no-hero-academia",
		"stvssn":          "spy-x-family",
		"kslym":           "kimetsu-no-yaiba",
		"jjksn":           "jujutsu-kaisen",
		"atkslyr":         "attack-on-titan",
		"nruto":           "naruto",
		"nrtsppdn":        "naruto-shippuden",
		"blach":           "bleach",
		"dmnslyar":        "demon-slayer",
		"tokyo-revengers": "tokyo-revengers",
		"blue-lock":       "blue-lock",
		"windbreaker":     "wind-breaker",
		"mushoku-tensei":  "mushoku-tensei",
		"solo-leveling":   "solo-leveling",
		"overlord":        "overlord",
		"re-zero":         "re-zero",
		"konosuba":        "konosuba",
		"danmachi":        "danmachi",
		"tensura":         "tensei-shitara-slime-datta-ken",
	}

	if fullTitle, exists := titleMap[slug]; exists {
		return fullTitle
	}

	// If no mapping found, return slug as-is (might work for some anime)
	return slug
}
