package repository

import (
	"anime-tanyaayomi/internal/models"
	"database/sql"
	"fmt"
)

// EnrichmentRepository handles database operations for enriched metadata
type EnrichmentRepository struct {
	DB *sql.DB
}

// NewEnrichmentRepository creates a new repository instance
func NewEnrichmentRepository(db *sql.DB) *EnrichmentRepository {
	return &EnrichmentRepository{DB: db}
}

// GetByTitle retrieves enriched metadata by title and media type
// Returns nil error if not found (allowing fallback to Gemini)
func (r *EnrichmentRepository) GetByTitle(title string, mediaType string) (*models.EnrichedMetadata, error) {
	query := `
		SELECT id, title, media_type, slug, author, genre, type, rating, status, 
		       release_year, synopsis, source, last_updated_at, created_at
		FROM enriched_metadata
		WHERE title = $1 AND media_type = $2
		LIMIT 1
	`

	var data models.EnrichedMetadata
	err := r.DB.QueryRow(query, title, mediaType).Scan(
		&data.ID,
		&data.Title,
		&data.MediaType,
		&data.Slug,
		&data.Author,
		&data.Genre,
		&data.Type,
		&data.Rating,
		&data.Status,
		&data.ReleaseYear,
		&data.Synopsis,
		&data.Source,
		&data.LastUpdatedAt,
		&data.CreatedAt,
	)

	if err == sql.ErrNoRows {
		// Not found - return nil to allow fallback to Gemini
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error querying enriched metadata: %w", err)
	}

	return &data, nil
}

// Upsert inserts or updates enriched metadata
// Uses PostgreSQL's ON CONFLICT to handle duplicates
func (r *EnrichmentRepository) Upsert(data *models.EnrichedMetadata) error {
	query := `
		INSERT INTO enriched_metadata 
			(title, media_type, slug, author, genre, type, rating, status, release_year, synopsis, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (title, media_type) 
		DO UPDATE SET
			slug = EXCLUDED.slug,
			author = COALESCE(NULLIF(EXCLUDED.author, ''), enriched_metadata.author),
			genre = COALESCE(NULLIF(EXCLUDED.genre, ''), enriched_metadata.genre),
			type = COALESCE(NULLIF(EXCLUDED.type, ''), enriched_metadata.type),
			rating = COALESCE(NULLIF(EXCLUDED.rating, ''), enriched_metadata.rating),
			status = COALESCE(NULLIF(EXCLUDED.status, ''), enriched_metadata.status),
			release_year = COALESCE(NULLIF(EXCLUDED.release_year, ''), enriched_metadata.release_year),
			synopsis = COALESCE(NULLIF(EXCLUDED.synopsis, ''), enriched_metadata.synopsis),
			source = EXCLUDED.source,
			last_updated_at = NOW()
		RETURNING id
	`

	var id int
	err := r.DB.QueryRow(
		query,
		data.Title,
		data.MediaType,
		data.Slug,
		data.Author,
		data.Genre,
		data.Type,
		data.Rating,
		data.Status,
		data.ReleaseYear,
		data.Synopsis,
		data.Source,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("error upserting enriched metadata: %w", err)
	}

	data.ID = id
	return nil
}

// GetAll retrieves all enriched metadata (useful for admin/debugging)
func (r *EnrichmentRepository) GetAll(limit int) ([]models.EnrichedMetadata, error) {
	query := `
		SELECT id, title, media_type, slug, author, genre, type, rating, status, 
		       release_year, synopsis, source, last_updated_at, created_at
		FROM enriched_metadata
		ORDER BY last_updated_at DESC
		LIMIT $1
	`

	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying all enriched metadata: %w", err)
	}
	defer rows.Close()

	var results []models.EnrichedMetadata
	for rows.Next() {
		var data models.EnrichedMetadata
		err := rows.Scan(
			&data.ID,
			&data.Title,
			&data.MediaType,
			&data.Slug,
			&data.Author,
			&data.Genre,
			&data.Type,
			&data.Rating,
			&data.Status,
			&data.ReleaseYear,
			&data.Synopsis,
			&data.Source,
			&data.LastUpdatedAt,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning enriched metadata: %w", err)
		}
		results = append(results, data)
	}

	return results, nil
}
