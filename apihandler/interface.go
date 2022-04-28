package apihandler

import (
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

//go:generate mockgen -destination=./api_mocks.go -package=apihandler -source=interface.go

/* there's three areas where you can interface with metadata between the different layers:
file.MetaData - this can contain data about the file but also the current state
parentMetaData - this contains data about the parent i.e the parent ID
returned MetaData (map[string]interface{}) - this data is captured by the foldercrawler i.e the returned ID from page that is created
*/

// ApiController - interface for the API
type ApiController interface { // API just takes in public CRUD method - from there you can create an API to deal with file contents
	CRUD(file *filehandler.FileContents,
		parentMetaData map[string]interface{}) (map[string]interface{}, error)
}
