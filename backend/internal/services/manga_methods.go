package services

import (
	"anime-tanyaayomi/internal/models"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// extractSlugFromLink extracts slug from link field (e.g., "/manga/slug-name/" -> "slug-name")
func extractSlugFromLink(link string) string {
	// Remove leading/trailing slashes
	link = strings.Trim(link, "/")
	// Split by '/' and get the last part
	parts := strings.Split(link, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return link
}

// GetMangaHome fetches trending and popular manga from homepage
func (s *SankavollereiService) GetMangaHome() (*models.MangaListResponse, error) {
	var result struct {
		Trending []models.Manga `json:"trending"`
	}

	// Cache for 5 minutes
	err := s.makeRequestWithCache("comic/trending", &result, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get manga home data: %w", err)
	}

	// Extract slugs from link fields
	for i := range result.Trending {
		if result.Trending[i].Slug == "" && result.Trending[i].Link != "" {
			result.Trending[i].Slug = extractSlugFromLink(result.Trending[i].Link)
		}
	}

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Trending,
		},
	}, nil
}

// SearchManga searches for manga by keyword
func (s *SankavollereiService) SearchManga(keyword string) (*models.MangaListResponse, error) {
	var result struct {
		Data []models.Manga `json:"data"`
	}

	endpoint := fmt.Sprintf("comic/search?q=%s", url.QueryEscape(keyword))
	// Cache search results for 10 minutes
	err := s.makeRequestWithCache(endpoint, &result, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to search manga: %w", err)
	}

	// Extract slugs
	for i := range result.Data {
		if result.Data[i].Slug == "" && result.Data[i].Link != "" {
			result.Data[i].Slug = extractSlugFromLink(result.Data[i].Link)
		}
	}

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Data,
		},
	}, nil
}

// GetMangaGenre searches for manga by genre
func (s *SankavollereiService) GetMangaGenre(slug string) (*models.MangaListResponse, error) {
	var result struct {
		Comics []models.Manga `json:"comics"`
	}

	endpoint := fmt.Sprintf("comic/genre/%s", slug)
	// Cache genre results for 30 minutes
	err := s.makeRequestWithCache(endpoint, &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get genre manga: %w", err)
	}

	// Extract slugs
	for i := range result.Comics {
		if result.Comics[i].Slug == "" && result.Comics[i].Link != "" {
			result.Comics[i].Slug = extractSlugFromLink(result.Comics[i].Link)
		}
	}

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Comics,
		},
	}, nil
}

// GetOngoingManga fetches the list of recent/ongoing manga
func (s *SankavollereiService) GetOngoingManga(page int) (*models.MangaListResponse, error) {
	var result struct {
		Comics []models.Manga `json:"comics"`
	}

	// Comic API uses "terbaru" (recent) instead of "ongoing"
	endpoint := "comic/terbaru"
	if page > 1 {
		endpoint = fmt.Sprintf("comic/terbaru?page=%d", page)
	}

	// Cache ongoing results for 15 minutes
	err := s.makeRequestWithCache(endpoint, &result, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get ongoing manga: %w", err)
	}

	// Extract slugs
	for i := range result.Comics {
		if result.Comics[i].Slug == "" && result.Comics[i].Link != "" {
			result.Comics[i].Slug = extractSlugFromLink(result.Comics[i].Link)
		}
	}

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Comics,
		},
	}, nil
}

// GetCompleteManga fetches the list of popular manga
func (s *SankavollereiService) GetCompleteManga(page int) (*models.MangaListResponse, error) {
	var result struct {
		Comics []models.Manga `json:"comics"`
	}

	// Comic API uses "populer" instead of "complete"
	endpoint := "comic/populer"
	if page > 1 {
		endpoint = fmt.Sprintf("comic/populer?page=%d", page)
	}

	// Cache complete results for 60 minutes
	err := s.makeRequestWithCache(endpoint, &result, 60*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get complete manga: %w", err)
	}

	// Extract slugs
	for i := range result.Comics {
		if result.Comics[i].Slug == "" && result.Comics[i].Link != "" {
			result.Comics[i].Slug = extractSlugFromLink(result.Comics[i].Link)
		}
	}

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Comics,
		},
	}, nil
}

// GetMangaDetail fetches detailed information about a manga
func (s *SankavollereiService) GetMangaDetail(slug string) (*models.MangaDetailResponse, error) {
	var result models.MangaDetailResponse

	endpoint := fmt.Sprintf("comic/comic/%s", slug)
	// Cache manga details for 30 minutes
	err := s.makeRequestWithCache(endpoint, &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get manga detail: %w", err)
	}

	return &result, nil
}

// GetChapterImages fetches images for a manga chapter
func (s *SankavollereiService) GetChapterImages(chapterId string) (*models.ChapterResponse, error) {
	var result models.ChapterResponse

	endpoint := fmt.Sprintf("comic/chapter/%s", chapterId)
	// Cache chapter images for 30 minutes
	err := s.makeRequestWithCache(endpoint, &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapter images: %w", err)
	}

	return &result, nil
}

// GetTrendingManga fetches trending manga
func (s *SankavollereiService) GetTrendingManga() (*models.MangaListResponse, error) {
	var result struct {
		Trending []models.Manga `json:"trending"`
	}

	// Cache for 15 minutes
	err := s.makeRequestWithCache("comic/trending", &result, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending manga: %w", err)
	}

	// Extract slugs
	for i := range result.Trending {
		if result.Trending[i].Slug == "" && result.Trending[i].Link != "" {
			result.Trending[i].Slug = extractSlugFromLink(result.Trending[i].Link)
		}
	}

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Trending,
		},
	}, nil
}
