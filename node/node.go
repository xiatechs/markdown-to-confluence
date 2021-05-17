// Package node is to enable reading through a repo and create a tree of content on confluence
package node

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

var (
	numberOfFolders     float64
	foldersWithMarkdown float64
	rootDir             string                // will contain the root folderpath of the repo (without '/' and '.' in it)
	nodeAPIClient       *confluence.APIClient // api client will be stored here
)

// Node struct enables creation of a page tree
type Node struct {
	id       int                     // when page is created, page ID will be stored here.
	alive    bool                    // for tracking if the folder has any valid content within it asides more folders
	path     string                  // file / folderpath will be stored here
	root     *Node                   // the parent page node will be linked here
	branches []*Node                 // any children page nodes will be stored here (used to delete pages)
	children *confluence.PageResults // to store a snapshot of folder page & children pages (used to delete pages)
}

// Start method begins the generation of a tree of the repo for confluence
// first it validates whether the project path is a folder
// if yes then it sets the rootDir as the project path folder name
// then begins the recursive method generateMaster
// and returns bool - if true then it means pages have been created/updated/checked on confluence
// and there is markdown content in the folder
func (node *Node) Start(projectPath string, client *confluence.APIClient) bool {
	if isFolder(projectPath) {
		node.path = projectPath

		rootDir = strings.ReplaceAll(projectPath, `/github/workspace/`, "")

		rootDir = strings.ReplaceAll(rootDir, ".", "")

		rootDir = strings.ReplaceAll(rootDir, "/", "")

		nodeAPIClient = client

		node.generateMaster()

		var oneHundredPercent float64 = 100 // for calculating percentage of folders with markdown

		markDownPercentage := (foldersWithMarkdown / numberOfFolders) * oneHundredPercent

		percentageString := fmt.Sprintf("Folders with markdown percentage: %.2f%s", markDownPercentage, "%")

		node.generateTODOPage(percentageString)

		return true
	}

	return false
}

// iterate method is to scan through the files or folders in a folder.
// and takes in two bools (justChecking, foldersOnly)
// if justChecking is true then it will only check whether there is a valid file in folder
// and return true if there is
// if foldersOnly is true then it will only iterate through folders
func (node *Node) iterate(justChecking, foldersOnly bool) (validFile bool) {
	// Go 1.15 method: err := filepath.Walk(node.path, func(fpath string, info os.FileInfo, err error) error {
	// Go 1.16 method: err := filepath.WalkDir(node.path, func(fpath string, info os.DirEntry, err error) error {
	err := filepath.Walk(node.path, func(fpath string, info os.FileInfo, err error) error {
		if withinDirectory(node.path, fpath) {
			validFile = node.fileInDirectoryCheck(fpath, justChecking, foldersOnly)
			if validFile {
				return io.EOF
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	return validFile
}

// Delete method starts loop through node.branches
// and calls this method on each subnode of the node
// if node.id != 0 (i.e not the root node) then
// it calls method findPagesToDelete
func (node *Node) Delete() {
	if node.id != 0 {
		id := strconv.Itoa(node.id)
		node.findPagesToDelete(id)
	}

	for index := range node.branches {
		node.branches[index].Delete()
	}
}
