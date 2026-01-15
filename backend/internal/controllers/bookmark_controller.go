package controllers

import (
	"anime-tanyaayomi/internal/models"
	"anime-tanyaayomi/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BookmarkController struct {
	Service *services.BookmarkService
}

func NewBookmarkController() *BookmarkController {
	return &BookmarkController{
		Service: services.NewBookmarkService(),
	}
}

// Helper to get UserID (In real app, extract from JWT)
func getUserID(ctx *fiber.Ctx) int {
	// For demo: get user_id from header or query param.
	// WARNING: Insecure. Access control should be implemented properly.
	uid, _ := strconv.Atoi(ctx.Get("X-User-ID"))
	if uid == 0 {
		return 0
	}
	return uid
}

func (c *BookmarkController) AddBookmark(ctx *fiber.Ctx) error {
	userID := getUserID(ctx)
	if userID == 0 {
		return ctx.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var b models.Bookmark
	if err := ctx.BodyParser(&b); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	b.UserID = userID

	if err := c.Service.AddBookmark(b); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(201).JSON(fiber.Map{"message": "Bookmark added"})
}

func (c *BookmarkController) GetBookmarks(ctx *fiber.Ctx) error {
	userID := getUserID(ctx)
	if userID == 0 {
		return ctx.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	bookmarks, err := c.Service.GetBookmarks(userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": bookmarks})
}

func (c *BookmarkController) RemoveBookmark(ctx *fiber.Ctx) error {
	userID := getUserID(ctx)
	if userID == 0 {
		return ctx.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.Service.RemoveBookmark(userID, id); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"message": "Bookmark removed"})
}
