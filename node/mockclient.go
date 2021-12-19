package node

import (
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
type iterator struct { // enables pointer arithmetic
	mockiter int
	isroot   bool
}

type mockclient struct {
	i *iterator
}

var s = make(chan bool, 1)

//nolint: staticcheck,ineffassign // is fine
func (m mockclient) CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error) {
	s <- true // race blocker

	m.i.mockiter++

	a := m.i.mockiter

	if !m.i.isroot {
		m.i.isroot = true

		isroot = true
	}

	<-s

	return a, nil
}

func (m mockclient) DeletePage(pageID int) error {
	return nil
}

func (m mockclient) UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
	originalPage confluence.PageResults) (bool, error) {
	return true, nil
}

func (m mockclient) FindPage(title string, many bool) (*confluence.PageResults, error) {
	return nil, nil
}

func (m mockclient) UploadAttachment(filename string, id int, index bool, indexid int) error {
	return nil
}
