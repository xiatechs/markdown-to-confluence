package common

import "github.com/xiatechs/markdown-to-confluence/filehandler"

// FileState - during the generation of files, these fields refer to different states of files
type FileState struct {
	CurrentPageID    int
	OutputPageID     int
	ParentPageID     int
	CurrentPageTitle string
	IsRoot           bool
	IsIndexPage      bool
	FileType         string
	Alive            bool
}

func CaptureState(file *filehandler.FileContents, parentMetaData map[string]interface{}) *FileState {
	fileState := &FileState{}

	fileState.CurrentPageID = func() int {
		if value, ok := file.MetaData["id"].(int); ok {
			return value
		}

		return 0
	}()

	fileState.CurrentPageTitle = func() string {
		if str, ok := file.MetaData["title"].(string); ok { // what is the current page title
			return str
		}
		return ""
	}()

	fileState.IsRoot = func() bool {
		if b, ok := file.MetaData["root"].(bool); ok { // is this the root page?
			return b
		}

		return false
	}()

	fileState.IsIndexPage = func() bool {
		if b, ok := file.MetaData["indexPage"].(bool); ok { // is this an 'index' page i.e a README.md
			return b
		}

		return false
	}()

	fileState.ParentPageID = func() int {
		if value, ok := parentMetaData["id"].(int); ok { // what is the parent page ID?
			return value
		}

		return 0
	}()

	fileState.FileType = func() string {
		if value, ok := file.MetaData["type"].(string); ok { // what is the parent page ID?
			return value
		}

		return ""
	}()

	fileState.Alive = func() bool {
		if value, ok := file.MetaData["alive"].(bool); ok { // what is the parent page ID?
			return value
		}

		return false
	}()

	return fileState
}
