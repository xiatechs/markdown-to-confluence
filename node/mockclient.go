package node

import (
	"fmt"
	"sort"

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
	pages    []*Page
}

// Page - test pages are generated and stored in here
type Page struct {
	title  string
	isroot bool
	body   string
	root   int
	id     int
}

type mockclient struct {
	i *iterator
}

var s = make(chan bool, 1)

func (m mockclient) Print() {
	sort.Slice(m.i.pages, func(i, j int) bool {
		return m.i.pages[i].id > m.i.pages[j].id
	})

	for _, page := range m.i.pages {
		fmt.Printf("-----------------------------------------\n")
		fmt.Printf("PAGE TITLE %s:\n\n", page.title)
		fmt.Println(page.body)
		fmt.Println("\nPAGE DETAILS:")
		fmt.Printf("\ntop: %t, id %d, root: %d\n", page.isroot, page.id, page.root)
		fmt.Printf("-----------------------------------------\n\n")
	}
}

var mocksem = make(chan bool, 1)

func (m mockclient) GetPages() []Page {
	mocksem <- true
	pages := []Page{}
	for _, page := range m.i.pages {
		pages = append(pages, *page)
	}
	<-mocksem
	return pages
}

//nolint: ineffassign // is ok
func (i *iterator) append(root int, contents *markdown.FileContents, isroot bool) int {
	var exists bool

	for index := range i.pages {
		if i.pages[index].title == contents.MetaData["title"].(string) {
			exists = true

			i.pages[index].body = string(contents.Body)

			return i.pages[index].id
		}
	}

	if !exists {
		i.pages = append(i.pages, &Page{
			title:  contents.MetaData["title"].(string),
			body:   string(contents.Body),
			root:   root,
			isroot: isroot,
			id:     i.mockiter,
		})
	}

	i.mockiter++

	return i.mockiter
}

func (m mockclient) CreatePage(root int, contents *markdown.FileContents, _ bool) (int, error) {
	s <- true // race blocker

	var isroot bool

	if !m.i.isroot {
		m.i.isroot = true

		isroot = true
	} else {
		isroot = false
	}

	id := m.i.append(root, contents, isroot)

	<-s

	return id, nil
}

func (m mockclient) DeletePage(pageID int) error {
	return nil
}

func (m mockclient) UpdatePage(pageID int, pageVersion int64, pageContents *markdown.FileContents,
	originalPage confluence.PageResults) (bool, error) {
	s <- true // race blocker

	var isroot bool

	if !m.i.isroot {
		m.i.isroot = true

		isroot = true
	} else {
		isroot = false
	}

	_ = m.i.append(pageID, pageContents, isroot)

	<-s

	return true, nil
}

func (m mockclient) FindPage(title string, many bool) (*confluence.PageResults, error) {
	return nil, nil
}

func (m mockclient) UploadAttachment(filename string, id int, index bool, indexid int) error {
	return nil
}
