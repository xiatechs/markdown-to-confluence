// Package control is the object that handles the iterating through files in a github repository
//nolint: staticcheck // err shadowing isn't always bad
package control

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

// Node - a folder in the repository
type Node struct {
	// for locking/unlocking when multiple goroutines are working on same node
	mu                    *sync.RWMutex    

	// you can return anything from API response here         
	responseMetaData      map[string]interface{}   
	
	// the parent node meta data will be stored here & passed through to the API
	parentMetaData        map[string]interface{}   
	
	// the name of the page (for confluence - page names have to be unique)
	indexName             string      
	
	// the full path to the file
	filePath              string
	
	// is the file a folder?
	isFolder              bool   

	// does the folder have markdown                 
	hasMarkDown           bool   

	// this will be the last folder above this folder that had markdown in it                   
	lastAliveParentFolder *Node   
	
	// any (live) files underneath will be mapped here by filePath
	subFiles              []*Node
	
	// if the folder had a README.md in it - this will be the file contents
	readMeFile            *filehandler.FileContents 

	// if it's any other file - this will be the file contents i.e images... whatever really
	fileContents          *filehandler.FileContents 
}

func (node *Node) validate(c *Controller) (alive bool) {
	pageTitle, _ := node.generatePaths()

	err := filepath.Walk(node.filePath, func(fpath string, info os.FileInfo, err error) error {
		if !withinDirectory(node.filePath, fpath) {
			return nil
		}

		if checkIfFolder(fpath) {
			return nil
		}

		if isVendorOrGit(fpath) {
			return nil
		}

		valid := isMarkdownFile(fpath)
		if valid {
			if isReadMeFile(fpath) {
				fileContents, fiErr := c.FH.ConvertMarkdown(fpath, pageTitle, node.parentMetaData)
				if fiErr != nil {
					return err
				}

				node.readMeFile = fileContents
			}

			alive = true
		}

		return nil
	})
	if err != nil {
		if err != io.EOF {
			node.usefulLogError("validate", err)
		}

		return false
	}

	return alive
}

