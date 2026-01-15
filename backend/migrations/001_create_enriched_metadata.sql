-- Migration: Create enriched_metadata table for storing AI-enriched anime/manga metadata
-- This reduces Gemini API calls by ~90% through persistent storage

CREATE TABLE IF NOT EXISTS enriched_metadata (
    id SERIAL PRIMARY KEY,
    
    -- Unique identifier fields
    title VARCHAR(500) NOT NULL,
    media_type VARCHAR(10) NOT NULL CHECK (media_type IN ('anime', 'manga')),
    slug VARCHAR(500),
    
    -- Enriched fields from Gemini AI
    author VARCHAR(255),
    genre TEXT,
    type VARCHAR(50),
    rating VARCHAR(10),
    status VARCHAR(50),
    release_year VARCHAR(10),
    synopsis TEXT,
    
    -- Metadata tracking
    source VARCHAR(50) DEFAULT 'gemini' CHECK (source IN ('gemini', 'manual', 'api')),
    last_updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Ensure one entry per title+media_type combination
    CONSTRAINT unique_title_media UNIQUE(title, media_type)
);

-- Indexes for fast lookup
CREATE INDEX IF NOT EXISTS idx_enriched_title ON enriched_metadata(title);
CREATE INDEX IF NOT EXISTS idx_enriched_slug ON enriched_metadata(slug);
CREATE INDEX IF NOT EXISTS idx_enriched_type ON enriched_metadata(media_type);
CREATE INDEX IF NOT EXISTS idx_enriched_updated ON enriched_metadata(last_updated_at);

-- Comment on table
COMMENT ON TABLE enriched_metadata IS 'Stores AI-enriched metadata for anime/manga to reduce API calls to Gemini';
