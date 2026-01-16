package services

import (
	"anime-tanyaayomi/internal/database"
	"anime-tanyaayomi/internal/models"
	"anime-tanyaayomi/internal/repository"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Package-level repository for AI enrichment (initialized on first use)
var (
	enrichmentRepo     *repository.EnrichmentRepository
	enrichmentRepoOnce sync.Once
)

// getEnrichmentRepo returns the singleton enrichment repository
func getEnrichmentRepo() *repository.EnrichmentRepository {
	enrichmentRepoOnce.Do(func() {
		if database.DB != nil {
			enrichmentRepo = repository.NewEnrichmentRepository(database.DB)
			fmt.Println("[EnrichmentRepo] ✅ Initialized with database connection")
		} else {
			fmt.Println("[EnrichmentRepo] ⚠️  Database not available, enrichment will use cache-only mode")
		}
	})
	return enrichmentRepo
}

// enrichAnimeWithGemini uses Gemini AI to fill missing data (Year, Rating, etc.) for Anime
// Now with database persistence to reduce API costs by ~90%
func (s *SankavollereiService) enrichAnimeWithGemini(items []models.Anime) []models.Anime {
	var wg sync.WaitGroup
	maxConcurrency := 5
	sem := make(chan struct{}, maxConcurrency)

	// Get repository (will use database if available)
	repo := getEnrichmentRepo()

	for i := range items {
		// Only enrich if missing data (upstream usually sends "ReleaseDate" as empty or non-year string)
		// Usually upstream has specific format, but if it's missing or we want year:
		// We'll trust Gemini if ReleaseDate is missing.
		if items[i].ReleaseDate == "" {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				// Use database-aware enrichment (DB → Cache → Gemini)
				data := EnrichDataWithDB(items[idx].Title, "anime", repo)
				if data.Year != "" {
					items[idx].ReleaseDate = data.Year
				}
				if items[idx].Status == "" && data.Status != "" {
					items[idx].Status = data.Status
				}
				if items[idx].Rating == "" && data.Rating != "" {
					items[idx].Rating = data.Rating
				}
			}(i)
		}

		// User requested Type to be "Anime" always, instead of "TV"/"Movie" or empty.
		// This applies to all lists (Home, Ongoing, Completed, Search, Genre)
		items[i].Type = "Anime"
	}
	wg.Wait()
	return items
}

// GetHome fetches ongoing and completed anime from homepage
func (s *SankavollereiService) GetHome() (*models.HomeResponse, error) {
	// Check cache first (cache the enriched data, not raw API response)
	cacheKey := s.Prefix + "home_enriched"
	if cached, found := s.Cache.Get(cacheKey); found {
		if result, ok := cached.(*models.HomeResponse); ok {
			return result, nil
		}
	}

	// Fetch raw data from API (no cache for raw data)
	var result models.HomeResponse
	err := s.makeRequest("home", &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get home data: %w", err)
	}

	// AI Enrichment (this now happens BEFORE caching)
	result.Data.Ongoing.AnimeList = s.enrichAnimeWithGemini(result.Data.Ongoing.AnimeList)
	result.Data.Completed.AnimeList = s.enrichAnimeWithGemini(result.Data.Completed.AnimeList)

	// Cache the enriched result for 5 minutes
	s.Cache.Set(cacheKey, &result, 5*time.Minute)

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

	// Optional: Enrich search results? Might be slow. Let's skip for speed for now or do it?
	// User didn't complain about search. Let's stick to lists.

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

	// AI Enrichment
	result.Data.AnimeList = s.enrichAnimeWithGemini(result.Data.AnimeList)

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

	// AI Enrichment
	result.Data.AnimeList = s.enrichAnimeWithGemini(result.Data.AnimeList)

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

	// AI Enrichment
	result.Data.AnimeList = s.enrichAnimeWithGemini(result.Data.AnimeList)

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

	// AI Enrichment (Gemini)
	// Check if key data is missing. Otakudesu usually returns empty strings for Genre/Rating/Author/Type.
	if result.Data.Author == "" || result.Data.Genre == "" {
		func() {
			// Get repository (will use database if available)
			repo := getEnrichmentRepo()

			// Run enrichment with database support (DB → Cache → Gemini)
			enriched := EnrichDataWithDB(result.Data.Title, "anime", repo)

			if enriched.Year != "" && result.Data.ReleaseDate == "" {
				result.Data.ReleaseDate = enriched.Year
			}
			if enriched.Author != "" && result.Data.Author == "" {
				result.Data.Author = enriched.Author
			}
			if enriched.Genre != "" && result.Data.Genre == "" {
				result.Data.Genre = enriched.Genre
			}
			if enriched.Rating != "" && result.Data.Rating == "" {
				result.Data.Rating = enriched.Rating
			}
			if enriched.Status != "" && result.Data.Status == "" {
				result.Data.Status = enriched.Status
			}
		}()
	}

	// User requested "Type" to be "Anime" always, instead of "TV"/"Movie" or empty.
	result.Data.Type = "Anime"

	return &result, nil
}

