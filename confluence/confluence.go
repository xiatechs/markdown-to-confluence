// Package confluence provides functionality for interacting with the confluence APIClient
// Specifically managing pages
package confluence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/markdown"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


// CreatePage in confluence
func (a *APIClient) CreatePage() error {
	return nil
}

// UpdatePage updates a confluence page with our newly created data and increases the
// version by 1 each time.
func (a *APIClient) UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents) error {

	fmt.Println("running update now....") //todo: remove
	newPageJson := PutPageContent{
		Type:  "page",
		Title: pageContents.MetaData["title"].(string),
		Version: VersionObj{
			Number: int(pageVersion) + 1,
		},
		Body: BodyObj{
			Storage: StorageObj{
				Value:          string(pageContents.Body),
				Representation: "storage",
			},
		},
	}

	URL := fmt.Sprintf("%s/wiki/rest/api/content/%d", a.BaseURL, pageID)

	b, err := json.Marshal(newPageJson)
	if err != nil {
		return err
	}

	fmt.Println(string(b)) //todo:remove

	req, err := retryablehttp.NewRequest(http.MethodPut, URL, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return err
	}

	req.SetBasicAuth(a.Username, a.Password)

	req.Header.Set("Content-Type", "application/json")

	fmt.Println("request:   ", req)

	resp, err := a.Client.Do(req)
	if err != nil {
		fmt.Println("error was: ", resp.Status, err)
		return fmt.Errorf("failed to do the request: %w", err)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error ioutil", err)
	}
	fmt.Println("response: ", string(r))
	defer func() { _ = resp.Body.Close() }()

	return nil
}

// FindPage in confluence
// Docs for this API endpoint are here
// https://developer.atlassian.com/cloud/confluence/rest/api-group-content/#api-api-content-get
func (a *APIClient) FindPage(title string) (int, int64, bool, error) {
	fmt.Printf("%+v", a) //todo remove
	lookUpURL := fmt.Sprintf("%s/wiki/rest/api/content?expand=version&type=page&spaceKey=%s&title=%s", a.BaseURL, a.Space, title)

	req, err := retryablehttp.NewRequest(http.MethodGet, lookUpURL, nil)
	if err != nil {
		return 0, 0, false, err
	}

	req.SetBasicAuth(a.Username, a.Password)

	resp, err := a.Client.Do(req)
	if err != nil {
		return 0, 0, false, fmt.Errorf("failed to do the request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	fmt.Println("req: ", req, lookUpURL) //todo remove

	r := findPageResult{}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, 0, false, err
	}

	fmt.Println("response decoded: ", r.Results) //todo remove
	spew.Dump(r)                                 //todo remove

	if len(r.Results) == 0 {
		return 0, 0, false, fmt.Errorf("no page present")
	}

	pageId, err := strconv.Atoi(r.Results[0].ID)
	if err != nil {
		fmt.Errorf("error converting ID to int value")
		return 0, 0, false, err
	}
	return pageId, r.Results[0].Version.Number, true, nil
}
