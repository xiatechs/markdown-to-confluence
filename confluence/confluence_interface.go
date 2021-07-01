package confluence

import "github.com/xiatechs/markdown-to-confluence/markdown"

//go:generate mockgen -destination=./apiclient_mock_test.go -package=confluence -source=confluence_interface.go
type APIClienter interface {
	CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error)
	DeletePage(pageID int) error
	UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
		originalPage PageResults) (bool, error)
	FindPage(title string, many bool) (*PageResults, error)
	UploadAttachment(filename string, id int) error
}
