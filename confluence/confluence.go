// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

// newPageResults function takes in a http response and
// decodes the response body into a PageResults struct that is returned
func newPageResults(resp *http.Response) (*PageResults, error) {
	pageResultVar := PageResults{}

	if err := json.NewDecoder(resp.Body).Decode(&pageResultVar); err != nil {
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

// grabPageContents method takes in page contents, the parent page id (root) and boolean
// confirming whether or not the page is a parent page folder (isroot)
// and returns byte of page contents
func (a *APIClient) grabPageContents(contents *markdown.FileContents, root int, isroot bool) ([]byte, error) {
	if isroot {
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

	newPageContent := Page{
		Type:      "page",
		Title:     contents.MetaData["title"].(string),
		Space:     SpaceObj{Key: a.Space},
		Ancestors: []AncestorObj{{ID: root}},
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

// CreatePage method takes root (root page id) and page contents and bool (is page root?)
// and generates a page in confluence and returns the generated page ID
func (a *APIClient) CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error) {
	newPageContentsJSON, err := a.grabPageContents(contents, root, isroot)
	if err != nil {
		return 0, err
	}

	URL := fmt.Sprintf("%s/wiki/rest/api/content", a.BaseURL)

	req, err := retryablehttp.NewRequest(http.MethodPost, URL, newPageContentsJSON)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Println("error was: ", resp.Status, err)
		return 0, fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to create confluence page: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	var output struct {
        	id int `json:"id"`
    	}

	err = decoder.Decode(&output)
	if err != nil {
		log.Println("error was: ", resp.Status, err)
		return 0, nil
	}

	return output.id, nil
}

// updatePageContents method updates the page contents and return as a []byte JSON to be used
func (a *APIClient) updatePageContents(pageVersion int64, contents *markdown.FileContents) ([]byte, error) {
	newPageContent := Page{
		Type:  "page",
		Title: contents.MetaData["title"].(string),
		Version: VersionObj{
			Number: int(pageVersion) + 1,
		},
		Space: SpaceObj{Key: a.Space},
		Body: BodyObj{
			Storage: StorageObj{
				Value:          string(contents.Body),
				Representation: "storage",
			},
		},
	}

	newPageContentsJSON, err := json.Marshal(newPageContent)
	if err != nil {
		return nil, err
	}

	return newPageContentsJSON, nil
}

// DeletePage deletes a confluence page by page ID
func (a *APIClient) DeletePage(pageID int) error {
	URL := fmt.Sprintf("%s/wiki/rest/api/content/%d", a.BaseURL, pageID)

	req, err := retryablehttp.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Println("error was: ", resp.Status, err)
		return fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		log.Println("error ioutil", err)
	}

	return nil
}

// UpdatePage updates a confluence page with our newly created data and increases the
// version by 1 each time.
func (a *APIClient) UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents) error {
	newPageContentsJSON, err := a.updatePageContents(pageVersion, pageContents)
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
		log.Println("error was: ", resp.Status, err)
		return fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		log.Println("error ioutil", err)
	}

	return nil
}

// createFindPageRequest method takes in a title (page title) and searches for page
// in confluence
func (a *APIClient) createFindPageRequest(title string) (*retryablehttp.Request, error) {
	lookUpURL := fmt.Sprintf("%s/wiki/rest/api/content?expand=version&type=page&spaceKey=%s&title=%s",
		a.BaseURL, a.Space, title)

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(a.Username, a.Password)

	return req, nil
}

// createFindPagesRequest method takes in a page ID and searches for page
// in confluence as well as children pages
func (a *APIClient) createFindPagesRequest(id string) (*retryablehttp.Request, error) {
	targetURL := fmt.Sprintf(common.ConfluenceBaseURL + "/wiki/rest/api/content/" + id + "/child/page")

	req, err := retryablehttp.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(a.Username, a.Password)

	return req, nil
}

// findPageRequest method takes in page title (and bool for if we are collecting multiple pages i.e
// parent and child pages
func (a *APIClient) findPageRequest(title string, many bool) (*retryablehttp.Request, error) {
	var req *retryablehttp.Request

	var err error

	if many {
		req, err = a.createFindPagesRequest(title)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = a.createFindPageRequest(title)
		if err != nil {
			return nil, err
		}
	}

	return req, err
}

// FindPage in confluence
// Docs for this API endpoint are here
// https://developer.atlassian.com/cloud/confluence/rest/api-group-content/#api-api-content-get
func (a *APIClient) FindPage(title string, many bool) (*PageResults, error) {
	req, err := a.findPageRequest(title, many)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	results, err := newPageResults(resp)
	if err != nil {
		log.Println(err)
	}

	return results, nil
}

func newfileUploadRequest(uri string, paramName, path string) (*retryablehttp.Request, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"comment":   "file uploaded using markdown-github-action",
		"minorEdit": "true",
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("PUT", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

// UploadAttachment to a page identify by page ID
// you need the page ID to upload the attachment(file path)
func (a *APIClient) UploadAttachment(filename string, id int) error {
	targetURL := fmt.Sprintf(common.ConfluenceBaseURL+"/wiki/rest/api/content/%d/child/attachment", id)

	req, err := newfileUploadRequest(targetURL, "file", filename)
	if err != nil {
		return err
	}

	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("Accept", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload attachment: %s", resp.Status)
	}

	func() { _ = resp.Body.Close() }()

	return nil
}
