// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"os"
	"strconv"
)

// New returns an APIClient with dependencies defaulted to sane values
func NewAPIClient() (*APIClient, bool) {
	username, ok := os.LookupEnv("INPUT_USERNAME")
	if !ok {
		fmt.Println("env INPUT_USERNAME not set")
		return nil, false
	}

	password, ok := os.LookupEnv("INPUT_PASSWORD")
	if !ok {
		fmt.Println("env INPUT_PASSWORD not set")

		return nil, false
	}

	space, ok := os.LookupEnv("INPUT_SPACE")
	if !ok {
		fmt.Println("env INPUT_SPACE not set")

		return nil, false
	}

	return &APIClient{
		BaseURL:  "https://xiatech.atlassian.net",
		Space:    space,
		Username: username,
		Password: password,
	}, true
}

// CreatePage in confluence
func (a *APIClient) CreatePage() error {
	return nil
}

// UpdatePage in confluence
func (a *APIClient) UpdatePage(pageID, pageVersion int, pageContents bytes.Buffer) error {
	pageVersion++

	return nil
}

// FindPage in confluence
// Docs for this API endpoint are here
// https://developer.atlassian.com/cloud/confluence/rest/api-group-content/#api-api-content-get
func (a *APIClient) FindPage(title string) (int, int, bool, error) {
	fmt.Printf("%+v", a)
	lookUpURL := fmt.Sprintf("%s/wiki/rest/api/content?type=page&spaceKey=%s&title=%s", a.BaseURL, a.Space, title)

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil {
		return 0, 0, false, err
	}

	req.SetBasicAuth(a.Username, a.Password)

	resp, err := retryablehttp.NewClient().Do(req)
	if err != nil {
		return 0, 0, false, fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	fmt.Println("req: ", req, lookUpURL)

	fmt.Println("response: ", resp.Body)
	r := findPageResult{}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, 0, false, err
	}

	fmt.Println("response decoded: ", r.Results)
	spew.Dump(r)

	pageId, err := strconv.Atoi(r.Results[0].ID)
	if err != nil {
		fmt.Errorf("error converting ID to int value")
		return 0, 0, false, err
	}
	return pageId, r.Results[0].Version, true, nil
}
