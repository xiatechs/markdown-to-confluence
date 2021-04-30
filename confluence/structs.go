package confluence

// PageResults contains the returned page values
type PageResults struct {
	Results []Page `json:"results"`
}

// Page holds returned confluence data
type Page struct {
	ID        string        `json:"id,omitempty"`
	Type      string        `json:"type"`
	Status    string        `json:"status"`
	Title     string        `json:"title"`
	Space     SpaceObj      `json:"space,omitempty"`
	Version   VersionObj    `json:"version,omitempty"`
	Ancestors []AncestorObj `json:"ancestors,omitempty"`
	Body      BodyObj       `json:"body,omitempty"`
}

// AncestorObj contains the page ID of a parent page
type AncestorObj struct {
	ID int `json:"id"`
}

// SpaceObj contains the confluence space value
type SpaceObj struct {
	Key string `json:"key,omitempty"`
}

// BodyObj stores body object
type BodyObj struct {
	Storage StorageObj `json:"storage"`
}

// StorageObj stores storage object
type StorageObj struct {
	Value          string `json:"value"`
	Representation string `json:"representation,omitempty"`
}

// VersionObj stores page version increased by 1 for PUT request
type VersionObj struct {
	Number int `json:"number"`
}
