package services

import (
	"anime-tanyaayomi/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const MangaBaseURL = "https://www.sankavollerei.com/comic"

type MangaService struct {
	Client *RateLimitedClient
}

func NewMangaService() *MangaService {
	return &MangaService{
		Client: NewRateLimitedClient(),
	}
}

func (s *MangaService) GetJSON(endpoint string, target interface{}) error {
	url := fmt.Sprintf("%s%s", MangaBaseURL, endpoint)
	// Use the helper Get method from RateLimitedClient which handles the request creation internally
	// BUT RateLimitedClient.Get takes just URL. We need headers.
	// So let's use Do.

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (s *MangaService) SearchManga(query string) ([]models.Manga, error) {
	// Replacing space with + or %20 might be needed, but usually http client handles it if encoded properly,
	// here we just use string formatting, so better replace spaces.
	// Python script used /search?q={keyword}
	query = strings.ReplaceAll(query, " ", "+")
	endpoint := fmt.Sprintf("/search?q=%s", query)

	var response struct {
		Data []models.Manga `json:"data"`
	}

	if err := s.GetJSON(endpoint, &response); err != nil {
		return nil, err
	}

	for i := range response.Data {
		if response.Data[i].Cover == "" && response.Data[i].Thumbnail != "" {
			response.Data[i].Cover = response.Data[i].Thumbnail
		}
	}

	return response.Data, nil
}

func (s *MangaService) GetMangaDetail(slug string) (*models.MangaDetail, error) {
	endpoint := fmt.Sprintf("/comic/%s", slug)
	var response struct {
		Data models.MangaDetail `json:"data"`
	}

	if err := s.GetJSON(endpoint, &response); err != nil {
		return nil, err
	}

	// Normalize chapters if needed (handle both ChapterList and Chapters fields)
	if len(response.Data.ChapterList) == 0 && len(response.Data.Chapters) > 0 {
		response.Data.ChapterList = response.Data.Chapters
	}

	// Normalize covers
	if response.Data.Cover == "" && response.Data.Thumbnail != "" {
		response.Data.Cover = response.Data.Thumbnail
	}

	return &response.Data, nil
}

func (s *MangaService) GetChapterImages(slug string) ([]string, error) {
	endpoint := fmt.Sprintf("/chapter/%s", slug)
	var response struct {
		Images []string `json:"images"`
	}

	if err := s.GetJSON(endpoint, &response); err != nil {
		return nil, err
	}
	return response.Images, nil
}
