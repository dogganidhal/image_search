package domain

type ImageLabelDocument struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type ImageDocument struct {
	Id                 string  `json:"id"`
	OriginalUrl        string  `json:"original_url"`
	OriginalLandingURL string  `json:"original_landing_url"`
	AuthorProfileURL   string  `json:"author_profile_url"`
	Author             string  `json:"author"`
	Title              string  `json:"title"`
	OriginalSize       float64 `json:"original_size"`
	OriginalMD5        string  `json:"original_md_5"`
	Thumbnail300KUrl   string  `json:"thumbnail_300k_url"`
}
