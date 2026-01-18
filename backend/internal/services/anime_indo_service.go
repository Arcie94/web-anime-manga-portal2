package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// AnimeIndoService handles communication with Anime Indo Stream API (sankavollerei.com)
type AnimeIndoService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// ========== Response Structs ==========

// LatestEpisode represents a recent episode release
type LatestEpisode struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Poster  string `json:"poster"`
	Episode string `json:"episode"`
}

// AnimeCard represents an anime in list view (popular, search results)
type AnimeCard struct {
	Title  string   `json:"title"`
	Slug   string   `json:"slug"`
	Poster string   `json:"poster"`
	Genres []string `json:"genres"`
}

// Episode represents a single episode in anime detail
type Episode struct {
	EpisodeTitle string `json:"eps_title"`
	EpisodeSlug  string `json:"eps_slug"`
}

// AnimeDetail represents full anime information
type AnimeDetail struct {
	Title    string    `json:"title"`
	Poster   string    `json:"poster"`
	Synopsis string    `json:"synopsis"`
	Genres   []string  `json:"genres"`
	Episodes []Episode `json:"episodes"`
}

// StreamLink represents a streaming server option
type StreamLink struct {
	Server string `json:"server"`
	URL    string `json:"url"`
}

// DownloadLink represents a download option
type DownloadLink struct {
	Server string `json:"server"`
	URL    string `json:"url"`
}

// EpisodeStreamData represents episode stream response data
type EpisodeStreamData struct {
	Title         string         `json:"title"`
	StreamLinks   []StreamLink   `json:"stream_links"`
	DownloadLinks []DownloadLink `json:"download_links"`
}

// ========== API Response Wrappers ==========

type LatestResponse struct {
	Status  int             `json:"status"`
	Creator string          `json:"creator"`
	Page    int             `json:"page"`
	Data    []LatestEpisode `json:"data"`
}

type PopularResponse struct {
	Status  int         `json:"status"`
	Creator string      `json:"creator"`
	Data    []AnimeCard `json:"data"`
}

type SearchResponse struct {
	Status  int         `json:"status"`
	Creator string      `json:"creator"`
	Query   string      `json:"query"`
	Data    []AnimeCard `json:"data"`
}

type AnimeDetailResponse struct {
	Status  int         `json:"status"`
	Creator string      `json:"creator"`
	Data    AnimeDetail `json:"data"`
}

type EpisodeStreamResponse struct {
	Status int               `json:"status"`
	Data   EpisodeStreamData `json:"data"`
}

// ========== Constructor ==========

// NewAnimeIndoService creates a new Anime Indo service client
func NewAnimeIndoService() *AnimeIndoService {
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

	return &AnimeIndoService{
		BaseURL:    "https://www.sankavollerei.com",
		HTTPClient: client,
	}
}

// ========== Public API Methods ==========

