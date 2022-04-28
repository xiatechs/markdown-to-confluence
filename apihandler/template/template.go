package template

import (
	"fmt"
	"log"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

type Example struct {
	// client HALOclient
}

var iter = 0

// CRUD - Create, Update, Delete
// - if exist: update, if not: create, else delete.
func (e *Example) CRUD(file *filehandler.FileContents, parentMetaData map[string]interface{}) (map[string]interface{}, error) {
	if file == nil { // this means the filehandler step returned a nil file somehow
		return nil, fmt.Errorf("a nil file was passed to the API")
	}

	state := common.CaptureState(file, parentMetaData)
	/*
		type FileState struct {
			CurrentPageID    int // the page ID of the current page (if it's an INDEX page)
			ParentPageID     int // the page ID of the parent page
			OutputPageID     int // the page ID of created page
			CurrentPageTitle string // the page title of current page
			IsRoot           bool // if this the root page
			IsIndexPage      bool // is this an INDEX page (an index page is a page created for a folder using the README.MD file)
		}
	*/
	iter++

	state.OutputPageID = iter

	log.Println(state)

	output := make(map[string]interface{})

	output["id"] = iter

	return output, nil
}
