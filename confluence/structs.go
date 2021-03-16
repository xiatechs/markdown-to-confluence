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
	Version Num     `json:"version"`
	Body    BodyObj `json:"body,omitempty"`
}

type BodyObj struct {
	Storage StorageObj `json:"storage"`
}

type StorageObj struct {
	Value          string `json:"value"`
	Representation string `json:"representation,omitempty"`
}

type Num struct {
	Number int64 `json:"number,omitempty"`
}
type findPageResult struct {
	Results []Page `json:"results"`
}

type PutPageContent struct {
	Type    string     `json:"type"`
	Title   string     `json:"title,omitempty"`
	Version VersionObj `json:"version"`
	Body    BodyObj    `json:"body"`
}

type VersionObj struct {
	Number int `json:"number"`
}
