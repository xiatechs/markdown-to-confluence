package confluence

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var Iterator = 0

type Local struct {
	client *APIClient
	mu     *sync.RWMutex
}

func NewAPIClient() *Local {
	client, err := CreateAPIClient()
	if err != nil {
		log.Fatal(err) // just die here
	}

	l := &Local{
		client: client,
		mu:     &sync.RWMutex{},
	}

	return l
}

func (l *Local) Delete(titles map[string]struct{}, id string) {
	children, err := l.client.FindPage(id, true)
	if err != nil {
		log.Printf("error finding page: %s", err)
	}

	if children != nil {
		l.mu.RLock()
		defer l.mu.RUnlock()
		for index := range children.Results {
			log.Println(children.Results[index].Title)
			if _, ok := titles[children.Results[index].Title]; ok {
				l.Delete(titles, children.Results[index].ID)
			} else {
				go func() {
					pageID, err := strconv.Atoi(children.Results[index].ID)
					if err != nil {
						log.Printf("error getting page ID: %s", err)
						return
					}
					err = l.client.DeletePage(pageID)
					if err != nil {
						log.Printf("error deleting page: %s", err)
					}
				}()
			}
		}
	}
}

// CRUD - Create, Update - do other stuff (process files locally? generate docs for go files? etc)
// - if exist: update, if not: create, else delete.
func (l *Local) CRUD(file *filehandler.FileContents, parentMetaData map[string]interface{}) (map[string]interface{}, error) {
	if file == nil { // this means the filehandler step returned a nil file somehow
		return nil, fmt.Errorf("a nil file was passed to the API")
	}

	state := common.CaptureState(file, parentMetaData)

	output := make(map[string]interface{})

	/////////////////////////////////////////////////// check if the file already exists
	pageResults, err := l.client.FindPage(state.CurrentPageTitle, false)
	if err != nil {
		return nil, err
	}

	/////////////////////////////////////////////////// upload an attachment or create / update a file
	switch state.FileType {
	case "image":
		err = l.handleImage(pageResults, state, file, parentMetaData)
		if err != nil {
			return nil, err
		}
	case "markdown":
		err = l.handleMarkDownFile(pageResults, state, file, parentMetaData)
		if err != nil {
			return nil, err
		}
	case "indexPage":
		err = l.handleMarkDownFile(pageResults, state, file, parentMetaData)
		if err != nil {
			return nil, err
		}
	case "folderpage":
		err = l.handleFolderPage(pageResults, state, file, parentMetaData)
		if err != nil {
			return nil, err
		}
	}

	output["title"] = file.MetaData["title"]

	if pageResults == nil {
		output["id"] = state.OutputPageID
		return output, nil
	}

	output["id"], err = strconv.Atoi(pageResults.Results[0].ID)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (l *Local) handleImage(pageResults *PageResults,
	state *common.FileState,
	file *filehandler.FileContents,
	parentMetaData map[string]interface{}) error {
	var err error

	log.Println("image", state.FilePath,
		state.CurrentPageID, state.IsIndexPage, state.ParentPageID)

	if pageResults != nil {
		state.ParentPageID, _ = strconv.Atoi(pageResults.Results[0].ID)
	}

	err = l.client.UploadAttachment(state.FilePath, state.ParentPageID)
	if err != nil {
		return err
	}

	return nil
}

func (l *Local) handleFolderPage(pageResults *PageResults,
	state *common.FileState,
	file *filehandler.FileContents,
	parentMetaData map[string]interface{}) error {
	log.Println("folderpage", state.ParentPageID, file.MetaData, state.IsRoot)
	var err error

	if pageResults == nil {
		state.OutputPageID, err = l.client.CreatePage(state.ParentPageID,
			file, state.IsRoot)
		if err != nil {
			return err
		}

		return nil
	}

	noChanges, err := l.client.UpdatePage(state.CurrentPageID,
		int64(pageResults.Results[0].Version.Number),
		file, *pageResults)
	if err != nil {
		return err
	}

	if noChanges {
		log.Println("No changes to this page ID:", state.CurrentPageTitle)
	}

	state.OutputPageID = state.CurrentPageID

	return nil
}

func (l *Local) handleMarkDownFile(pageResults *PageResults,
	state *common.FileState,
	file *filehandler.FileContents,
	parentMetaData map[string]interface{}) error {
	log.Println("markdown", state.ParentPageID,
		file.MetaData, state.IsRoot, state.CurrentPageID)
	var err error

	if pageResults == nil {
		state.OutputPageID, err = l.client.CreatePage(state.ParentPageID,
			file, state.IsRoot)
		if err != nil {
			return err
		}

		return nil
	}

	if pageResults != nil {
		state.CurrentPageID, _ = strconv.Atoi(pageResults.Results[0].ID)
	}

	noChanges, err := l.client.UpdatePage(state.CurrentPageID,
		int64(pageResults.Results[0].Version.Number), file, *pageResults)
	if err != nil {
		return err
	}

	if noChanges {
		log.Println("No changes to this page ID:", state.CurrentPageTitle)
	}

	state.OutputPageID = state.CurrentPageID

	return nil
}
