package node

import (
	"log"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

/*
	CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error)
	DeletePage(pageID int) error
	UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
		originalPage confluence.PageResults) (bool, error)
	FindPage(title string, many bool) (*confluence.PageResults, error)
	UploadAttachment(filename string, id int, index bool, indexid int) error
*/

var mockiter = 0
var toproot bool

type mockclient struct {
}

func (m mockclient) CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error) {
	mockiter++

	if !toproot {
		isroot = true
		toproot = true
	}

	log.Printf("CREATING PAGE:\n%s\nroot [%d]\nisRoot [%t]\n ", string(contents.Body), root, isroot)

	return mockiter, nil
}

func (m mockclient) DeletePage(pageID int) error {
	return nil
}

func (m mockclient) UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
	originalPage confluence.PageResults) (bool, error) {
	log.Println("UPDATING PAGE")
	return true, nil
}

func (m mockclient) FindPage(title string, many bool) (*confluence.PageResults, error) {
	return nil, nil
}

func (m mockclient) UploadAttachment(filename string, id int, index bool, indexid int) error {
	log.Printf("UPLOADING: name:[%s], id:[%d], index:[%t], indexID:[%d]", filename, id, index, indexid)
	return nil
}
