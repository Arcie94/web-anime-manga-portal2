package models

// Base Anime model (existing - keep as is)
type Anime struct {
	Title         string      `json:"title"`
	Slug          string      `json:"slug"`
	AnimeID       string      `json:"animeId"`
	Cover         string      `json:"cover"`
	Poster        string      `json:"poster"`    // Upstream uses this for Anime
	Thumbnail     string      `json:"thumbnail"` // Upstream fallback
	Image         string      `json:"image"`     // Another potential fallback
	Synopsis      interface{} `json:"synopsis"`
	Genre         string      `json:"genre"`
	ReleaseDate   string      `json:"releaseDate"`
	TotalEpisodes string      `json:"totalEpisodes"`
	Status        string      `json:"status"` // Added for Sankavollerei
	Rating        string      `json:"rating"` // Added for Sankavollerei
	Author        string      `json:"author"` // Added for AI Enrichment
	Type          string      `json:"type"`   // Added for AI Enrichment (TV, Movie, OVA)
}

// Episode model (existing - keep as is)
type Episode struct {
	Title     string      `json:"title"`
	EpisodeID string      `json:"episodeId"`
	Slug      string      `json:"slug"`
	Episode   string      `json:"episode"`
	Eps       interface{} `json:"eps"` // Can be string or number from API
}

// AnimeDetail model (existing - keep as is)
type AnimeDetail struct {
	Anime
	EpisodeList []Episode `json:"episodeList"`
}

// StreamData model (existing - keep as is)
type StreamData struct {
	DefaultStreamingUrl string `json:"defaultStreamingUrl"`
	StreamLink          string `json:"stream_link"`
	Url                 string `json:"url"`
	Title               string `json:"title"`
	AnimeID             string `json:"animeId"`
}

// APIResponse model (existing - keep as is)
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ============ NEW MODELS FOR SANKAVOLLEREI ============

// AnimeListWrapper wraps the animeList array for ongoing/completed
type AnimeListWrapper struct {
	Href         string  `json:"href"`
	OtakudesuUrl string  `json:"otakudesuUrl"`
	AnimeList    []Anime `json:"animeList"`
}

// HomeResponse - Response from /anime/home
type HomeResponse struct {
	Status string `json:"status"`
	Data   struct {
		Ongoing   AnimeListWrapper `json:"ongoing"`
		Completed AnimeListWrapper `json:"completed"`
	} `json:"data"`
}

// SearchResponse - Response from /anime/search/:keyword
type SearchResponse struct {
	Data struct {
		AnimeList []Anime `json:"animeList"`
	} `json:"data"`
}

// AnimeDetailResponse - Response from /anime/anime/:slug
type AnimeDetailResponse struct {
	Data AnimeDetail `json:"data"`
}

// StreamServer represents a streaming server option in a quality's serverList
type StreamServer struct {
	Title    string `json:"title"`
	ServerID string `json:"serverId"`
	Href     string `json:"href"`
}

// QualityOption represents a video quality with its server list
type QualityOption struct {
	Title      string         `json:"title"`
	ServerList []StreamServer `json:"serverList"`
}

// ServerData represents the server object containing qualities
type ServerData struct {
	Qualities []QualityOption `json:"qualities"`
}

// DownloadURL represents a download option
type DownloadURL struct {
	Quality string `json:"quality"`
	URL     string `json:"url"`
}

// StreamResponse - Response from /anime/episode/:episodeId
type StreamResponse struct {
	Data struct {
		Title               string      `json:"title"`
		DefaultStreamingUrl string      `json:"defaultStreamingUrl"`
		StreamLink          string      `json:"stream_link"`
		URL                 string      `json:"url"`
		AnimeID             string      `json:"animeId"`
		Server              ServerData  `json:"server"`      // Changed from []StreamServer to ServerData
		DownloadURL         interface{} `json:"downloadUrl"` // Can be array or null
	} `json:"data"`
}

// LatestEpisode represents a recent episode update
type LatestEpisode struct {
	Title     string `json:"title"`
	EpisodeID string `json:"episodeId"`
	Slug      string `json:"slug"`
	Poster    string `json:"poster"`
	AnimeID   string `json:"animeId"`
	Source    string `json:"source"` // Which site (otakudesu, samehadaku, etc.)
}

// LatestResponse - Response from /anime/stream/latest
type LatestResponse struct {
	Data struct {
		Episodes []LatestEpisode `json:"episodes"`
	} `json:"data"`
}

// ============ MANGA MODELS ============
// Manga uses similar structure as Anime, but with some differences
type Manga struct {
	Title       string      `json:"title"`
	Link        string      `json:"link"` // For slug extraction: "/manga/slug-name/"
	Slug        string      `json:"slug"` // Extracted from link
	Cover       string      `json:"cover"`
	Image       string      `json:"image"` // Sankavollerei uses 'image' field
	Poster      string      `json:"poster"`
	Thumbnail   string      `json:"thumbnail"`
	Synopsis    interface{} `json:"synopsis"`
	Genre       string      `json:"genre"`
	Status      string      `json:"status"`
	Chapter     string      `json:"chapter"` // Latest chapter from list
	TimeAgo     string      `json:"time_ago"`
	ReleaseDate string      `json:"releaseDate"`
}

type Chapter struct {
	Title     string `json:"title"`
	Chapter   string `json:"chapter"` // Display name like "Chapter 123"
	ChapterID string `json:"chapterId"`
	Slug      string `json:"slug"`
	Date      string `json:"date"`
}

type MangaDetail struct {
	Title    string      `json:"title"`
	Image    string      `json:"image"`
	Synopsis interface{} `json:"synopsis"`
	Metadata struct {
		Author string `json:"author"`
		Status string `json:"status"`
		Type   string `json:"type"`
	} `json:"metadata"`
	Chapters []Chapter `json:"chapters"` // Sankavollerei uses 'chapters' not 'chapterList'
}

// MangaDetailResponse - Response from /comic/comic/:slug
type MangaDetailResponse struct {
	Title    string      `json:"title"`
	Image    string      `json:"image"`
	Synopsis interface{} `json:"synopsis"`
	Metadata struct {
		Author string `json:"author"`
		Status string `json:"status"`
		Type   string `json:"type"`
	} `json:"metadata"`
	Chapters []Chapter `json:"chapters"`
}

// MangaList Response - Response from manga list endpoints
type MangaListResponse struct {
	Data struct {
		MangaList []Manga `json:"mangaList"`
	} `json:"data"`
}

// ChapterResponse - Response from /comic/chapter/:chapterId
type ChapterResponse struct {
	Title     string   `json:"title"`
	ChapterID string   `json:"chapterId"`
	MangaID   string   `json:"mangaId"`
	Images    []string `json:"images"` // Array of image URLs for the chapter
	NextSlug  string   `json:"nextSlug,omitempty"`
	PrevSlug  string   `json:"prevSlug,omitempty"`
}
