// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"os"
)


// APIClient for interacting with confluence
type APIClient struct{
	BaseURL string
	headers header
}

type header struct{
	Authorization string
}

// New returns an APIClient with dependencies defaulted to sane values
func NewAPIClient() *APIClient {
	return &APIClient{
		BaseURL: "https://xiatech.atlassian.net",
		headers: header{
			Authorization: "Basic " + getHeader(),
		},
	}
}

func getHeader() string {
	uname, ok := os.LookupEnv("USERNAME")
	if !ok {
		return ""
	}
	pword, ok := os.LookupEnv("PASSWORD")
	if !ok {
		return ""
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(uname + ":" +  pword))

	return encoded
}

// CreatePage in confluence
func (a *APIClient) CreatePage() error {
	return nil
}

// UpdatePage in confluence
func (a *APIClient) UpdatePage() error {
	return nil
}

// FindPage in confluence
func (a *APIClient) FindPage () (error, bool) {
	NewAPIClient()

	lookUpURL := a.BaseURL + "/wiki/rest/api/content/"

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil{
		return err, false
	}

	resp, err :=  retryablehttp.NewClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to do the request: %w", err), false
	}

	defer func() { _ = resp.Body.Close() }()

	fmt.Println("req: ", req, lookUpURL)

	fmt.Println("response: ", resp.Body)
	return nil, true
}
