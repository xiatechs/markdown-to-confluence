package confluence

import (
	"fmt"
	"log"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var Iterator = 0

type Local struct {
	client *APIClient
}

func NewAPIClient() *Local {
	client, err := CreateAPIClient()
	if err != nil {
		log.Fatal(err) // just die here
	}

	l := &Local{
		client: client,
	}

	return l
}

// CRUD - Create, Update, Delete - do other stuff (process files locally? generate docs for go files? etc)
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
		log.Println("image", state.CurrentPageTitle,
			state.CurrentPageID, state.IsIndexPage, state.ParentPageID)
		err = l.client.UploadAttachment(state.CurrentPageTitle,
			state.CurrentPageID, state.IsIndexPage, state.ParentPageID)
		if err != nil {
			return nil, err
		}
	case "markdown":
		log.Println("markdown", state.ParentPageID,
			file.MetaData, state.IsRoot)
		if pageResults == nil {
			state.OutputPageID, err = l.client.CreatePage(state.ParentPageID,
				file, state.IsRoot)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := l.client.UpdatePage(state.CurrentPageID,
				int64(pageResults.Results[0].Version.Number), file, *pageResults)
			if err != nil {
				return nil, err
			}

			state.OutputPageID = state.CurrentPageID
		}
	case "folderpage":
		log.Println("folderpage", state.ParentPageID, file.MetaData, state.IsRoot)
		if pageResults == nil {
			state.OutputPageID, err = l.client.CreatePage(state.ParentPageID,
				file, state.IsRoot)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := l.client.UpdatePage(state.CurrentPageID,
				int64(pageResults.Results[0].Version.Number),
				file, *pageResults)
			if err != nil {
				return nil, err
			}

			state.OutputPageID = state.CurrentPageID
		}
	}

	output["id"] = state.OutputPageID

	return output, nil
}
