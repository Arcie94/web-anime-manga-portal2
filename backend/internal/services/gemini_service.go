package services

import (
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
	Year        string `json:"year"`
	Rating      string `json:"rating"`
	Synopsis    string `json:"synopsis"`
	Status      string `json:"status"` // e.g. "Completed", "Ongoing"
	LastUpdated int64  `json:"-"`
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
func EnrichData(title string, mediaType string) EnrichedData {
	// Check Cache First
	cacheMutex.RLock()
	cached, exists := enrichmentCache[title]
	cacheMutex.RUnlock()
	if exists {
		// Simple TTL (e.g. 24 hours) - though for years/static data it could be forever
		if time.Now().Unix()-cached.LastUpdated < 86400 {
			return cached
		}
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = "AIzaSyBWZ4cXMHDd5U2-zpaCn54sCExFSDQcH6U" // Fallback/User Provided provided in session
	}

	// Construct Prompt
	prompt := fmt.Sprintf(`Identify the %s "%s". 
    Return a strictly valid JSON object (no markdown formatting) with these fields:
    - "year": (string) Release year (e.g. "2023").
    - "rating": (string) Average score 0-10 (e.g. "8.5").
    - "status": (string) "Ongoing" or "Completed".
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

			// Update Cache
			cacheMutex.Lock()
			enrichmentCache[title] = result
			cacheMutex.Unlock()

			return result
		}
	}

	return EnrichedData{}
}
