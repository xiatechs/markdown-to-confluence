package node

// check - methods for checking various conditions

import (
	"log"
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

// checkIfRootAlive method checks if root node is alive,
// and creates a subnode (alive = markdown files exist in folder)
// if root is alive then root is the parent node for subnode
// else the root root node is the parent node for subnode
// then it calls generateMaster method on subnode
func (node *Node) checkIfRootAlive(fpath string) {
	if node.path != fpath {
		subNode := newNode()
		subNode.path = fpath

		if node.alive {
			subNode.root = node.root
		} else {
			if node.root != nil {
				subNode.root = node.root.root
			} else {
				subNode.root = node.root
			}
		}

		node.branches = append(node.branches, subNode)

		subNode.generateMaster()
	}
}

// fileInDirectoryCheck method takes file path and two bools (checking / folders)
// if folders is false & checking is true then it returns true if it finds a markdown file in a folder
// if folders is false & checking is false then it processes markdown files via checkIfMarkDown method
// if folders is true then it processes other file types via checkOtherFileTypes method
func (node *Node) fileInDirectoryCheck(fpath string, checking, folders bool) bool {
	if folders {
		node.checkOtherFileTypes(fpath) // you can process other file types inside this method

		return false
	}

	validFile := node.checkIfMarkDown(fpath, checking) // for checking & processing markdown files / images etc

	return validFile && checking
}

// checkIfMarkDown method checks is a folder or not, and if not
// passes file to checkIfMarkDownFile method
func (node *Node) checkIfMarkDown(fpath string, checking bool) bool {
	if !isFolder(fpath) {
		if ok := node.checkIfMarkDownFile(checking, fpath); ok {
			node.alive = true
			return true
		}
	}

	return false
}

// checkIfMarkDownFile method checks whether file is a markdown file or not
// checking bool is for whether we are just checking returning bool, or
// if we are doing work on file
func (node *Node) checkIfMarkDownFile(checking bool, name string) bool {
	if strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".MD") {
		if !checking {
			err := node.processMarkDown(name)
			if err != nil {
				log.Println(err)
			}

			foldersWithMarkdown++
		}

		return true
	}

	return false
}

// checkIfFolder method checks filepath is a folder or not
// and returns bool
func (node *Node) checkIfFolder(fpath string) bool {
	if isFolder(fpath) && !isVendorOrGit(fpath) {
		numberOfFolders++

		node.checkIfRootAlive(fpath)

		return true
	}

	return false
}

// checkOtherFileTypes method checks if file is a folder
// and if not, checks for if it is a go or image file
func (node *Node) checkOtherFileTypes(fpath string) {
	if !node.checkIfFolder(fpath) {
		node.checkIfGoFile(fpath)
		node.checkForImages(fpath)
	}
}

// checkIfGoFile method checks to see if the file is
// a golang file
func (node *Node) checkIfGoFile(name string) {
	if strings.HasSuffix(name, ".go") {
		err := node.processGoFile(name)
		if err != nil {
			log.Println(err)
		}
	}
}

// checkForImages method checks to see if the file is
// an image file
func (node *Node) checkForImages(name string) {
	if node.alive {
		validFiles := []string{".png", ".jpg", ".jpeg", ".gif"}

		for index := range validFiles {
			if strings.Contains(name, validFiles[index]) {
				node.checkNodeRootIsNil(name)
				return
			}
		}
	}
}

// checkNodeRootIsNil method checks whether the
// node root is nil before calling uploadFile method
func (node *Node) checkNodeRootIsNil(name string) {
	if node.root != nil {
		node.uploadFile(name)
	}
}

// checkConfluencePages method runs through the CRUD operations for confluence
func (node *Node) checkConfluencePages(newPageContents *markdown.FileContents) error {
	if nodeAPIClient == nil {
		return nil
	}

	pageTitle := strings.Join(strings.Split(newPageContents.MetaData["title"].(string), " "), "+")

	pageResult, err := nodeAPIClient.FindPage(pageTitle, false)
	if err != nil {
		return err
	}

	if pageResult == nil {
		err := node.generatePage(newPageContents)
		if err != nil {
			return err
		}
	} else {
		err = node.checkPageID(*pageResult)
		if err != nil {
			return err
		}
		err = nodeAPIClient.UpdatePage(node.id, int64(pageResult.Results[0].Version.Number), newPageContents)
		if err != nil {
			return err
		}
	}

	masterTitles = append(masterTitles, newPageContents.MetaData["title"].(string))

	return nil
}

// checkPageID method checks the pageID of the page contents and
// sets the node id to the page id
func (node *Node) checkPageID(pageResult confluence.PageResults) error {
	var err error

	if len(pageResult.Results) > 0 {
		node.id, err = strconv.Atoi(pageResult.Results[0].ID)
		if err != nil {
			return err
		}
	}

	return nil
}
