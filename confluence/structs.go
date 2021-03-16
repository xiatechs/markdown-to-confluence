package confluence

// APIClient for interacting with confluence
type APIClient struct {
	BaseURL  string
	Space    string
	Username string
	Password string
}

type Page struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
	Title   string  `json:"title"`
	Version int     `json:"version"`
	Body    BodyObj `json:"body,omitempty"`
}

type BodyObj struct {
	Storage StorageObj `json:"storage"`
}

type StorageObj struct {
	Value string `json:"value"`
}

type Num struct {
	Number int `json:"number,omitempty"`
}
type findPageResult struct {
	Results []Page `json:"results"`
}
