package node

// methods for processing/reading/uploading files & iterating through folders

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/markdown"
	"github.com/xiatechs/markdown-to-confluence/todo"
)

// processGoFile method takes in a go file contents and
// calls method todo.ParseGo on the file contents with the
// file path
func (node *Node) processGoFile(fpath string) error {
	_, abs := node.generateTitles()

	contents, err := ioutil.ReadFile(filepath.Clean(fpath))
	if err != nil {
		return fmt.Errorf("absolute path [%s] - file [%s] - read file error: %w",
			abs, fpath, err)
	}

	fullpath := strings.Replace(fpath, ".", "", 2)

	fullpath = strings.TrimPrefix(fullpath, "/")

	todo.ParseGo(contents, fullpath) // not used atm but will be maybe in future

	return nil
}

// processMarkDownIndex method takes in index file contents
// and parses the markdown file
func (node *Node) processMarkDownIndex(path string) (*markdown.FileContents, error) {
	fpath, abs := node.generateTitles()

	contents, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("absolute path [%s] - file [%s] - read file error: %w",
			abs, path, err)
	}

	parsedContents, err := markdown.ParseMarkdown(func() int {
		if node.root == nil {
			return 0
		}

		return node.root.id
	}(), contents, node.indexPage, node.id)
	if err != nil {
		return nil, fmt.Errorf("absolute path [%s] - file [%s] - parse markdown error: %w",
			abs, path, err)
	}

	parsedContents.MetaData["title"] = fpath

	return parsedContents, nil
}

// processMarkDown method takes in file contents
// and parses the markdown file before calling
// checkConfluencePages method
func (node *Node) processMarkDown(path string) error {
	fpath, abs := node.generateTitles()

	contents, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("absolute path [%s] - file [%s] - read file error: %w",
			abs, path, err)
	}

	parsedContents, err := markdown.ParseMarkdown(func() int {
		if node.root == nil {
			return 0
		}

		return node.root.id
	}(), contents, node.indexPage, node.id)
	if err != nil {
		return fmt.Errorf("absolute path [%s] - file [%s] - parse markdown error: %w",
			abs, path, err)
	}

	parsedContents.MetaData["title"] = parsedContents.MetaData["title"].(string) + "-" + fpath

	err = node.checkConfluencePages(parsedContents)
	if err != nil {
		return fmt.Errorf("absolute path [%s] - file [%s] - confluence check error: %w",
			abs, path, err)
	}

	return nil
}

// uploadFile method takes in file and
// uploads the file to a page by parent page ID (node.root.id)
func (node *Node) uploadFile(path string, isIndexPage bool) {
	_, abs := node.generateTitles()

	err := nodeAPIClient.UploadAttachment(filepath.Clean(path), node.root.id, isIndexPage, node.id)
	if err != nil {
		log.Printf("absolute path [%s] - local path [%s] - file upload error: %v",
			path, abs, err)
	}
}
