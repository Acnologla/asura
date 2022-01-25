package interaction

type Footer struct {
}

type Image struct {
	URL string `json:"url"`
}

type Thumbnail struct {
}

type Author struct {
}

type Field struct {
}

type Embed struct {
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Color       int    `json:"color"`
	Image       *Image `json:"image"`
}
