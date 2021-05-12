package node

// check - methods for checking various conditions

import (
	"log"
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

// checkIfRootAlive method - if root is alive, root is the parent node.
// else the root root is the parent node
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

func (node *Node) fileInDirectoryCheck(fpath string, checking, folders bool) bool {
	if !folders {
		validFile := node.checkIfMarkDown(fpath, checking) // for checking & processing markdown files / images etc
		if validFile && checking {
			return true
		}
	} else {
		node.checkOtherFileTypes(fpath) // you can process other file types inside this method
	}

	return false
}

func (node *Node) checkIfMarkDown(fpath string, checking bool) bool {
	if !isFolder(fpath) {
		if ok := node.checkIfMarkDownFile(checking, fpath); ok {
			node.alive = true
			foldersWithMarkdown++

			return true
		}
	}

	return false
}

func (node *Node) checkIfMarkDownFile(checking bool, name string) bool {
	if strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".MD") {
		if !checking {
			err := node.processMarkDown(name)
			if err != nil {
				log.Println(err)
			}
		}

		return true
	}

	return false
}

func (node *Node) checkIfFolder(fpath string) bool {
	if isFolder(fpath) && !isVendorOrGit(fpath) {
		node.checkIfRootAlive(fpath)
		return true
	}

	return false
}

func (node *Node) checkOtherFileTypes(fpath string) {
	if !node.checkIfFolder(fpath) {
		node.checkIfGoFile(fpath)
		node.checkForImages(fpath)
	}
}

func (node *Node) checkIfGoFile(name string) {
	if strings.HasSuffix(name, ".go") {
		err := node.processGoFile(name)
		if err != nil {
			log.Println(err)
		}
	}
}

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

func (node *Node) checkNodeRootIsNil(name string) {
	if node.root != nil {
		node.uploadFile(name)
	}
}

// checkConfluencePages runs through the CRUD operations for confluence
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

// checkPageID
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