// GetLatestEpisodes fetches recently released episodes
func (s *AnimeIndoService) GetLatestEpisodes(page int) ([]LatestEpisode, error) {
	endpoint := fmt.Sprintf("%s/anime/stream/latest?page=%d", s.BaseURL, page)

	var resp LatestResponse
	if err := s.doRequest(endpoint, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetPopularAnime fetches popular ongoing anime
func (s *AnimeIndoService) GetPopularAnime() ([]AnimeCard, error) {
	endpoint := fmt.Sprintf("%s/anime/stream/popular", s.BaseURL)

	var resp PopularResponse
	if err := s.doRequest(endpoint, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// SearchAnime searches for anime by query
func (s *AnimeIndoService) SearchAnime(query string) ([]AnimeCard, error) {
	endpoint := fmt.Sprintf("%s/anime/stream/search/%s", s.BaseURL, strings.ReplaceAll(query, " ", "%20"))

	var resp SearchResponse
	if err := s.doRequest(endpoint, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GetAnimeDetail fetches full anime details including episode list
func (s *AnimeIndoService) GetAnimeDetail(slug string) (*AnimeDetail, error) {
	endpoint := fmt.Sprintf("%s/anime/stream/anime/%s", s.BaseURL, slug)

	fmt.Printf("[AnimeIndo] Fetching detail for slug: %s\n", slug)
	fmt.Printf("[AnimeIndo] Endpoint: %s\n", endpoint)

	var resp AnimeDetailResponse
	if err := s.doRequest(endpoint, &resp); err != nil {
		fmt.Printf("[AnimeIndo] Error fetching detail: %v\n", err)
		return nil, err
	}

	fmt.Printf("[AnimeIndo] Got detail: %s with %d episodes\n", resp.Data.Title, len(resp.Data.Episodes))

	return &resp.Data, nil
}

// GetEpisodeStream fetches episode streams (with Otakudesu slug compatibility)
func (s *AnimeIndoService) GetEpisodeStream(otakudesuSlug string, otakudesuTitle string) (map[string]string, error) {
	// Candidate 1: Raw Slug (as provided by frontend/Otakudesu/Anime Indo)
	// This is important because Anime Indo sometimes puts "-tamat" or "-end" in the slug
	slugCandidates := []string{otakudesuSlug}

	// Candidate 2: Cleaned/Converted Slug
	// This handles "wpoiec-episode-123" -> "one-piece-episode-123"
	cleanedSlug := s.convertSlugToAnimeIndo(otakudesuSlug, otakudesuTitle)
	if cleanedSlug != otakudesuSlug {
		slugCandidates = append(slugCandidates, cleanedSlug)
	}

	// Candidate 3: Cleaned + "-sub-indo"
	// Many Anime Indo slugs utilize this suffix if the raw one fails
	if !strings.HasSuffix(cleanedSlug, "-sub-indo") {
		slugCandidates = append(slugCandidates, cleanedSlug+"-sub-indo")
	}

	// Candidate 4: Raw + "-sub-indo"
	if !strings.HasSuffix(otakudesuSlug, "-sub-indo") {
		slugCandidates = append(slugCandidates, otakudesuSlug+"-sub-indo")
	}

	var resp EpisodeStreamResponse
	var lastErr error

	for _, slug := range slugCandidates {
		endpoint := fmt.Sprintf("%s/anime/stream/episode/%s", s.BaseURL, slug)
		fmt.Printf("[AnimeIndo] Trying stream slug: '%s'\n", slug)

		lastErr = s.doRequest(endpoint, &resp)
		if lastErr == nil && len(resp.Data.StreamLinks) > 0 {
			fmt.Printf("[AnimeIndo] Success with slug: '%s'\n", slug)
			// Break on first success
			break
		}
	}

	if lastErr != nil || len(resp.Data.StreamLinks) == 0 {
		return nil, fmt.Errorf("failed to get stream after trying candidates: %v", lastErr)
	}

	// Convert to map format for frontend
	streamMap := make(map[string]string)
	for _, link := range resp.Data.StreamLinks {
		// Trim any whitespace or special characters from URL
		cleanURL := strings.TrimSpace(link.URL)
		fmt.Printf("[AnimeIndo] Found Stream: %s -> %s\n", link.Server, cleanURL)
		streamMap[link.Server] = cleanURL
	}

	// Set default to B-TUBE if available (most reliable), otherwise first stream
	if streamMap["B-TUBE"] != "" {
		streamMap["default"] = streamMap["B-TUBE"]
	} else if len(resp.Data.StreamLinks) > 0 {
		streamMap["default"] = resp.Data.StreamLinks[0].URL
	}

	return streamMap, nil
}

// ========== Helper Methods ==========

// doRequest performs HTTP GET and decodes JSON response
func (s *AnimeIndoService) doRequest(endpoint string, v interface{}) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")

	client := s.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: status=%s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// convertSlugToAnimeIndo converts Otakudesu format to Anime Indo format
func (s *AnimeIndoService) convertSlugToAnimeIndo(otakudesuSlug string, otakudesuTitle string) string {
	// Example: "wpoiec-episode-1155-sub-indo" -> "one-piece-episode-1155"

	re := regexp.MustCompile(`^(.*?)-episode-(\d+)(?:-.*)?$`)
	matches := re.FindStringSubmatch(otakudesuSlug)

	if len(matches) < 3 {
		return otakudesuSlug
	}

	animeSlug := matches[1]  // "wpoiec"
	episodeNum := matches[2] // "1155"

	// Map weird slugs to real titles
	mappedTitle := s.mapSlugToFullTitle(animeSlug)

	// Construct new slug: {title}-episode-{num}
	return fmt.Sprintf("%s-episode-%s", mappedTitle, episodeNum)
}

// mapSlugToFullTitle maps common Otakudesu anime slugs to readable titles
func (s *AnimeIndoService) mapSlugToFullTitle(slug string) string {
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

	return slug
}
