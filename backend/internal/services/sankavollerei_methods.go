package services

import (
	"anime-tanyaayomi/internal/models"
	"fmt"
	"time"
)

// GetHome fetches ongoing and completed anime from homepage
func (s *SankavollereiService) GetHome() (*models.HomeResponse, error) {
	var result models.HomeResponse

	// Cache for 5 minutes
	err := s.makeRequestWithCache("home", &result, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get home data: %w", err)
	}

	return &result, nil
}

// Search searches for anime by keyword
func (s *SankavollereiService) Search(keyword string) (*models.SearchResponse, error) {
	var result models.SearchResponse

	endpoint := fmt.Sprintf("search/%s", keyword)
	// Cache search results for 10 minutes
	err := s.makeRequestWithCache(endpoint, &result, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to search anime: %w", err)
	}

	return &result, nil
}

// GetGenre searches for anime by genre
func (s *SankavollereiService) GetGenre(slug string) (*models.SearchResponse, error) {
	var result models.SearchResponse

	endpoint := fmt.Sprintf("genre/%s", slug)
	// Cache genre results for 30 minutes
	err := s.makeRequestWithCache(endpoint, &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get genre anime: %w", err)
	}

	return &result, nil
}

// GetOngoingAnime fetches the list of ongoing anime from specific page
func (s *SankavollereiService) GetOngoingAnime(page int) (*models.SearchResponse, error) {
	var result models.SearchResponse

	endpoint := fmt.Sprintf("ongoing-anime/page/%d", page)
	if page <= 1 {
		endpoint = "ongoing-anime"
	}

	// Cache ongoing results for 15 minutes
	err := s.makeRequestWithCache(endpoint, &result, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get ongoing anime: %w", err)
	}

	return &result, nil
}

// GetCompleteAnime fetches the list of completed anime from specific page
func (s *SankavollereiService) GetCompleteAnime(page int) (*models.SearchResponse, error) {
	var result models.SearchResponse

	endpoint := fmt.Sprintf("complete-anime/page/%d", page)
	if page <= 1 {
		endpoint = "complete-anime"
	}

	// Cache complete results for 60 minutes
	err := s.makeRequestWithCache(endpoint, &result, 60*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get complete anime: %w", err)
	}

	return &result, nil
}

// GetAnimeDetail fetches detailed information about an anime
func (s *SankavollereiService) GetAnimeDetail(slug string) (*models.AnimeDetailResponse, error) {
	var result models.AnimeDetailResponse

	endpoint := fmt.Sprintf("anime/%s", slug)
	// Cache anime details for 30 minutes (they don't change often)
	err := s.makeRequestWithCache(endpoint, &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get anime detail: %w", err)
	}

	return &result, nil
}

// GetEpisodeStream fetches streaming URLs for an episode
func (s *SankavollereiService) GetEpisodeStream(episodeId string) (*models.StreamResponse, error) {
	var result models.StreamResponse

	endpoint := fmt.Sprintf("episode/%s", episodeId)
	// Cache episode streams for 15 minutes
	err := s.makeRequestWithCache(endpoint, &result, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get episode stream: %w", err)
	}

	return &result, nil
}

// GetServerURL fetches specific server embed URL
func (s *SankavollereiService) GetServerURL(serverId string) (string, error) {
	var result struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	endpoint := fmt.Sprintf("server/%s", serverId)
	// Cache server URLs for 20 minutes
	err := s.makeRequestWithCache(endpoint, &result, 20*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to get server URL: %w", err)
	}

	return result.Data.URL, nil
}

// GetLatest fetches latest episodes from all sources
func (s *SankavollereiService) GetLatestEpisodes() (*models.LatestResponse, error) {
	var result models.LatestResponse

	// Cache for 3 minutes (latest updates change frequently)
	err := s.makeRequestWithCache("stream/latest", &result, 3*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest episodes: %w", err)
	}

	return &result, nil
}
