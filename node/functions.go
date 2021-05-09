package node

// helper & factory functions

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

// newNode - create a new node object
func newNode() *Node {
	node := Node{}
	return &node
}

// newPageResults - create a new confluence.PageResults object
func newPageResults() *confluence.PageResults {
	results := confluence.PageResults{}
	return &results
}

// checks to see if the file is within the folder
func withinDirectory(base, path string) bool {
	return strings.Count(path, "/")-strings.Count(base, "/") == 1
}

// checking if file is vendor / git folder
func isVendorOrGit(name string) bool {
	if strings.Contains(name, "vendor") || strings.Contains(name, ".github") || strings.Contains(name, ".git") {
		return true
	}

	return false
}

// isFolder checks whether a file is a folder
func isFolder(name string) bool {
	file, err := os.Open(filepath.Clean(name))
	if err != nil {
		log.Println(err)
		return false
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return false
	}

	if fileInfo.IsDir() {
		return true
	}

	return false
}
