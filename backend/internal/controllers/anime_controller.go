package controllers

import (
	"anime-tanyaayomi/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AnimeController struct {
	Service *services.SankavollereiService
}

func NewAnimeController() *AnimeController {
	return &AnimeController{
		// Default to Otakudesu (empty prefix)
		Service: services.NewSankavollereiService(""),
	}
}

// GetHome returns ongoing and completed anime
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

// Search searches for anime by keyword
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

// GetGenre returns anime list by genre
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

// GetOngoing returns ongoing anime list
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

// GetComplete returns complete anime list (used for Popular section)
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

// GetDetail returns anime details
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

// GetStream returns episode streaming data
func (c *AnimeController) GetStream(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	result, err := c.Service.GetEpisodeStream(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Format response to match frontend expectations
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

// GetLatest returns latest episodes from all sources
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
