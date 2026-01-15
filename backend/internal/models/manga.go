package models

type Manga struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Cover     string `json:"cover"`
	Thumbnail string `json:"thumbnail"` // Upstream uses this
	Link      string `json:"link,omitempty"`
	Type      string `json:"type,omitempty"`
}

type MangaDetail struct {
	Manga
	Description string    `json:"description,omitempty"`
	Author      string    `json:"author,omitempty"`
	Status      string    `json:"status,omitempty"`
	ChapterList []Chapter `json:"chapterList"`
	Chapters    []Chapter `json:"chapters,omitempty"` // Fallback as per python script
}

type Chapter struct {
	Chapter string `json:"chapter"`
	Slug    string `json:"slug"`
	Date    string `json:"date,omitempty"`
}

type ChapterImages struct {
	Images []string `json:"images"`
}
