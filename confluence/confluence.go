// Package confluence provides functionality for interacting with the confluence API
// Specifically managing pages
package confluence

// API for interacting with confluence
type API struct{}

// New returns an API with dependencies defaulted to sane values
func New() *API {
	return &API{}
}

// CreatePage in confluence
func (a *API) CreatePage() error {
	return nil
}

// UpdatePage in confluence
func (a *API) UpdatePage() error {
	return nil
}

// FindPage in confluence
func (a *API) FindPage() error {
	return nil
}
