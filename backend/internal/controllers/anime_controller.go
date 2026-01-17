package controllers

import (
	"anime-tanyaayomi/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AnimeController struct {
	Service          *services.SankavollereiService
	AnimeIndoService *services.AnimeIndoService
}

func NewAnimeController() *AnimeController {
	return &AnimeController{
		// Keep Sankavollerei for episode streaming (existing functionality)
		Service: services.NewSankavollereiService(""),
		// Add Anime Indo for anime data (new functionality)
		AnimeIndoService: services.NewAnimeIndoService(),
	}
}

// ========== New Anime Indo Endpoints ==========

// GetLatestEpisodes returns recently released episodes
func (c *AnimeController) GetLatestEpisodes(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)

	episodes, err := c.AnimeIndoService.GetLatestEpisodes(page)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": episodes,
	})
}

// GetPopularAnime returns popular ongoing anime
func (c *AnimeController) GetPopularAnime(ctx *fiber.Ctx) error {
	animeList, err := c.AnimeIndoService.GetPopularAnime()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": animeList,
	})
}

// SearchAnimeIndo searches for anime using Anime Indo API
func (c *AnimeController) SearchAnimeIndo(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Query param 'q' is required"})
	}

	results, err := c.AnimeIndoService.SearchAnime(query)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": results,
	})
}

// GetAnimeDetailIndo returns anime details from Anime Indo API
func (c *AnimeController) GetAnimeDetailIndo(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	detail, err := c.AnimeIndoService.GetAnimeDetail(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": detail,
	})
}

// ========== Legacy Otakudesu/Sankavollerei Endpoints (Keep for backward compat) ==========

// GetHome returns ongoing and completed anime (Sankavollerei)
func (c *AnimeController) GetHome(ctx *fiber.Ctx) error {
	result, err := c.Service.GetHome()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"ongoing":   result.Data.Ongoing.AnimeList,
			"completed": result.Data.Completed.AnimeList,
		},
	})
}

// Search searches for anime by keyword (Sankavollerei)
func (c *AnimeController) Search(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Query param 'q' is required"})
	}

	result, err := c.Service.Search(query)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"animeList": result.Data.AnimeList,
		},
	})
}

// GetGenre returns anime list by genre (Sankavollerei)
func (c *AnimeController) GetGenre(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	result, err := c.Service.GetGenre(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"animeList": result.Data.AnimeList,
		},
	})
}

// GetOngoing returns ongoing anime list (Sankavollerei)
func (c *AnimeController) GetOngoing(ctx *fiber.Ctx) error {
	result, err := c.Service.GetOngoingAnime(1)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"animeList": result.Data.AnimeList,
		},
	})
}

// GetComplete returns complete anime list (Sankavollerei)
func (c *AnimeController) GetComplete(ctx *fiber.Ctx) error {
	result, err := c.Service.GetCompleteAnime(1)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"animeList": result.Data.AnimeList,
		},
	})
}

// GetDetail returns anime details (Sankavollerei)
func (c *AnimeController) GetDetail(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	result, err := c.Service.GetAnimeDetail(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": result.Data})
}

// GetStream returns episode streaming data (Sankavollerei + Anime Indo)
func (c *AnimeController) GetStream(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	// Try Sankavollerei first (handles Otakudesu slugs)
	result, err := c.Service.GetEpisodeStream(slug)
	if err == nil {
		// Success with Sankavollerei
		return ctx.JSON(fiber.Map{
			"data": fiber.Map{
				"title":               result.Data.Title,
				"defaultStreamingUrl": result.Data.DefaultStreamingUrl,
				"stream_link":         result.Data.StreamLink,
				"url":                 result.Data.URL,
				"animeId":             result.Data.AnimeID,
				"server":              result.Data.Server,
				"downloadUrl":         result.Data.DownloadURL,
			},
		})
	}

	// Fallback: Try Anime Indo API (handles clean slugs like "one-piece-episode-1155")
	animeIndoStreams, animeIndoErr := c.AnimeIndoService.GetEpisodeStream(slug, "")
	if animeIndoErr == nil && len(animeIndoStreams) > 0 {
		// Extract anime slug (remove episode number) for animeId
		animeSlug := slug
		if idx := strings.Index(slug, "-episode-"); idx != -1 {
			animeSlug = slug[:idx]
		}

		// Convert Anime Indo response to Sankavollerei format
		return ctx.JSON(fiber.Map{
			"data": fiber.Map{
				"title":               slug,
				"defaultStreamingUrl": animeIndoStreams["default"],
				"stream_link":         animeIndoStreams,
				"url":                 animeIndoStreams["default"],
				"animeId":             animeSlug, // Clean anime slug without episode
				"server": fiber.Map{
					"qualities": []fiber.Map{},
				},
				"downloadUrl": "",
			},
		})
	}

	// Both failed
	return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
}

// GetLatest returns latest episodes from Sankavollerei (Legacy)
func (c *AnimeController) GetLatest(ctx *fiber.Ctx) error {
	result, err := c.Service.GetLatestEpisodes()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"episodes": result.Data.Episodes,
		},
	})
}
