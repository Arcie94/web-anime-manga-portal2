package services

import (
	"anime-tanyaayomi/internal/database"
	"anime-tanyaayomi/internal/models"
)

type BookmarkService struct{}

func NewBookmarkService() *BookmarkService {
	return &BookmarkService{}
}

func (s *BookmarkService) AddBookmark(b models.Bookmark) error {
	query := `INSERT INTO bookmarks (user_id, type, slug, title, cover_image) 
              VALUES ($1, $2, $3, $4, $5) 
              ON CONFLICT (user_id, type, slug) DO NOTHING`
	_, err := database.DB.Exec(query, b.UserID, b.Type, b.Slug, b.Title, b.CoverImage)
	return err
}

func (s *BookmarkService) GetBookmarks(userID int) ([]models.Bookmark, error) {
	query := `SELECT id, user_id, type, slug, title, cover_image, created_at FROM bookmarks WHERE user_id = $1`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		if err := rows.Scan(&b.ID, &b.UserID, &b.Type, &b.Slug, &b.Title, &b.CoverImage, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

func (s *BookmarkService) RemoveBookmark(userID int, bookmarkID int) error {
	query := `DELETE FROM bookmarks WHERE id = $1 AND user_id = $2`
	_, err := database.DB.Exec(query, bookmarkID, userID)
	return err
}
