package node

// check - methods for checking various conditions

import (
	"fmt"
	"log"
	"path/filepath"
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
		subNode.treeLink = node.treeLink
		subNode.path = fpath
		subNode.images = node.images

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
//
// if folders is false & checking is true then it returns true if it finds a markdown file in a folder
//
// if folders is false & checking is false then it processes markdown files via checkIfMarkDown method
//
// if folders is true & checking is true then it returns true if it finds an image file in a folder
//
// if folders is true and checking is false then it processes uploads images and processes
// other file types via checkOtherFileTypes method
func (node *Node) fileInDirectoryCheck(fpath string, checking, folders bool) bool {
	if folders {
		return node.checkOtherFileTypes(fpath, checking) // you can process other file types inside this method
	}

	validFile := node.checkIfMarkDown(fpath, checking) // for checking & processing markdown files

	return validFile && checking
}

// checkIfMarkDown method checks is a folder or not, and if not
// passes file to checkIfMarkDownFile method
func (node *Node) checkIfMarkDown(fpath string, checking bool) bool {
	if !isFolder(fpath) {
		if ok := node.checkIfMarkDownFile(checking, fpath); ok {
			if checking {
				node.alive = true
			}

			return true
		}
	}

	return false
}

// checkIfMarkDownFile method checks whether file is a markdown file or not
// checking bool is for whether we are just checking returning bool, or
// if we are doing work on file
func (node *Node) checkIfMarkDownFile(checking bool, name string) bool {
	fileName := filepath.Base(name)

	if strings.HasSuffix(strings.ToLower(fileName), ".md") {
		if !checking {
			if strings.ToLower(fileName) == indexName { // we don't want to process index.md here
				return true
			}

			err := node.processMarkDown(name, fileName)
			if err != nil {
				log.Println(err)
			}

			return true
		}

		return true
	}

	return false
}

// checkIfFolder method checks filepath is a folder or not
// and returns bool
func (node *Node) checkIfFolder(fpath string) bool {
	if isFolder(fpath) && !isVendorOrGit(fpath) {
		node.checkIfRootAlive(fpath)

		return true
	}

	return false
}

// checkOtherFileTypes method checks if file is a folder
// and if not, checks for if it is a go or image file
func (node *Node) checkOtherFileTypes(fpath string, checking bool) bool {
	if !node.checkIfFolder(fpath) {
		node.checkIfGoFile(fpath)
		return node.checkForImages(fpath, checking)
	}

	return false
}

// checkIfGoFile method checks to see if the file is
// a golang file
func (node *Node) checkIfGoFile(name string) {
	if !strings.HasSuffix(name, ".go") {
		return
	}

	err := node.processGoFile(name)
	if err != nil {
		log.Println(err)
	}
}

// checkForImages method checks to see if the file is an image file
func (node *Node) checkForImages(name string, checking bool) bool {
	validFiles := []string{".png", ".jpg", ".jpeg", ".gif"}

	for index := range validFiles {
		if strings.HasSuffix(strings.ToLower(name), validFiles[index]) {
			if checking {
				node.alive = true
			} else {
				node.checkNodeRootIsNil(name)
			}

			return true
		}
	}

	return false
}

// checkNodeRootIsNil method checks whether the
// node root is nil before calling uploadFile method
func (node *Node) checkNodeRootIsNil(name string) {
	if node.root != nil {
		if node.imageToBeUploaded(name) {
			node.uploadFile(name, node.indexPage)
		}
	}
}

// imageToBeUploaded method checks if this imag ewas previously uploaded
func (node *Node) imageToBeUploaded(name string) bool {
	if _, exists := node.images[name]; exists {
		return false
	}

	node.images[name] = true

	return true
}

// checkConfluencePages method runs through the CRUD operations for confluence
func (node *Node) checkConfluencePages(newPageContents *markdown.FileContents, filepath string) error {
	_, abs := node.generateTitles()

	if newPageContents == nil {
		return fmt.Errorf("checkConfluencePages error for folder path [%s]: the markdown file was nil",
			abs)
	}

	pageTitle := strings.Join(strings.Split(newPageContents.MetaData["title"].(string), " "), "+")

	pageResult, err := nodeAPIClient.FindPage(pageTitle, false)
	if err != nil {
		return fmt.Errorf("find page error for folder path [%s] - page title [%s]: %w",
			abs, pageTitle, err)
	}

	if pageResult == nil {
		err := node.newPage(newPageContents)
		if err != nil {
			return fmt.Errorf("create page error for folder path [%s] - page title [%s]: %w",
				abs, pageTitle, err)
		}

		mapSem <- struct{}{}

		id := strconv.Itoa(node.id)

		node.treeLink.branches[filepath] = id

		log.Printf("processed file - id: [%d]", node.id)
		<-mapSem

		return nil
	}

	err = node.createOrUpdatePage(newPageContents, pageResult)
	if err != nil {
		return fmt.Errorf("create/update page error for folder path [%s] - page title [%s]: %w",
			abs, pageTitle, err)
	}

	mapSem <- struct{}{}

	id := strconv.Itoa(node.id)

	node.treeLink.branches[filepath] = id

	log.Printf("processed file - id: [%d]", node.id)
	<-mapSem

	return nil
}

func (node *Node) newPage(newPageContents *markdown.FileContents) error {
	_, abs := node.generateTitles()

	if newPageContents == nil {
		return fmt.Errorf("newPage error for folder path [%s]: the markdown file was nil",
			abs)
	}

	err := node.generatePage(newPageContents)
	if err != nil {
		return err
	}

	node.addContents(newPageContents)

	return nil
}

func (node *Node) createOrUpdatePage(newPageContents *markdown.FileContents,
	pageResult *confluence.PageResults) error {
	_, abs := node.generateTitles()

	if newPageContents == nil {
		return fmt.Errorf("createOrUpdatePage error for folder path [%s]: the newPageContents param was nil",
			abs)
	}

	if pageResult == nil {
		return fmt.Errorf("createOrUpdatePage pageResult error for folder path [%s]: the pageResult param was nil",
			abs)
	}

	err := node.checkPageID(*pageResult)
	if err != nil {
		return err
	}

	if len(pageResult.Results) > 0 {
		addToList, err := nodeAPIClient.UpdatePage(node.id, int64(pageResult.Results[0].Version.Number),
			newPageContents, *pageResult)
		if err != nil {
			return err
		}

		if addToList {
			node.addContents(newPageContents)
		}
	}

	return nil
}

// addContents adds the page title to either the parent page titles slice, or the node slice
// multiple goroutines could access same titles (or node.root.titles) slice so locking is required
func (node *Node) addContents(newPageContents *markdown.FileContents) {
	_, abs := node.generateTitles()

	if newPageContents.MetaData == nil {
		log.Printf("createOrUpdatePage warning for folder path [%s]: metadata was nil", abs)

		return
	}

	if node.root != nil {
		node.root.mu.Lock()

		defer node.root.mu.Unlock()

		node.root.titles = append(node.root.titles, newPageContents.MetaData["title"].(string))
	}

	node.mu.Lock()

	defer node.mu.Unlock()

	node.titles = append(node.titles, newPageContents.MetaData["title"].(string))
}

// checkPageID method checks the pageID of the page contents and
// sets the node id to the page id
// multiple goroutines could access same id field so locking is required
func (node *Node) checkPageID(pageResult confluence.PageResults) error {
	node.mu.RLock()

	defer node.mu.RUnlock()

	var err error

	if len(pageResult.Results) > 0 {
		node.id, err = strconv.Atoi(pageResult.Results[0].ID)
		if err != nil {
			return err
		}
	}

	return nil
}
