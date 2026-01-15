package services

import (
	"anime-tanyaayomi/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://www.sankavollerei.com"

type AnimeService struct {
	Client *RateLimitedClient
}

func NewAnimeService() *AnimeService {
	return &AnimeService{
		Client: NewRateLimitedClient(),
	}
}

func (s *AnimeService) GetJSON(endpoint string, target interface{}) error {
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)
	resp, err := s.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[AnimeService] Error fetching %s: Status %d\n", url, resp.StatusCode)
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// log.Printf("[AnimeService] Response: %s", string(body))
	fmt.Printf("[AnimeService] Response Body: %s\n", string(body))

	if err := json.Unmarshal(body, target); err != nil {
		fmt.Printf("[AnimeService] Unmarshal Error: %v\n", err)
		return err
	}
	return nil
}

func (s *AnimeService) SearchAnime(query string) ([]models.Anime, error) {
	endpoint := fmt.Sprintf("/anime/search/%s", query)
	var response struct {
		Data struct {
			AnimeList []models.Anime `json:"animeList"`
		} `json:"data"`
	}

	if err := s.GetJSON(endpoint, &response); err != nil {
		return nil, err
	}

	for i := range response.Data.AnimeList {
		anime := &response.Data.AnimeList[i]
		if anime.Cover == "" {
			if anime.Poster != "" {
				anime.Cover = anime.Poster
			} else if anime.Thumbnail != "" {
				anime.Cover = anime.Thumbnail
			} else if anime.Image != "" {
				anime.Cover = anime.Image
			}
		}
		// Fix missing Slug
		if anime.Slug == "" && anime.AnimeID != "" {
			anime.Slug = anime.AnimeID
		}
	}

	return response.Data.AnimeList, nil
}

func (s *AnimeService) GetAnimeDetail(slug string) (*models.AnimeDetail, error) {
	endpoint := fmt.Sprintf("/anime/anime/%s", slug)
	var response struct {
		Data models.AnimeDetail `json:"data"`
	}

	if err := s.GetJSON(endpoint, &response); err != nil {
		return nil, err
	}

	if response.Data.Cover == "" {
		if response.Data.Poster != "" {
			response.Data.Cover = response.Data.Poster
		} else if response.Data.Thumbnail != "" {
			response.Data.Cover = response.Data.Thumbnail
		} else if response.Data.Image != "" {
			response.Data.Cover = response.Data.Image
		}
	}

	return &response.Data, nil
}

func (s *AnimeService) GetEpisodeStream(slug string) (*models.StreamData, error) {
	endpoint := fmt.Sprintf("/anime/episode/%s", slug)
	var response struct {
		Data models.StreamData `json:"data"`
	}

	if err := s.GetJSON(endpoint, &response); err != nil {
		return nil, err
	}
	return &response.Data, nil
}
