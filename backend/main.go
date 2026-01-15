package main

import (
	"log"

	"anime-tanyaayomi/internal/database"
	"anime-tanyaayomi/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Tidak ada file .env ditemukan")
	}

	// Database
	database.Connect()

	app := fiber.New()

	// Middleware
	app.Use(cors.New())

	// Routes
	routes.SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Anime-TanyaAyomi Backend Running")
	})

	log.Fatal(app.Listen(":3000"))
}
