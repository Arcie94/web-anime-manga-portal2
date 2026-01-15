package controllers

import (
	"anime-tanyaayomi/internal/models"
	"anime-tanyaayomi/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		Service: services.NewUserService(),
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.Service.Register(req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(201).JSON(fiber.Map{"message": "User registered successfully"})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req models.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := c.Service.Login(req)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// In a real app, generate JWT here.
	return ctx.JSON(fiber.Map{"message": "Login successful", "user": user})
}
