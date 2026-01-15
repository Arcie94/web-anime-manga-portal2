package models

import "time"

// EnrichedMetadata represents AI-enriched metadata stored in PostgreSQL database
// This reduces Gemini API calls by storing results permanently
type EnrichedMetadata struct {
	ID        int    `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	MediaType string `db:"media_type" json:"mediaType"` // "anime" or "manga"
	Slug      string `db:"slug" json:"slug,omitempty"`

	// Enriched fields from Gemini AI
	Author      string `db:"author" json:"author,omitempty"`
	Genre       string `db:"genre" json:"genre,omitempty"`
	Type        string `db:"type" json:"type,omitempty"`
	Rating      string `db:"rating" json:"rating,omitempty"`
	Status      string `db:"status" json:"status,omitempty"`
	ReleaseYear string `db:"release_year" json:"releaseYear,omitempty"`
	Synopsis    string `db:"synopsis" json:"synopsis,omitempty"`

	// Metadata tracking
	Source        string    `db:"source" json:"source"` // "gemini", "manual", or "api"
	LastUpdatedAt time.Time `db:"last_updated_at" json:"lastUpdatedAt"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}
