// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/markdown"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	confluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
	confluenceAPIKeyEnv   = "INPUT_CONFLUENCE_API_KEY"
	confluenceSpaceEnv    = "INPUT_CONFLUENCE_SPACE"
)

// New returns an APIClient with dependencies defaulted to sane values
func NewAPIClient() (*APIClient, bool) {
	password, exists := os.LookupEnv(confluenceAPIKeyEnv)
	if !exists {
		fmt.Printf("Environment variable not set for %s", confluenceAPIKeyEnv)
	} else {
		log.Printf("API KEY: %s", password)
	}

	username, exists := os.LookupEnv(confluenceUsernameEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceUsernameEnv)
	} else {
		log.Printf("API KEY: %s", username)
	}

	space, exists := os.LookupEnv(confluenceSpaceEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceSpaceEnv)
	} else {
		log.Printf("SPACE: %s", space)
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

// UpdatePage updates a confluence page with our newly created data and increases the
// version by 1 each time.
func (a *APIClient) UpdatePage(pageID, pageVersion int, pageContents *markdown.FileContents) error {

	var newPageJson PutPageContent
	newPageJson.Type = "page"
	newPageJson.Title = pageContents.MetaData["title"].(string)
	newPageJson.Version.Number = pageVersion + 1
	newPageJson.Body.Storage.Value = string(pageContents.Body)

	fmt.Println(pageID, newPageJson) //todo: remove

	URL := fmt.Sprintf("%s/wiki/rest/api/content/%d?expand=%v", a.BaseURL, pageID, newPageJson)

	fmt.Println(URL) //todo: remove

	req, err := retryablehttp.NewRequest(http.MethodPut, URL, nil)
	if err != nil {
		return err
	}

	fmt.Println(req) //todo: remove

	req.SetBasicAuth(a.Username, a.Password)

	resp, err := retryablehttp.NewClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	fmt.Println("1") //todo: remove

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
