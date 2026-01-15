package services

import (
	"anime-tanyaayomi/internal/models"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync" // Added sync
	"time"
)

// Helper to remove resize/quality params from image URL
func cleanImageURL(imgUrl string) string {
	if imgUrl == "" {
		return ""
	}
	// Regex to remove ?resize=... or &resize=... and ?quality=... or &quality=...
	reResize := regexp.MustCompile(`[?&]resize=[^&]+`)
	reQuality := regexp.MustCompile(`[?&]quality=[^&]+`)

	clean := reResize.ReplaceAllString(imgUrl, "")
	clean = reQuality.ReplaceAllString(clean, "")
	return clean
}

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

// enrichMangaListURLs fetches detailed info for each manga in parallel to get the high-quality portrait image
func (s *SankavollereiService) enrichMangaListURLs(items []models.Manga) []models.Manga {
	var wg sync.WaitGroup
	// Limit concurrency to avoid overwhelming the upstream server or local resources
	maxConcurrency := 10
	sem := make(chan struct{}, maxConcurrency)

	for i := range items {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Ensure slug exists
			if items[idx].Slug == "" && items[idx].Link != "" {
				items[idx].Slug = extractSlugFromLink(items[idx].Link)
			}
			if items[idx].Slug == "" {
				return
			}

			// Fetch detail (utilizes existing cache on GetMangaDetail)
			detail, err := s.GetMangaDetail(items[idx].Slug)
			if err == nil && detail.Image != "" {
				// Replace low-res/landscape image with high-quality portrait from detail
				// GetMangaDetail already cleans the URL
				items[idx].Image = detail.Image
				items[idx].Cover = detail.Image
				items[idx].Poster = detail.Image
				items[idx].Thumbnail = detail.Image
			} else {
				// Fallback: just clean the existing URL
				items[idx].Image = cleanImageURL(items[idx].Image)
				items[idx].Cover = cleanImageURL(items[idx].Cover)
				items[idx].Poster = cleanImageURL(items[idx].Poster)
				items[idx].Thumbnail = cleanImageURL(items[idx].Thumbnail)
			}
		}(i)
	}
	wg.Wait()
	return items
}