// GetEpisodeStream fetches streaming URLs for an episode from multiple sources
// Now aggregates both Otakudesu (Server 1) and Oploverz (Quality Selector)
func (s *SankavollereiService) GetEpisodeStream(episodeId string) (*models.StreamResponse, error) {
	var result models.StreamResponse

	// Step 1: Get Otakudesu stream (Server 1 - Sub Indo, default quality)
	endpoint := fmt.Sprintf("episode/%s", episodeId)
	err := s.makeRequestWithCache(endpoint, &result, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get episode stream: %w", err)
	}

	// Step 2: Try to get Oploverz streams for quality selector
	oploverzService := NewOploverzService()
	oploverzStreams, oploverzErr := oploverzService.GetEpisodeStream(episodeId)

	if oploverzErr == nil && len(oploverzStreams) > 0 {
		// Successfully got Oploverz streams, add as quality options
		fmt.Printf("[Stream] ✅ Oploverz qualities available for %s\n", episodeId)

		// Build quality options from Oploverz response
		qualityOptions := []models.QualityOption{}

		// Map of quality order for sorting
		qualityOrder := map[string]int{
			"360p":  1,
			"480p":  2,
			"720p":  3,
			"1080p": 4,
		}

		// Group streams by quality
		for quality, url := range oploverzStreams {
			if quality == "default" {
				continue // Skip default, we'll use specific qualities
			}

			// Create server entry for this quality
			// Oploverz typically provides multiple mirrors per quality
			serverList := []models.StreamServer{
				{
					Title:    "Oploverz",
					ServerID: fmt.Sprintf("oploverz_%s", quality),
					Href:     url,
				},
			}

			qualityOptions = append(qualityOptions, models.QualityOption{
				Title:      quality,
				ServerList: serverList,
			})
		}

		// Sort qualities by resolution (360p first, then 480p, 720p, 1080p)
		sort.Slice(qualityOptions, func(i, j int) bool {
			orderI := qualityOrder[qualityOptions[i].Title]
			orderJ := qualityOrder[qualityOptions[j].Title]
			return orderI < orderJ
		})

		// Update result with quality options
		if len(qualityOptions) > 0 {
			result.Data.Server.Qualities = qualityOptions
			fmt.Printf("[Stream] ✅ Added %d quality options from Oploverz\n", len(qualityOptions))
		}
	} else {
		fmt.Printf("[Stream] ⚠️ Oploverz unavailable for %s: %v\n", episodeId, oploverzErr)
		// Continue without Oploverz - Otakudesu stream still available
	}

	// Step 3: Update Server 1 title to indicate Otakudesu
	// Rename existing qualities to "Server 1 - Otakudesu"
	// This step is only relevant if Oploverz did NOT overwrite the qualities.
	// If Oploverz successfully provided streams, result.Data.Server.Qualities would have been replaced.
	// If Oploverz failed, the original Otakudesu qualities remain, and this update applies.
	if oploverzErr != nil && len(result.Data.Server.Qualities) > 0 {
		for i := range result.Data.Server.Qualities {
			// Only update the first quality (Otakudesu)
			if i == 0 {
				result.Data.Server.Qualities[i].Title = fmt.Sprintf("Server 1 - Otakudesu (Sub Indo) - %s", result.Data.Server.Qualities[i].Title)
				break
			}
		}
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
