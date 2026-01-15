package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Bookmark struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Type       string    `json:"type"` // 'anime' or 'manga'
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	CoverImage string    `json:"cover_image"`
	CreatedAt  time.Time `json:"created_at"`
}
