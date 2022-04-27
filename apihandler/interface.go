package apihandler

import (
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

type ApiController interface { // API just takes in public CRUD method - from there you can create an API to deal with file contents
	CRUD(file *filehandler.FileContents,
		parentMetaData map[string]interface{}) (map[string]interface{}, error)
}
