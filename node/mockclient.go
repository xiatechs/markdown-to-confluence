package node

import (
	"log"
	"sync"

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

type mockclient struct {
	mockiter int
	toproot  bool
	mu       *sync.RWMutex // for locking/unlocking when multiple goroutines are working on same node
}

func (m mockclient) CreatePage(root int, contents *markdown.FileContents, isroot bool) (int, error) {
	m.mu.Lock()
	m.mockiter++

	if !m.toproot {
		isroot = true
		m.toproot = true
	}

	log.Printf("CREATING PAGE:\n%s\nroot [%d]\nisRoot [%t]\n ", string(contents.Body), root, isroot)
	m.mu.Unlock()
	return m.mockiter, nil
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
	m.mu.Lock()
	log.Printf("UPLOADING: name:[%s], id:[%d], index:[%t], indexID:[%d]", filename, id, index, indexid)
	m.mu.Unlock()
	return nil
}
