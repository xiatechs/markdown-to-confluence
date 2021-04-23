// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

// close body response
func httpResponseClose(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
}

// grab the page contents and return as a []byte to be used
func (a *APIClient) grabPageContentsAsJSON(contents *markdown.FileContents) ([]byte, error) {
	newPageContent := Page{
		Type:  "page",
		Title: contents.MetaData["title"].(string),
		Space: SpaceObj{Key: a.Space},
		Body: BodyObj{Storage: StorageObj{
			Value:          string(contents.Body),
			Representation: "storage",
		}},
	}

	newPageContentsJSON, err := json.Marshal(newPageContent)
	if err != nil {
		return nil, err
	}

	return newPageContentsJSON, nil
}

// CreatePage in confluence
// todo: function not tested live on confluence yet! test written on expected results
func (a *APIClient) CreatePage(contents *markdown.FileContents) error {
	newPageContentsJSON, err := a.grabPageContentsAsJSON(contents)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/wiki/rest/api/content", a.BaseURL)

	req, err := retryablehttp.NewRequest(http.MethodPost, URL, newPageContentsJSON)
	if err != nil {
		return err
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}

	defer httpResponseClose(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create confluence page: %s", resp.Status)
	}

	return nil
}

// UpdatePage updates a confluence page with our newly created data and increases the
// version by 1 each time.
func (a *APIClient) UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents) error {
	newPageContentsJSON, err := a.grabPageContentsAsJSON(pageContents)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/wiki/rest/api/content/%d", a.BaseURL, pageID)

	req, err := retryablehttp.NewRequest(http.MethodPut, URL, bytes.NewBuffer(newPageContentsJSON))
	if err != nil {
		log.Println(err)
		return err
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		fmt.Println("error was: ", resp.Status, err)
		return fmt.Errorf("failed to do the request: %w", err)
	}

	defer httpResponseClose(resp)

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error ioutil", err)
	}

	fmt.Println("response: ", string(r))

	return nil
}

// FindPage in confluence
// Docs for this API endpoint are here
// https://developer.atlassian.com/cloud/confluence/rest/api-group-content/#api-api-content-get
func (a *APIClient) FindPage(title string) (*PageResults, error) {
	lookUpURL := fmt.Sprintf("%s/wiki/rest/api/content?expand=version&type=page&spaceKey=%s&title=%s",
		a.BaseURL, a.Space, title)

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(a.Username, a.Password)

	resp, err := a.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer httpResponseClose(resp)

	pageResultVar := PageResults{}

	if err = json.NewDecoder(resp.Body).Decode(&pageResultVar); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}

		return nil, err
	}

	if len(pageResultVar.Results) == 0 {
		return nil, nil
	}

	return &pageResultVar, nil
}
