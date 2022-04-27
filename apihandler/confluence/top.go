package confluence

import (
	"log"

	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var Iterator = 0

type Local struct {
	client *APIClient
}

func NewAPIClient() *Local {
	client, err := CreateAPIClient()
	if err != nil {
		log.Println(err)
		return nil
	}

	l := &Local{
		client: client,
	}

	return l
}

// CRUD - Create, Update, Delete
// - if exist: update, if not: create, else delete.
func (l *Local) CRUD(file *filehandler.FileContents, parentMetaData map[string]interface{}) (map[string]interface{}, error) {
	if file == nil { // this means the filehandler step returned a nil file somehow
		return nil, nil
	}

	// variables captured via metadata being passed from file
	var (
		//input
		err           error
		currentPageID = func() int {
			if value, ok := file.MetaData["id"].(int); ok {
				return value
			}

			return 0
		}()
		currentPageTitle = func() string {
			if str, ok := file.MetaData["title"].(string); ok { // what is the current page title
				return str
			}
			return ""
		}()
		isRoot = func() bool {
			if b, ok := file.MetaData["root"].(bool); ok { // is this the root page?
				return b
			}

			return false
		}()
		isIndexPage = func() bool {
			if b, ok := file.MetaData["indexPage"].(bool); ok { // is this an 'index' page i.e a README.md
				return b
			}

			return false
		}()
		parentPageID = func() int {
			if value, ok := parentMetaData["id"].(int); ok { // what is the parent page ID?
				return value
			}

			return 0
		}()
		outputPageID int
	)

	output := make(map[string]interface{})

	pageResults, err := l.client.FindPage(currentPageTitle, false)
	if err != nil {
		return nil, err
	}

	switch file.MetaData["type"] {
	case "attachment":
		log.Println("attachment", currentPageTitle, currentPageID, isIndexPage, parentPageID)
		err = l.client.UploadAttachment(currentPageTitle, currentPageID, isIndexPage, parentPageID)
		if err != nil {
			return nil, err
		}
	case "markdown":
		log.Println("markdown", parentPageID, file.MetaData, isRoot)
		if pageResults == nil {
			outputPageID, err = l.client.CreatePage(parentPageID, file, isRoot)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := l.client.UpdatePage(currentPageID,
				int64(pageResults.Results[0].Version.Number), file, *pageResults)
			if err != nil {
				return nil, err
			}

			outputPageID = currentPageID
		}
		//Iterator++
	case "folderpage":
		log.Println("folderpage", parentPageID, file.MetaData, isRoot)
		if pageResults == nil {
			outputPageID, err = l.client.CreatePage(parentPageID, file, isRoot)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := l.client.UpdatePage(currentPageID,
				int64(pageResults.Results[0].Version.Number), file, *pageResults)
			if err != nil {
				return nil, err
			}

			outputPageID = currentPageID
		}
	}

	output["id"] = outputPageID

	return output, nil
}
