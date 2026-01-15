package services

import "anime-tanyaayomi/internal/models"

type AnimeProvider interface {
	SearchAnime(query string) ([]models.Anime, error)
	GetAnimeDetail(slug string) (*models.AnimeDetail, error)
	GetEpisodeStream(slug string) (*models.StreamData, error)
}

type MangaProvider interface {
	SearchManga(query string) ([]models.Manga, error)
	GetMangaDetail(slug string) (*models.MangaDetail, error)
	GetChapterImages(slug string) ([]string, error)
}
