package controllers

import (
	"anime-tanyaayomi/internal/services"

	"github.com/gofiber/fiber/v2"
)

type MangaController struct {
	Service *services.SankavollereiService
}

func NewMangaController() *MangaController {
	return &MangaController{
		Service: services.NewSankavollereiService(""),
	}
}

// GetHome fetches manga home page data (trending + ongoing)
func (c *MangaController) GetHome(ctx *fiber.Ctx) error {
	result, err := c.Service.GetMangaHome()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}

// Search searches for manga by keyword
func (c *MangaController) Search(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Query param 'q' is required"})
	}

	result, err := c.Service.SearchManga(query)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}

// GetGenre fetches manga by genre
func (c *MangaController) GetGenre(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	result, err := c.Service.GetMangaGenre(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}

// GetTrending fetches trending manga
func (c *MangaController) GetTrending(ctx *fiber.Ctx) error {
	result, err := c.Service.GetTrendingManga()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}

// GetOngoing fetches ongoing manga
func (c *MangaController) GetOngoing(ctx *fiber.Ctx) error {
	result, err := c.Service.GetOngoingManga(1)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}

// GetDetail fetches manga detail
func (c *MangaController) GetDetail(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	result, err := c.Service.GetMangaDetail(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}

// GetChapter fetches chapter images
func (c *MangaController) GetChapter(ctx *fiber.Ctx) error {
	chapterId := ctx.Params("chapterId")
	if chapterId == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Chapter ID is required"})
	}

	result, err := c.Service.GetChapterImages(chapterId)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(result)
}