// filterBlacklistedManga removes unwanted items (e.g. "APK") from the list
func (s *SankavollereiService) filterBlacklistedManga(items []models.Manga) []models.Manga {
	var filtered []models.Manga
	for _, item := range items {
		title := strings.ToLower(item.Title)
		if strings.Contains(title, "apk") || strings.Contains(title, "komiku plus") {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

// enrichMangaWithGemini uses Gemini AI to fill missing data (Year, Rating, etc.)
func (s *SankavollereiService) enrichMangaWithGemini(items []models.Manga) []models.Manga {
	var wg sync.WaitGroup
	// Limit Concurrency for AI to 5 to be safe with rate limits
	maxConcurrency := 5
	sem := make(chan struct{}, maxConcurrency)

	for i := range items {
		// Only enrich if missing data
		if items[i].ReleaseDate == "" {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				data := EnrichData(items[idx].Title, "manga")
				if data.Year != "" {
					items[idx].ReleaseDate = data.Year
				}
				if items[idx].Status == "" && data.Status != "" {
					items[idx].Status = data.Status
				}
				// We can also add rating if we add the field to Manga struct, but let's stick to ReleaseDate for now
			}(i)
		}
	}
	wg.Wait()
	return items
}

// GetMangaHome fetches trending and popular manga from homepage
func (s *SankavollereiService) GetMangaHome() (*models.MangaListResponse, error) {
	var result struct {
		Trending []models.Manga `json:"trending"`
	}

	// Cache for 30 minutes (Increased from 5 due to heavy enrichment)
	err := s.makeRequestWithCache("comic/trending", &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get manga home data: %w", err)
	}

	// Filter unwanted items
	result.Trending = s.filterBlacklistedManga(result.Trending)

	// Enrich images in parallel
	result.Trending = s.enrichMangaListURLs(result.Trending)

	// AI Enrichment (Gemini)
	result.Trending = s.enrichMangaWithGemini(result.Trending)

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

	// Filter unwanted items
	result.Data = s.filterBlacklistedManga(result.Data)

	// For search, we might not want to enrich ALL results (slow), just clean them.
	// Or we can enrich, but let's stick to cleaning for search responsiveness.
	for i := range result.Data {
		if result.Data[i].Slug == "" && result.Data[i].Link != "" {
			result.Data[i].Slug = extractSlugFromLink(result.Data[i].Link)
		}
		result.Data[i].Cover = cleanImageURL(result.Data[i].Cover)
		result.Data[i].Poster = cleanImageURL(result.Data[i].Poster)
		result.Data[i].Thumbnail = cleanImageURL(result.Data[i].Thumbnail)
		result.Data[i].Image = cleanImageURL(result.Data[i].Image)
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

	// Filter unwanted items
	result.Comics = s.filterBlacklistedManga(result.Comics)

	// Enrich images in parallel
	result.Comics = s.enrichMangaListURLs(result.Comics)

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

	// Filter unwanted items
	result.Comics = s.filterBlacklistedManga(result.Comics)

	// Enrich images in parallel
	result.Comics = s.enrichMangaListURLs(result.Comics)

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

	// Filter unwanted items
	result.Comics = s.filterBlacklistedManga(result.Comics)

	// Enrich images in parallel
	result.Comics = s.enrichMangaListURLs(result.Comics)

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

	result.Image = cleanImageURL(result.Image)

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

	// Logic to find Next/Prev Slug
	// 1. Extract Manga Slug from Chapter Slug directly
	re := regexp.MustCompile(`^(.*)-chapter-.*$`)
	matches := re.FindStringSubmatch(chapterId)
	if len(matches) > 1 {
		mangaSlug := matches[1]

		if result.MangaID == "" {
			result.MangaID = mangaSlug
		}
		if result.ChapterID == "" {
			result.ChapterID = chapterId
		}

		// 2. Fetch Manga Details
		detail, err := s.GetMangaDetail(mangaSlug)
		if err != nil {
			// Fallback 1: Try common prefix/suffix pattern "komik-{slug}-indo"
			// This avoids Search latency/timeout for common cases like One Piece
			heuristicSlug := fmt.Sprintf("komik-%s-indo", mangaSlug)
			detail, err = s.GetMangaDetail(heuristicSlug)
			if err != nil {
				// Fallback 2: Search for the slug to find the correct canonical slug
				// e.g. "one-piece" -> "One Piece" -> Search -> "komik-one-piece-indo"
				searchQuery := strings.ReplaceAll(mangaSlug, "-", " ")
				searchRes, searchErr := s.SearchManga(searchQuery)
				if searchErr == nil && searchRes != nil && len(searchRes.Data.MangaList) > 0 {
					var canonicalSlug string

					// Smart Selection: Find exact title match first
					for _, m := range searchRes.Data.MangaList {
						if strings.EqualFold(m.Title, searchQuery) {
							canonicalSlug = m.Slug
							break
						}
					}

					// If no exact match, fallback to first result
					if canonicalSlug == "" {
						canonicalSlug = searchRes.Data.MangaList[0].Slug
					}

					// Retry Detail Fetch
					detail, err = s.GetMangaDetail(canonicalSlug)
				}
			}
		}

		if detail != nil {
			reNum := regexp.MustCompile(`-chapter-([\d\.]+)`)
			numMatches := reNum.FindStringSubmatch(chapterId)
			var targetNum string
			if len(numMatches) > 1 {
				targetNum = numMatches[1]
			}

			for i, ch := range detail.Chapters {
				iterSlug := strings.Trim(ch.Slug, "/")
				cleanID := strings.Trim(chapterId, "/")

				isMatch := false
				if iterSlug == cleanID {
					isMatch = true
				} else if targetNum != "" {
					if strings.HasSuffix(iterSlug, "-"+targetNum) || strings.Contains(iterSlug, "-chapter-"+targetNum) {
						isMatch = true
					}
				}

				if isMatch {
					if i > 0 {
						result.NextSlug = detail.Chapters[i-1].Slug
					}
					if i < len(detail.Chapters)-1 {
						result.PrevSlug = detail.Chapters[i+1].Slug
					}
					break
				}
			}
		}
	}

	return &result, nil
}

// GetTrendingManga fetches trending manga
func (s *SankavollereiService) GetTrendingManga() (*models.MangaListResponse, error) {
	var result struct {
		Trending []models.Manga `json:"trending"`
	}

	// Cache for 30 minutes (Increased due to enrichment)
	err := s.makeRequestWithCache("comic/trending", &result, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending manga: %w", err)
	}

	// Filter unwanted items
	result.Trending = s.filterBlacklistedManga(result.Trending)

	// Enrich images in parallel
	result.Trending = s.enrichMangaListURLs(result.Trending)

	// AI Enrichment
	result.Trending = s.enrichMangaWithGemini(result.Trending)

	return &models.MangaListResponse{
		Data: struct {
			MangaList []models.Manga `json:"mangaList"`
		}{
			MangaList: result.Trending,
		},
	}, nil
}
