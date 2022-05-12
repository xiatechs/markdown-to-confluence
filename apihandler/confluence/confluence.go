// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/common"
	fh "github.com/xiatechs/markdown-to-confluence/filehandler"
)

// newPageResults function takes in a http response and
// decodes the response body into a PageResults struct that is returned
func newPageResults(resp *http.Response) (*PageResults, error) {
	pageResultVar := PageResults{}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("newPageResults statuscode error - not good request: status code [%d]", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("newPageResults decode error: %w", err)
	}

	err = json.Unmarshal(contents, &pageResultVar)
	if err != nil {
		return nil, fmt.Errorf("newPageresults json unmarshal error: %w", err)
	}

	if len(pageResultVar.Results) == 0 { // we want to return nil to skip this result
		return nil, nil
	}

	return &pageResultVar, nil
}

// grabPageContents method takes in page contents, the parent page id (root) and boolean
// confirming whether or not the page is a parent page folder (isroot)
// and returns byte of page contents
func (a *APIClient) grabPageContents(contents *fh.FileContents, root int, isroot bool) ([]byte, error) {
	title, ok := contents.MetaData["title"]

	if !ok {
		return nil, fmt.Errorf("grabPageContents err - title is empty")
	}

	newPageContent := Page{
		Type:  "page",
		Title: title.(string),
		Space: SpaceObj{Key: a.Space},
		Body: BodyObj{Storage: StorageObj{
			Value:          string(contents.Body),
			Representation: "storage",
		}},
	}

	if !isroot {
		newPageContent.Ancestors = []AncestorObj{{ID: root}}
	}

	newPageContentsJSON, err := json.Marshal(newPageContent)
	if err != nil {
		return nil, fmt.Errorf("grabPageContents failed to marshal new page contents: %w", err)
	}

	return newPageContentsJSON, nil
}

