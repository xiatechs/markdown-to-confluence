package node

import (
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

//go:generate mockgen -destination=./apiclient_mock_test.go -package=node -source=client_interface.go

// APIClienter is interface for confluence API client and mock tests
type APIClienter interface {
	/* CreatePage - create a page

	root = the root page ID
	contents = the contents of the page that we are looking to upload to the wiki
	isroot = bool, if it's the root page i.e the top page this is 'true'.
	*/
	CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error)

	/* DeletePage - delete a page

	root = the page ID (the page you want to delete)
	*/
	DeletePage(pageID int) error

	/* UpdatePage - update a page

	originalPage confluence.PageResults
	pageID = the page ID that you are updating
	pageVersion = the version of the page i.e if you update the page, update the version
	pageContents = the new page contents you're going to upload / update
	originalPage = the original page that's been collected
	*/
	UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
		originalPage confluence.PageResults) (bool, error)

	/* FindPage(title string, many bool) (*confluence.PageResults, error)

	FindPage - find a page
	title = the title of the page
	many = if we are collecting many pages i.e the parent page and all it's child pages
	*/
	FindPage(title string, many bool) (*confluence.PageResults, error)

	/* UploadAttachment - upload an attachment to the client

	filename = the name of the file
	id = the page that the file is being uploaded to
	index = is the page you're uploading the file to the index page?
	indexid = the id of the index page
	*/
	UploadAttachment(filename string, id int, index bool, indexid int) error
}
