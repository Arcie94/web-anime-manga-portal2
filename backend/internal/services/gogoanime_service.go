package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ZoroService handles requests to Zoro.to via Consumet API
// Using Zoro instead of Gogoanime due to better API reliability
type ZoroService struct {
	BaseURL string
	Client  *http.Client
}

// NewZoroService creates a new Zoro service instance
func NewZoroService() *ZoroService {
	return &ZoroService{
		BaseURL: "https://api.consumet.org", // Official Consumet API
		Client: &http.Client{
			Timeout: 15 * time.Second, // Slightly longer timeout for Zoro
		},
	}
}

// ZoroStreamSource represents a single stream quality option
type ZoroStreamSource struct {
	Quality string `json:"quality"` // "1080p", "720p", "480p", "360p", "default"
	URL     string `json:"url"`
	IsM3U8  bool   `json:"isM3U8"`
}

// ZoroStreamResponse represents the response from Zoro watch endpoint
type ZoroStreamResponse struct {
	Headers struct {
		Referer string `json:"Referer"`
	} `json:"headers"`
	Sources []ZoroStreamSource `json:"sources"`
}

// ZoroEpisode represents an episode in the anime
type ZoroEpisode struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// ZoroAnimeInfo represents anime search result
type ZoroAnimeInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Image string `json:"image"`
	URL   string `json:"url"`
}

// ZoroSearchResponse represents Zoro search results
type ZoroSearchResponse struct {
	CurrentPage int             `json:"currentPage"`
	HasNextPage bool            `json:"hasNextPage"`
	Results     []ZoroAnimeInfo `json:"results"`
}

// GetEpisodeStream fetches streaming URLs for a specific anime episode
// Parameters:
//   - animeTitle: Title of the anime (e.g., "One Piece")
//   - episodeNumber: Episode number (e.g., 1155)
//
// Returns: Map of quality → URL
func (s *ZoroService) GetEpisodeStream(animeTitle string, episodeNumber int) (map[string]string, error) {
	// Step 1: Search for anime on Zoro to get anime ID
	searchQuery := strings.ToLower(strings.ReplaceAll(animeTitle, " ", "-"))
	searchURL := fmt.Sprintf("%s/anime/zoro/%s?page=1", s.BaseURL, searchQuery)

	searchResp, err := s.Client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search Zoro: %w", err)
	}
	defer searchResp.Body.Close()

	if searchResp.StatusCode != 200 {
		return nil, fmt.Errorf("zoro search returned status %d", searchResp.StatusCode)
	}

	searchBody, _ := io.ReadAll(searchResp.Body)
	var searchResults ZoroSearchResponse
	if err := json.Unmarshal(searchBody, &searchResults); err != nil {
		return nil, fmt.Errorf("failed to parse Zoro search: %w", err)
	}

	if len(searchResults.Results) == 0 {
		return nil, fmt.Errorf("anime not found on Zoro")
	}

	// Get the first result (usually most relevant)
	animeID := searchResults.Results[0].ID

	// Step 2: Get anime info to find episode ID
	infoURL := fmt.Sprintf("%s/anime/zoro/info?id=%s", s.BaseURL, animeID)
	infoResp, err := s.Client.Get(infoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get Zoro info: %w", err)
	}
	defer infoResp.Body.Close()

	if infoResp.StatusCode != 200 {
		return nil, fmt.Errorf("zoro info returned status %d", infoResp.StatusCode)
	}

	infoBody, _ := io.ReadAll(infoResp.Body)

	// Parse to get episodes list
	var animeInfo struct {
		Episodes []ZoroEpisode `json:"episodes"`
	}
	if err := json.Unmarshal(infoBody, &animeInfo); err != nil {
		return nil, fmt.Errorf("failed to parse Zoro info: %w", err)
	}

	// Find the specific episode
	var episodeID string
	for _, ep := range animeInfo.Episodes {
		if ep.Number == episodeNumber {
			episodeID = ep.ID
			break
		}
	}

	if episodeID == "" {
		return nil, fmt.Errorf("episode %d not found on Zoro", episodeNumber)
	}

	// Step 3: Get stream URLs for the episode
	watchURL := fmt.Sprintf("%s/anime/zoro/watch?episodeId=%s", s.BaseURL, episodeID)
	watchResp, err := s.Client.Get(watchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Zoro stream: %w", err)
	}
	defer watchResp.Body.Close()

	if watchResp.StatusCode != 200 {
		return nil, fmt.Errorf("zoro watch returned status %d", watchResp.StatusCode)
	}

	watchBody, _ := io.ReadAll(watchResp.Body)
	var streamData ZoroStreamResponse
	if err := json.Unmarshal(watchBody, &streamData); err != nil {
		return nil, fmt.Errorf("failed to parse Zoro stream: %w", err)
	}

	// Step 4: Extract quality → URL mapping
	qualityMap := make(map[string]string)
	for _, source := range streamData.Sources {
		if source.URL != "" {
			qualityMap[source.Quality] = source.URL
		}
	}

	if len(qualityMap) == 0 {
		return nil, fmt.Errorf("no stream sources found")
	}

	return qualityMap, nil
}

// ParseOtakudesuSlug extracts anime title and episode number from Otakudesu slug
// Example: "wpoiec-episode-1155-sub-indo" → ("One Piece", 1155)
func ParseOtakudesuSlug(slug string) (string, int, error) {
	// Remove "-sub-indo" suffix
	slug = strings.TrimSuffix(slug, "-sub-indo")

	// Extract episode number using regex
	re := regexp.MustCompile(`episode-(\d+)`)
	matches := re.FindStringSubmatch(slug)
	if len(matches) < 2 {
		return "", 0, fmt.Errorf("could not extract episode number from slug: %s", slug)
	}

	episodeNumber := 0
	fmt.Sscanf(matches[1], "%d", &episodeNumber)

	// Extract anime title slug (before "-episode-")
	parts := strings.Split(slug, "-episode-")
	if len(parts) < 1 {
		return "", 0, fmt.Errorf("invalid slug format: %s", slug)
	}

	titleSlug := parts[0]

	// Map known anime title slugs to proper names
	animeTitle := mapSlugToTitle(titleSlug)

	return animeTitle, episodeNumber, nil
}

// mapSlugToTitle converts Otakudesu slug to proper anime title
// Expandable mapping for popular anime
func mapSlugToTitle(slug string) string {
	knownTitles := map[string]string{
		"wpoiec":           "One Piece",
		"boruto":           "Boruto: Naruto Next Generations",
		"bleach":           "Bleach",
		"naruto":           "Naruto",
		"naruto-shippuden": "Naruto: Shippuuden",
		"black-clover":     "Black Clover",
		"demon-slayer":     "Kimetsu no Yaiba",
		"jujutsu-kaisen":   "Jujutsu Kaisen",
		"haikyuu":          "Haikyuu!!",
		// Add more as needed
	}

	if title, exists := knownTitles[slug]; exists {
		return title
	}

	// Fallback: Convert slug to title case
	// "jujutsu-kaisen" → "Jujutsu Kaisen"
	words := strings.Split(slug, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
