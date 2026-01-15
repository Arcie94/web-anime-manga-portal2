package routes

import (
	"anime-tanyaayomi/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	animeController := controllers.NewAnimeController()

	anime := api.Group("/anime")
	anime.Get("/home", animeController.GetHome) // NEW: Home endpoint
	anime.Get("/ongoing", animeController.GetOngoing)
	anime.Get("/complete", animeController.GetComplete)
	anime.Get("/search", animeController.Search)
	anime.Get("/genres/:slug", animeController.GetGenre) // Use plural 'genres' to match service logic but exposed as 'genres' or 'genre'? Let's keeps consistent. Service logic mapped 'genres/%s'. Let's use /genre/:slug for clean API.
	anime.Get("/genre/:slug", animeController.GetGenre)
	anime.Get("/:slug", animeController.GetDetail)
	anime.Get("/episode/:slug", animeController.GetStream)

	mangaController := controllers.NewMangaController()
	manga := api.Group("/manga")
	manga.Get("/home", mangaController.GetHome)         // NEW: Home endpoint
	manga.Get("/trending", mangaController.GetTrending) // NEW: Trending endpoint
	manga.Get("/ongoing", mangaController.GetOngoing)   // NEW: Ongoing endpoint
	manga.Get("/search", mangaController.Search)
	manga.Get("/genres/:slug", mangaController.GetGenre) // NEW: Genres endpoint (plural)
	manga.Get("/genre/:slug", mangaController.GetGenre)  // Alias for consistency
	manga.Get("/:slug", mangaController.GetDetail)
	manga.Get("/chapter/:chapterId", mangaController.GetChapter)

	userController := controllers.NewUserController()
	auth := api.Group("/auth")
	auth.Post("/register", userController.Register)
	auth.Post("/login", userController.Login)

	bookmarkController := controllers.NewBookmarkController()
	bookmarks := api.Group("/bookmarks")
	bookmarks.Post("/", bookmarkController.AddBookmark)
	bookmarks.Get("/", bookmarkController.GetBookmarks)
	bookmarks.Delete("/:id", bookmarkController.RemoveBookmark)
}
