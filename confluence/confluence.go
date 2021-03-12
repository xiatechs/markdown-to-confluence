// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"os"
)


// APIClient for interacting with confluence
type APIClient struct{
	BaseURL string
	Space string
	Username string
	Password string
}

type Page struct{
	ID string `json:"id"`
	Type string `json:"type"`
	Status string `json:"status"`
	Title string `json:"title"`
}

type findPageResult struct {
	Results []Page `json:"results"`
}


// New returns an APIClient with dependencies defaulted to sane values
func NewAPIClient() (*APIClient,bool) {
	username, ok := os.LookupEnv("INPUT_USERNAME")
	if !ok {
		return nil, false
	}

	password, ok := os.LookupEnv("INPUT_PASSWORD")
	if !ok {
		return nil, false
	}

	space, ok := os.LookupEnv("INPUT_SPACE")
	if !ok {
		return nil, false
	}

	return &APIClient{
		BaseURL: "https://xiatech.atlassian.net",
		Space: space,
		Username: username,
		Password: password,
	}, true
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
// Docs for this API endpoint are here
// https://developer.atlassian.com/cloud/confluence/rest/api-group-content/#api-api-content-get
func (a *APIClient) FindPage (title string) (error, bool) {
	fmt.Printf("%+v", a)
	lookUpURL := fmt.Sprintf("%s/wiki/rest/api/content?type=page&spaceKey=%s&title=%s",a.BaseURL,a.Space, title)

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil{
		return err, false
	}

	req.SetBasicAuth(a.Username,a.Password)

	resp, err :=  retryablehttp.NewClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to do the request: %w", err), false
	}

	defer func() { _ = resp.Body.Close() }()

	fmt.Println("req: ", req, lookUpURL)

	fmt.Println("response: ", resp.Body)
	r := findPageResult{}

	if err := json.NewDecoder(resp.Body).Decode(&r); err !=nil {
		return err,true
	}

	spew.Dump(r)
	return nil, true
}