// CreatePage method takes root (root page id) and page contents and bool (is page root?)
// and generates a page in confluence and returns the generated page ID
//nolint: gocyclo // 11 is just about fine
func (a *APIClient) CreatePage(root int, contents *fh.FileContents, isroot bool) (int, error) {
	if contents == nil {
		return 0, fmt.Errorf("createpage error: contents parameter is nil")
	}

	title, ok := contents.MetaData["title"]
	if ok {
		log.Printf("start creating page with title [%s]", title.(string))
	}

	newPageContentsJSON, err := a.grabPageContents(contents, root, isroot)
	if err != nil {
		return 0, fmt.Errorf("createpage error: %w", err)
	}

	if newPageContentsJSON == nil {
		return 0, fmt.Errorf("createpage error: newPageContentsJSON is nil")
	}

	URL := fmt.Sprintf("%s/wiki/rest/api/content", a.BaseURL)

	req, err := retryablehttp.NewRequest(http.MethodPost, URL, newPageContentsJSON)
	if err != nil {
		return 0, fmt.Errorf("createpage error: %w", err)
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("body close error: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to create confluence page: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)

	var output struct {
		ID string `json:"id"`
	}

	err = decoder.Decode(&output)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(output.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// updatePageContents method updates the page contents and return as a []byte JSON to be used
func (a *APIClient) updatePageContents(pageVersion int64, contents *fh.FileContents) ([]byte, *Page, error) {
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
		return nil, nil, fmt.Errorf("json marshal error: %w", err)
	}

	return newPageContentsJSON, &newPageContent, nil
}

// DeletePage deletes a confluence page by page ID
func (a *APIClient) DeletePage(pageID int) error {
	URL := fmt.Sprintf("%s/wiki/rest/api/content/%d", a.BaseURL, pageID)

	req, err := retryablehttp.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		return fmt.Errorf("deletepage error: %w", err)
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("deletepage error: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("body close error: %w", err))
		}
	}()

	if err != nil {
		log.Println("deletepage error", resp.Status, err)
	}

	return nil
}

// UpdatePage updates a confluence page with our newly created data and increases the
// version by 1 each time.
func (a *APIClient) UpdatePage(pageID int, pageVersion int64, pageContents *fh.FileContents,
	originalPage PageResults) (bool, error) {
	newPageContentsJSON, newPage, err := a.updatePageContents(pageVersion, pageContents)
	if err != nil {
		return false, fmt.Errorf("updatePageContents error: %w", err)
	}

	if len(originalPage.Results) > 0 {
		if originalPage.Results[0].Body.Storage == newPage.Body.Storage {
			return true, nil
		}
	}

	URL := fmt.Sprintf("%s/wiki/rest/api/content/%d", a.BaseURL, pageID)

	req, err := retryablehttp.NewRequest(http.MethodPut, URL, bytes.NewBuffer(newPageContentsJSON))
	if err != nil {
		return false, fmt.Errorf("updatePageContents error: %w", err)
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("updatepage failed to do the request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("body close error: %w", err))
		}
	}()

	if err != nil {
		log.Println("updatepage error was: ", resp.Status, err)
	}

	return true, nil
}

// createFindPageRequest method takes in a title (page title) and searches for page
// in confluence
func (a *APIClient) createFindPageRequest(title string) (*retryablehttp.Request, error) {
	lookUpURL := fmt.Sprintf("%s/wiki/rest/api/content?expand=body.storage,version&type=page&spaceKey=%s&title=%s",
		a.BaseURL, a.Space, title)

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil {
		return nil, fmt.Errorf("createFindPageRequest error: %w", err)
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
		return nil, fmt.Errorf("createFindPagesRequest error: %w", err)
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
			return nil, fmt.Errorf("createFindPagesRequest error: %w", err)
		}
	} else {
		req, err = a.createFindPageRequest(title)
		if err != nil {
			return nil, fmt.Errorf("createFindPageRequest error: %w", err)
		}
	}

	return req, nil
}

// FindPage in confluence
// Docs for this API endpoint are here
// https://developer.atlassian.com/cloud/confluence/rest/api-group-content/#api-api-content-get
func (a *APIClient) FindPage(title string, many bool) (*PageResults, error) {
	req, err := a.findPageRequest(title, many)
	if err != nil {
		return nil, fmt.Errorf("find page request error: %w", err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("find page request error: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("body close error: %w", err))
		}
	}()

	results, err := newPageResults(resp)
	if err != nil {
		return nil, err
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
			log.Println(fmt.Errorf("file close error: %w", err))
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("create form file error: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("io copy error: %w", err)
	}

	params := map[string]string{
		"comment":   "file uploaded using fh-github-action",
		"minorEdit": "true",
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	req, err := retryablehttp.NewRequest("PUT", uri, body)
	if err != nil {
		return nil, fmt.Errorf("writer close error: %w", err)
	}

	data := []byte{}

	_, err = file.Read(data)
	if err != nil {
		log.Println(err)
	}

	sEnc := base64.StdEncoding.EncodeToString(data)
	req.ContentLength = int64(len(sEnc))

	req.Header = http.Header{
		"Content-Type": []string{writer.FormDataContentType()},
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("writer close error: %w", err)
	}

	return req, err
}

// UploadAttachment to a page identify by page ID
// you need the page ID to upload the attachment(file path)
func (a *APIClient) UploadAttachment(filename string, id int) error {
	var targetURL string

	targetURL = fmt.Sprintf(common.ConfluenceBaseURL+"/wiki/rest/api/content/%d/child/attachment", id)

	req, err := newfileUploadRequest(targetURL, "file", filename)
	if err != nil {
		return fmt.Errorf("file upload error: %w", err)
	}

	req.SetBasicAuth(a.Username, a.Password)
	req.Header.Set("Accept", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("upload attachment response error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload attachment response issue: %s", resp.Status)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(fmt.Errorf("body close error: %w", err))
		}
	}()

	return nil
}
