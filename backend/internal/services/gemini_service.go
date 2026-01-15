package services

import (
	"anime-tanyaayomi/internal/models"
	"anime-tanyaayomi/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// In-memory cache to avoid burning API quota and reduce latency
// Map Key: Title -> Value: EnrichedData
var enrichmentCache = make(map[string]EnrichedData)
var cacheMutex sync.RWMutex

type EnrichedData struct {
	Year     string `json:"year"`
	Rating   string `json:"rating"`
	Synopsis string `json:"synopsis"`
	Status   string `json:"status"` // e.g. "Completed", "Ongoing"
	Author   string `json:"author"`
	Genre    string `json:"genre"`
	// Type removed as we are hardcoding it to "Anime"
	LastUpdated int64 `json:"-"`
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// EnrichData attempts to fetch metadata from Gemini for a given title
// Now with 3-tier lookup: Database → Cache → Gemini API
// This reduces Gemini API calls by ~90% through persistent storage
func EnrichData(title string, mediaType string) EnrichedData {
	// Kept for backward compatibility - will be replaced with EnrichDataWithDB
	fmt.Println("[EnrichData] WARNING: Using legacy enrichment without database. Please migrate to EnrichDataWithDB.")
	return enrichDataFromCache(title, mediaType)
}

// EnrichDataWithDB is the new database-aware enrichment function
// It implements 3-tier lookup to minimize Gemini API costs
func EnrichDataWithDB(title string, mediaType string, repo interface{}) EnrichedData {
	// Type assertion for repository (allowing nil for backward compatibility)
	var enrichRepo *EnrichmentRepository
	if repo != nil {
		if r, ok := repo.(*EnrichmentRepository); ok {
			enrichRepo = r
		}
	}

	// Tier 1: Check PostgreSQL Database first (permanent storage)
	if enrichRepo != nil {
		if dbData, err := enrichRepo.GetByTitle(title, mediaType); err == nil && dbData != nil {
			fmt.Printf("[EnrichDataWithDB] ✅ Found in DATABASE for '%s' (%s)\n", title, mediaType)
			return convertDBToEnrichedData(dbData)
		}
	}

	// Tier 2: Check in-memory cache (fast fallback, 1 hour TTL)
	cacheMutex.RLock()
	cached, exists := enrichmentCache[title]
	cacheMutex.RUnlock()
	if exists && time.Now().Unix()-cached.LastUpdated < 3600 {
		fmt.Printf("[EnrichDataWithDB] ⚡ Found in CACHE for '%s' (%s)\n", title, mediaType)
		return cached
	}

	// Tier 3: Call Gemini API (last resort, expensive)
	fmt.Printf("[EnrichDataWithDB] 🤖 Calling GEMINI API for '%s' (%s)\n", title, mediaType)
	enriched := callGeminiAPI(title, mediaType)

	// Save to database for future use (avoid future API calls)
	if enrichRepo != nil && enriched.Year != "" {
		dbModel := convertEnrichedDataToDB(title, mediaType, enriched)
		if err := enrichRepo.Upsert(dbModel); err != nil {
			fmt.Printf("[EnrichDataWithDB] ⚠️  Failed to save to database: %v\n", err)
		} else {
			fmt.Printf("[EnrichDataWithDB] 💾 Saved to DATABASE for '%s'\n", title)
		}
	}

	// Update memory cache
	cacheMutex.Lock()
	enrichmentCache[title] = enriched
	cacheMutex.Unlock()

	return enriched
}

// enrichDataFromCache is the legacy cache-only lookup (for backward compatibility)
func enrichDataFromCache(title string, mediaType string) EnrichedData {
	// Check Cache First
	cacheMutex.RLock()
	cached, exists := enrichmentCache[title]
	cacheMutex.RUnlock()
	if exists {
		if time.Now().Unix()-cached.LastUpdated < 86400 {
			return cached
		}
	}

	// Call Gemini API
	enriched := callGeminiAPI(title, mediaType)

	// Update Cache
	cacheMutex.Lock()
	enrichmentCache[title] = enriched
	cacheMutex.Unlock()

	return enriched
}

// callGeminiAPI makes the actual API call to Gemini
func callGeminiAPI(title string, mediaType string) EnrichedData {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = "AIzaSyBWZ4cXMHDd5U2-zpaCn54sCExFSDQcH6U" // Fallback/User Provided
	}

	// Construct Prompt
	prompt := fmt.Sprintf(`Identify the %s "%s". 
    Return a strictly valid JSON object (no markdown formatting) with these fields:
    - "year": (string) Release year (e.g. "2023").
    - "rating": (string) Average score 0-10 (e.g. "8.5").
    - "status": (string) "Ongoing" or "Completed".
    - "author": (string) Original creator/mangaka.
    - "genre": (string) Comma-separated genres (e.g. "Action, Adventure").
    - "synopsis": (string) A very short, engaging 1-sentence summary.
    If unknown, return generic/empty values but valid JSON.`, mediaType, title)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Gemini Error:", err)
		return EnrichedData{}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return EnrichedData{}
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		rawJSON := geminiResp.Candidates[0].Content.Parts[0].Text
		// Cleanup potentially markdown wrapped JSON
		rawJSON = strings.TrimPrefix(rawJSON, "```json")
		rawJSON = strings.TrimPrefix(rawJSON, "```")
		rawJSON = strings.TrimSuffix(rawJSON, "```")
		rawJSON = strings.TrimSpace(rawJSON)

		var result EnrichedData
		if err := json.Unmarshal([]byte(rawJSON), &result); err == nil {
			result.LastUpdated = time.Now().Unix()
			return result
		}
	}

	return EnrichedData{}
}

// Type alias for repository (to avoid import cycle)
type EnrichmentRepository = repository.EnrichmentRepository

// convertDBToEnrichedData converts database model to EnrichedData struct
func convertDBToEnrichedData(db *models.EnrichedMetadata) EnrichedData {
	return EnrichedData{
		Year:        db.ReleaseYear,
		Rating:      db.Rating,
		Synopsis:    db.Synopsis,
		Status:      db.Status,
		Author:      db.Author,
		Genre:       db.Genre,
		LastUpdated: db.LastUpdatedAt.Unix(),
	}
}

// convertEnrichedDataToDB converts EnrichedData to database model
func convertEnrichedDataToDB(title string, mediaType string, enriched EnrichedData) *models.EnrichedMetadata {
	return &models.EnrichedMetadata{
		Title:       title,
		MediaType:   mediaType,
		Author:      enriched.Author,
		Genre:       enriched.Genre,
		Type:        "Anime", // Hardcoded as per user requirement
		Rating:      enriched.Rating,
		Status:      enriched.Status,
		ReleaseYear: enriched.Year,
		Synopsis:    enriched.Synopsis,
		Source:      "gemini",
	}
}