//nolint: gocognit // keep all logic in one place
func (node *Node) checkFolderPageGeneration(c *Controller) error {
	pageTitle, _ := node.generatePaths()

	if node.hasMarkDown { // if the page is alive...
		if isFolder(node.filePath) { // if this is a folder...
			if node.readMeFile != nil { // if this folder has a README - let's create an index page from it
				err := node.generateReadMeIndexPage(c)
				if err != nil {
					return err
				}
			} else {
				err := node.generateGenericIndexPage(c)
				if err != nil {
					return err
				}
			}
		}
	}

	node.indexName = filepath.Base(pageTitle)

	err := filepath.Walk(node.filePath, func(fpath string, info os.FileInfo, err error) error {
		if !withinDirectory(node.filePath, fpath) {
			return nil
		}

		if isVendorOrGit(fpath) {
			return nil
		}

		if isMarkdownFile(fpath) {
			return nil
		}

		if isFolder(fpath) { // if this subpath is a folder, we'll rinse and repeat
			wg.Add() // concurrently
			go func() {
				defer wg.Done()
				childNode := node.createChildNode(fpath, c)

				node.scanUpForParent(childNode)

				err := childNode.checkFolderPageGeneration(c)
				if err != nil {
					node.usefulLogError("childNode.checkFolderPageGeneration", err)
				}

				err = childNode.checkForFiles(c)
				if err != nil {
					node.usefulLogError("childNode.checkForFiles", err)
				}
			}()
			return nil
		}

		err = node.processOtherFiles(c, fpath)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if err != io.EOF {
			node.usefulLogError("checkFolderPageGeneration", err)
		}
	}

	node.end()

	return nil
}

func (node *Node) processOtherFiles(c *Controller, fpath string) error {
	otherFileNode := node.createChildNode(fpath, c)
	c.mu.RLock()
	defer c.mu.RUnlock()

	node.scanUpForParent(otherFileNode)

	otherFileNodeTitle, _ := otherFileNode.generatePaths()

	fileContents, err := c.FH.ProcessOtherFile(fpath, otherFileNodeTitle, node.responseMetaData)
	if err != nil {
		return err
	}

	if _, ok := fileContents.MetaData["type"]; !ok {
		// this file is not being handled by the current filehandler so ignore it
		return nil
	}

	otherFileNode.fileContents = fileContents

	otherFileNode.responseMetaData, err = c.API.CRUD(fileContents, node.responseMetaData)
	if err != nil {
		return err
	}

	if _, ok := otherFileNode.responseMetaData["title"]; !ok {
		return nil
	}

	c.titles[otherFileNode.responseMetaData["title"].(string)] = struct{}{}
	
	return nil
}

func (node *Node) generateReadMeIndexPage(c *Controller) error {
	node.readMeFile.MetaData["indexPage"] = true

	node.readMeFile.MetaData["alive"] = true

	if node.lastAliveParentFolder == nil {
		node.readMeFile.MetaData["root"] = true

		pageTitle, _ := node.generatePaths()

		node.readMeFile.MetaData["title"] = pageTitle
	}

	var err error

	node.responseMetaData, err = c.API.CRUD(node.readMeFile, node.parentMetaData)
	if err != nil {
		return err
	}

	if node.responseMetaData == nil {
		return nil
	}

	if _, ok := node.responseMetaData["title"]; !ok {
		return nil
	}

	c.mu.RLock()

	c.titles[node.responseMetaData["title"].(string)] = struct{}{}

	c.mu.RUnlock()

	return nil
}

func (node *Node) generateGenericIndexPage(c *Controller) error {
	pageTitle, _ := node.generatePaths()

	var err error

	// else, let's create a 'generic folder page' for indexing
	folderContents, err := c.FH.ConvertFolder(node.filePath, pageTitle, node.parentMetaData) 
	if err != nil {
		return err
	}

	folderContents.MetaData["alive"] = true

	node.fileContents = folderContents

	if node.lastAliveParentFolder == nil {
		node.fileContents.MetaData["root"] = true
	}

	node.responseMetaData, err = c.API.CRUD(folderContents, node.parentMetaData)
	if err != nil {
		return err
	}

	if node.responseMetaData == nil {
		return nil
	}

	if _, ok := node.responseMetaData["title"]; !ok {
		return nil
	}

	c.mu.RLock()

	c.titles[node.responseMetaData["title"].(string)] = struct{}{}

	c.mu.RUnlock()
	
	return nil
}

func (node *Node) processMarkdown(fpath string, c *Controller) error {
	var err error

	otherFileNode := node.createChildNode(fpath, c)

	node.scanUpForParent(otherFileNode)

	otherFileNodeTitle, _ := otherFileNode.generatePaths()

	fileContents, err := c.FH.ConvertMarkdown(fpath, otherFileNodeTitle, node.responseMetaData)
	if err != nil {
		return err
	}

	otherFileNode.fileContents = fileContents

	// create a webpage using the parent node meta data (i.e link this page to the parent page)
	otherFileNode.responseMetaData, err = c.API.CRUD(fileContents, node.responseMetaData)
	if err != nil {
		return err
	}

	if otherFileNode.responseMetaData == nil {
		return nil
	}

	if _, ok := otherFileNode.responseMetaData["title"]; !ok {
		return nil
	}

	c.mu.RLock()
	c.titles[otherFileNode.responseMetaData["title"].(string)] = struct{}{}
	c.mu.RUnlock()

	return nil
}

func (node *Node) checkForFiles(c *Controller) error {
	err := filepath.Walk(node.filePath, func(fpath string, info os.FileInfo, err error) error {
		if !withinDirectory(node.filePath, fpath) {
			return nil
		}

		if checkIfFolder(fpath) {
			return nil
		}

		if isVendorOrGit(fpath) {
			return nil
		}

		if !(isMarkdownFile(fpath) && !isReadMeFile(fpath)) {
			return nil
		}

		err = node.processMarkdown(fpath, c)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if err != io.EOF {
			node.usefulLogError("node.start", err)
		}
	}

	return nil
}

func (node *Node) createChildNode(fpath string, c *Controller) *Node {
	childNode := &Node{
		mu:       &sync.RWMutex{},
		filePath: fpath,
	}
	childNode.mu.Lock()
	defer childNode.mu.Unlock()

	alive := childNode.validate(c)
	if alive {
		childNode.hasMarkDown = true
	}

	node.scanUpForParent(childNode)

	return childNode
}

func (node *Node) scanUpForParent(theChildNode *Node) {
	if node.hasMarkDown {
		theChildNode.lastAliveParentFolder = node

		theChildNode.parentMetaData = node.responseMetaData

		node.mu.RLock()

		node.subFiles = append(node.subFiles, theChildNode)

		node.mu.RUnlock()

		return
	}

	if node.lastAliveParentFolder != nil {
		node.lastAliveParentFolder.scanUpForParent(theChildNode)
		return
	}

	theChildNode.lastAliveParentFolder = node

	theChildNode.parentMetaData = node.responseMetaData

	node.mu.RLock()

	node.subFiles = append(node.subFiles, theChildNode)

	node.mu.RUnlock()
}

func (node *Node) end() {
	node.mu.Lock()
	defer node.mu.Unlock()
	
	if !node.hasMarkDown {
		log.Printf("NODE: isFolder: [%t] - [%s] and has no markdown i.e dead",
			node.isFolder, node.filePath)
		return
	}

	if node.lastAliveParentFolder == nil {
		log.Printf("NODE: isFolder: [%t] - [%s] has [%d] sub files and is the ROOT folder",
			node.isFolder, node.filePath, len(node.subFiles))

		return
	}

	log.Printf("NODE: isFolder: [%t] - [%s] has [%d] sub files and the last living parent folder is [%s]",
		node.isFolder, node.filePath, len(node.subFiles), node.lastAliveParentFolder.filePath)
}

// generateTitles returns two strings
// string 1 - folder of the node
// string 2 - the absolute filepath to the node dir from root dir
func (node *Node) generatePaths() (string, string) {
	const nestedDepth = 2

	fullDir := strings.ReplaceAll(node.filePath, common.GitHubPrefix, "")

	fullDir = strings.ReplaceAll(fullDir, ".", "")

	fullDir = strings.TrimPrefix(fullDir, "/")

	dirList := strings.Split(fullDir, "/")

	dir := dirList[len(dirList)-1]

	if len(dirList) > nestedDepth {
		dir += "-"
		dir += dirList[len(dirList)-2]
	}

	if node.lastAliveParentFolder != nil {
		dir += "-"
		dir += rootDir
	}

	return dir, fullDir
}

func checkIfFolder(fpath string) bool {
	return (isFolder(fpath) && !isVendorOrGit(fpath))
}

func (node *Node) usefulLogError(functionName string, err error) {
	path, fullPath := node.generatePaths()
	log.Printf("within func [%s] - error at [%s] - path [%s]: %v",
		functionName, path, fullPath, err)
}
