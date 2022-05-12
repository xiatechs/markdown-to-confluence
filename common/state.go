package common

import "github.com/xiatechs/markdown-to-confluence/filehandler"

// FileState - during the generation of files, these fields refer to different states of files

// TODO: remove maps of interfaces with just this filestate so that the API is simpler to understand

type FileState struct {
	CurrentPageID    int    // the currentpage ID - if it already has an ID & has been created before - it'll be stored here
	OutputPageID     int    // the OutputPageID
	ParentPageID     int    // the parent page ID - if this page has a parent page, this is it's ID
	CurrentPageTitle string // the page title of the page you are creating
	IsRoot           bool   // is this the root page? true = yes
	IsIndexPage      bool   // is this an 'index page' (a page that will be used as an index that is a readme.md)
	FileType         string // type of file (current types are: indexPage, markdown, image, folderpage)
	FilePath         string // the file path of the file in question
	Alive            bool   // if this is a folderpage - is it alive i.e does it have markdown contained in it
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

	fileState.FilePath = func() string {
		if value, ok := file.MetaData["filepath"].(string); ok { // what is the parent page ID?
			return value
		}

		return ""
	}()

	return fileState
}
