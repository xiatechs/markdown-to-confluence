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
	contents, err := ioutil.ReadFile(filepath.Clean(fpath))
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	fullpath := strings.Replace(fpath, ".", "", 2)

	fullpath = strings.TrimPrefix(fullpath, "/")

	todo.ParseGo(contents, fullpath)

	return nil
}

// processMarkDown method takes in file contents
// and parses the markdown file before calling
// checkConfluencePages method
func (node *Node) processMarkDown(path string) error {
	contents, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}

	parsedContents, err := markdown.ParseMarkdown(node.root.id, contents)
	if err != nil {
		return fmt.Errorf("parse markdown error: %w", err)
	}

	err = node.checkConfluencePages(parsedContents)
	if err != nil {
		log.Printf("error completing confluence operations: %s", err)
	}

	return nil
}

// uploadFile method takes in file and
// uploads the file to a page by parent page ID (node.root.id)
func (node *Node) uploadFile(path string) {
	err := NodeAPIClient.UploadAttachment(filepath.Clean(path), node.root.id)
	if err != nil {
		log.Printf("error uploading attachment: %s", err)
	}
}
