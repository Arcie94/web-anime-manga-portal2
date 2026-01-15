package controllers

import (
	"anime-tanyaayomi/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MangaController struct {
	Service *services.MangaService
	Cache   *services.CacheService
}

func NewMangaController() *MangaController {
	return &MangaController{
		Service: services.NewMangaService(),
		Cache:   services.NewCacheService(),
	}
}

func (c *MangaController) Search(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Query param 'q' is required"})
	}

	cacheKey := "manga:search:" + query
	var cachedResults []interface{}

	if c.Cache.Get(cacheKey, &cachedResults) {
		return ctx.JSON(fiber.Map{"data": cachedResults, "source": "cache"})
	}

	results, err := c.Service.SearchManga(query)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cache.Set(cacheKey, results, 30*time.Minute)

	return ctx.JSON(fiber.Map{"data": results, "source": "upstream"})
}

func (c *MangaController) GetDetail(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	cacheKey := "manga:detail:" + slug
	var cachedDetail interface{}

	if c.Cache.Get(cacheKey, &cachedDetail) {
		return ctx.JSON(fiber.Map{"data": cachedDetail, "source": "cache"})
	}

	detail, err := c.Service.GetMangaDetail(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cache.Set(cacheKey, detail, 60*time.Minute)

	return ctx.JSON(fiber.Map{"data": detail, "source": "upstream"})
}

func (c *MangaController) GetChapter(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(400).JSON(fiber.Map{"error": "Slug is required"})
	}

	cacheKey := "manga:chapter:" + slug
	var cachedImages interface{}

	if c.Cache.Get(cacheKey, &cachedImages) {
		return ctx.JSON(fiber.Map{"data": cachedImages, "source": "cache"})
	}

	images, err := c.Service.GetChapterImages(slug)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cache.Set(cacheKey, images, 15*time.Minute)

	return ctx.JSON(fiber.Map{"data": images, "source": "upstream"})
}
