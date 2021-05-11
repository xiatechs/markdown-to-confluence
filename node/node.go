// Package node is to enable reading through a repo and create a tree of content on confluence
package node

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

var (
	rootDir       string                // will contain the root folderpath of the repo (without '/' and '.' in it)
	nodeAPIClient *confluence.APIClient // api client will be stored here
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

// Start begins the generation of a tree of the repo for confluence
// and starts the whole process from the top/root node
func (node *Node) Start(projectPath string, client *confluence.APIClient) bool {
	if isFolder(projectPath) {
		node.path = projectPath
		
		rootDir = strings.ReplaceAll(projectPath, `/github/workspace/`, "")
		
		rootDir = strings.ReplaceAll(rootDir, ".", "")
		
		rootDir = strings.ReplaceAll(rootDir, "/", "")

		nodeAPIClient = client

		node.generateMaster()

		node.generateTODOPage()

		return true
	}

	return false
}

// iterate method is to scan through the files or folders in a folder.
func (node *Node) iterate(justChecking, foldersOnly bool) (validFile bool) {
	// Go 1.15 method: err := filepath.Walk(node.path, func(fpath string, info os.FileInfo, err error) error {
	// Go 1.16 method: err := filepath.WalkDir(node.path, func(fpath string, info os.DirEntry, err error) error {
	err := filepath.WalkDir(node.path, func(fpath string, info os.DirEntry, err error) error {
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

// Delete method clears away any pages on confluence that shouldn't exist
// this method should be called from the top node
func (node *Node) Delete() {
	if node.id != 0 {
		id := strconv.Itoa(node.id)
		node.findPagesToDelete(id)
	}

	for index := range node.branches {
		node.branches[index].Delete()
	}
}
