package apihandler

import (
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

//go:generate mockgen -destination=./api_mocks.go -package=apihandler -source=interface.go

/*

inside CRUD use the method use common.CaptureState() i.e

state := common.CaptureState(file, parentMetaData)

and you'll be returned:

// FileState - during the generation of files, these fields refer to different states of files
type FileState struct {
	CurrentPageID    int
	OutputPageID     int
	ParentPageID     int
	CurrentPageTitle string
	IsRoot           bool
	IsIndexPage      bool
	FileType         string
	FilePath         string
	Alive            bool
}

*/

// ApiController - interface for the API
type ApiController interface {
	CRUD(file *filehandler.FileContents,
		parentMetaData map[string]interface{}) (map[string]interface{}, error)
	Delete(titles map[string]struct{}, id string)
}
