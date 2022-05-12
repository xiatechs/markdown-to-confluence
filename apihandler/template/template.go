package template

import (
	"fmt"
	"log"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var mockPageIDgenerator = 0

type Example struct {
}

func (l *Example) Delete(titles map[string]struct{}, id string) {
	// logic here will be:

	// delete any pages on website where the title is not in the 'titles' map.
}

// CRUD - Create, Update - do other stuff (process files locally? generate docs for go files? etc)
// - if exist: update, if not: create, else delete.
func (e *Example) CRUD(file *filehandler.FileContents, parentMetaData map[string]interface{}) (map[string]interface{}, error) {
	if file == nil { // this means the filehandler step returned a nil file somehow
		return nil, fmt.Errorf("a nil file was passed to the API")
	}

	state := common.CaptureState(file, parentMetaData)

	/////////////////////////////////////////////////// check if the file already exists
	//pageResults, err := e.client.FindPage(file *filehandler.FileContents)

	/////////////////////////////////////////////////// upload an attachment or create / update a file
	switch state.FileType {
	case "image":
		//id, err := e.client.UploadImage(file *filehandler.FileContents, parentMetaData)
		//state.OutputPageID = id
	case "markdown":
		//id, err := e.client.CreateOrUpdateMarkdown(file *filehandler.FileContents, parentMetaData)
		//state.OutputPageID = id
	case "folderpage":
		//id, err := e.client.CreateOrUpdateFolderPage(file *filehandler.FileContents, parentMetaData)
		//state.OutputPageID = id
	}

	mockPageIDgenerator++

	state.OutputPageID = mockPageIDgenerator

	log.Println(state.FileType, state)

	output := returnOutputs(state, file)

	return output, nil
}

func returnOutputs(state *common.FileState, file *filehandler.FileContents) map[string]interface{} {
	output := make(map[string]interface{})
	output["title"] = file.MetaData["title"]
	output["id"] = state.OutputPageID

	return output
}
