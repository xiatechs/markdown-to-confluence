package node

import (
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

//go:generate mockgen -destination=./apiclient_mock_test.go -package=node -source=mock_client.go

type APIClienter interface {
	CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error)
	DeletePage(pageID int) error
	UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
		originalPage confluence.PageResults) (bool, error)
	FindPage(title string, many bool) (*confluence.PageResults, error)
	UploadAttachment(filename string, id int) error
}
