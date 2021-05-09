package node

// methods for processing/reading/uploading files & iterating through folders

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/markdown"
	"github.com/xiatechs/markdown-to-confluence/todo"
)

func (node *Node) processGoFile(fpath string) error {
	contents, err := ioutil.ReadFile(filepath.Clean(fpath))
	if err != nil {
		return err
	}

	fullpath := strings.Replace(fpath, ".", "", 2)

	fullpath = strings.TrimPrefix(fullpath, "/")

	todo.ParseGo(contents, fullpath)

	return nil
}

func (node *Node) processMarkDown(path string) error {
	contents, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	parsedContents, err := markdown.ParseMarkdown(node.root.id, contents)
	if err != nil {
		return err
	}

	err = node.checkConfluencePages(parsedContents)
	if err != nil {
		log.Printf("error completing confluence operations: %s", err)
	}

	return nil
}

// uploadFile is for uploading files to a specific page by root node page id
func (node *Node) uploadFile(path string) {
	if nodeAPIClient != nil {
		err := nodeAPIClient.UploadAttachment(filepath.Clean(path), node.root.id)
		if err != nil {
			log.Printf("error uploading attachment: %s", err)
		}
	}
}
