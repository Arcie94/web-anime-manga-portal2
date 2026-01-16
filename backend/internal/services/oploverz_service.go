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

// GetEpisodeStream fetches episode streams from Oploverz API
// Returns map of Quality -> URL (e.g., "1080p": "https://...")
func (s *OploverzService) GetEpisodeStream(otakudesuSlug string, otakudesuTitle string) (map[string]string, error) {
	// 1. Convert Otakudesu slug/title to Oploverz slug
	oploverzSlug := s.convertSlugToOploverz(otakudesuSlug, otakudesuTitle)

	fmt.Printf("[Oploverz] Converted '%s' -> '%s'\n", otakudesuSlug, oploverzSlug)

	// 2. Fetch data from Oploverz
	// URL: https://www.sankavollerei.com/anime/oploverz/episode/{slug}
	endpoint := fmt.Sprintf("%s/anime/oploverz/episode/%s", s.BaseURL, oploverzSlug)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Add User-Agent to mimic browser (critical for Oploverz/Cloudflare)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Oploverz API error: status=%s", resp.Status)
	}

	// 3. Parse Response
	var apiResp OploverzEpisodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode Oploverz response: %w", err)
	}

	// 4. Extract Qualities
	qualityMap := make(map[string]string)

	// Helper to add quality if not exists
	addQuality := func(resolution, url string) {
		if url == "" {
			return
		}

		// If explicit iframe/embed URL, use it
		// Otherwise, convert Acefile download URLs to embeds
		finalURL := s.convertToEmbedURL(url)

		// Normalize resolution key (remove 'Op ' prefix etc if needed)
		qualityMap[resolution] = finalURL
	}

	// Iterate through Downloads to find streaming links
	for _, dl := range apiResp.Downloads {
		// Use Name to check provider (e.g. "Acefile", "Akira", "FileDon")
		label := strings.ToLower(dl.Name)

		// If it's a provider we like for streaming
		if strings.Contains(label, "acefile") || strings.Contains(label, "filedon") || strings.Contains(label, "akirabox") {
			addQuality(dl.Resolution, dl.URL)
		}
	}

	// Use first stream as default if available (still useful for internal default, but won't be in map)
	// Actually, if we don't add them to map, they won't appear.
	// For default key:
	if len(apiResp.Streams) > 0 {
		qualityMap["default"] = apiResp.Streams[0].URL
	} else if len(apiResp.Downloads) > 0 {
		// Fallback to first download if no streams
		qualityMap["default"] = apiResp.Downloads[0].URL
	}

	return qualityMap, nil
}

// convertToEmbedURL transforms download URLs to player embed URLs
func (s *OploverzService) convertToEmbedURL(url string) string {
	// Acefile: https://acefile.co/f/110881619 -> https://acefile.co/player/110881619
	if strings.Contains(url, "acefile.co/f/") {
		return strings.Replace(url, "/f/", "/player/", 1)
	}
	// Akirabox: https://akirabox.to/f/xxxxx -> https://akirabox.to/embed/xxxxx
	// (Note: verify akirabox format, assuming /embed/ based on typical patterns, or just leave as is if unsure)

	// Filedon: https://filedon.co/f/xxxxx -> https://filedon.co/view/xxxxx (hypothetical, need verification)
	// Actually filedon usually works with /v/ or /view/ or just /f/ might be download only.
	// For now, let's treat acefile as the primary target.

	return url
}

// convertSlugToOploverz attempts to convert Otakudesu slug/title to Oploverz format
func (s *OploverzService) convertSlugToOploverz(otakudesuSlug string, otakudesuTitle string) string {
	// Regex to extract anime slug and episode number
	// Otakudesu format: {anime-slug}-episode-{num}-subtitle-indonesia
	// Example: "wpoiec-episode-1155-sub-indo" -> anime: "wpoiec", num: "1155"

	re := regexp.MustCompile(`^(.*?)-episode-(\d+)(?:-.*)?$`)
	matches := re.FindStringSubmatch(otakudesuSlug)

	if len(matches) < 3 {
		// Fallback: return as-is if parsing fails
		return otakudesuSlug
	}

	animeSlug := matches[1]  // e.g., "wpoiec"
	episodeNum := matches[2] // e.g., "1155"

	// 1. Try Map lookup first (Manual override for weird slugs like 'wpoiec')
	mappedTitle := s.mapSlugToFullTitle(animeSlug)
	if mappedTitle != "" {
		return fmt.Sprintf("%s-episode-%s-subtitle-indonesia", mappedTitle, episodeNum)
	}

	// 2. Dynamic Fallback: Use Title if available
	// Otakudesu Title: "One Piece Episode 1155 Subtitle Indonesia"
	if otakudesuTitle != "" {
		// Clean title: Remove "Episode ...", "Subtitle Indonesia"
		cleanTitle := strings.ToLower(otakudesuTitle)

		// Remove "subtitle indonesia" or "sub indo"
		cleanTitle = strings.ReplaceAll(cleanTitle, "subtitle indonesia", "")
		cleanTitle = strings.ReplaceAll(cleanTitle, "sub indo", "")

		// Remove Episode/Ep/Eps + Number
		reEp := regexp.MustCompile(`(episode|ep|eps)\s*\d+.*`)
		cleanTitle = reEp.ReplaceAllString(cleanTitle, "")

		cleanTitle = strings.TrimSpace(cleanTitle)

		// Replace non-alphanumeric with hyphens
		reg, _ := regexp.Compile(`[^a-z0-9]+`)
		dynamicSlug := reg.ReplaceAllString(cleanTitle, "-")
		dynamicSlug = strings.Trim(dynamicSlug, "-")

		fmt.Printf("[Oploverz] Dynamic slug generated: '%s' from title '%s'\n", dynamicSlug, otakudesuTitle)
		return fmt.Sprintf("%s-episode-%s-subtitle-indonesia", dynamicSlug, episodeNum)
	}

	// 3. Last Resort: Use the Otakudesu anime slug as-is (unlikely to work for weird ones)
	return fmt.Sprintf("%s-episode-%s-subtitle-indonesia", animeSlug, episodeNum)
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
		"kaijuu-8-gou":    "kaijuu-8-gou",
		"dandadan":        "dandadan",
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
