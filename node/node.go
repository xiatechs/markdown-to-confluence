// Package node is to read through a repo and create a tree of content
package node

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

// Node struct enables creation of a page tree
// Each node is either a folder or a page.
type Node struct {
	index      int
	path       string
	rootFolder *Node
	rootID     int
	alive      bool
}

type details struct {
	index int
	info  string
}

var treeOverView = []details{}

// PrintOverview prints out the overview of the roots of tree
// so you can see what has content
func PrintOverview() {
	if len(treeOverView) == 0 {
		fmt.Println("Overview is empty")
		return
	}

	sort.SliceStable(treeOverView, func(i, j int) bool {
		return treeOverView[i].index < treeOverView[j].index
	})

	for index := range treeOverView {
		log.Println(treeOverView[index].info)
	}
}

// Instantiate begins the generation of a tree of the repo for confluence
func (node *Node) Instantiate(projectPath string) bool {
	if isFolder(projectPath) {
		node.index = 1
		node.path = projectPath
		node.generateMaster()

		return true
	}

	return false
}

func (node *Node) addToOverView() {
	var d = details{}

	if node.rootFolder != nil {
		toadd := fmt.Sprintf("Path: %s, ID: %d, Has content: %t", node.path, node.rootID, node.alive)
		d.info = toadd
		d.index = node.index
		treeOverView = append(treeOverView, d)
	} else {
		toadd := fmt.Sprintf("Path: %s, ID: %d, Has content: %t", node.path, node.rootID, node.alive)
		d.info = toadd
		d.index = 0
		treeOverView = append(treeOverView, d)
	}
}

func newFolder() *Node {
	folder := Node{}
	return &folder
}

// check to see if the name of the file ends with .md i.e it's a markdown file
func (node *Node) checkreadme(name string) {
	if strings.HasSuffix(name, ".md") {
		node.alive = true

		err := node.processFile(name)
		if err != nil {
			log.Println(err)
		}
	}
}

// check to see if the file is a puml or png image
func (node *Node) checkpuml(fpath, name string) {
	if node.alive && strings.Contains(name, ".puml") || strings.Contains(name, ".jpg") {
		if err := uploadFile(name, node.rootFolder.rootID); err != nil {
			log.Println(err)
		}
	}
}

// processFile is the function called on eligible files to handle uploads.
func (node *Node) processFile(path string) error {
	log.Println("Processing:", filepath.Clean(path))

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Printf("error opening file (%s): %v", path, err)
		return err
	}

	contents, err := markdown.ParseMarkdown(f)
	if err != nil {
		return err
	}

	err = node.checkConfluencePages(contents)
	if err != nil {
		log.Printf("error completing confluence operations: %s", err)
	}

	return nil
}

// uploadFile is for uploading files to a specific page by page ID
func uploadFile(path string, pageID int) error {
	log.Println("Processing:", filepath.Clean(path))

	Client, err := confluence.CreateAPIClient()
	if err != nil {
		log.Printf("error creating APIClient: %s", err)
		return err
	}

	err = Client.UploadAttachment(filepath.Clean(path), pageID)
	if err != nil {
		log.Printf("error uploading attachment: %s", err)
		return err
	}

	return nil
}

func (node *Node) generatePage(newPageContents *markdown.FileContents, client *confluence.APIClient) error {
	var err error

	if client != nil {
		if node.rootFolder == nil {
			node.rootID, err = client.CreatePage(0, newPageContents)
		} else {
			node.rootID, err = client.CreatePage(node.rootFolder.rootID, newPageContents)
		}
	}

	return err
}

// checkConfluencePages runs through the CRUD operations for confluence
func (node *Node) checkConfluencePages(newPageContents *markdown.FileContents) error {
	Client, err := confluence.CreateAPIClient()
	if err != nil {
		log.Printf("error creating APIClient: %s", err)
		return err
	}

	pageTitle := strings.Join(strings.Split(newPageContents.MetaData["title"].(string), " "), "+")

	pageResult, err := Client.FindPage(pageTitle)
	if err != nil {
		return err
	}

	if pageResult == nil {
		err := node.generatePage(newPageContents, Client)
		if err != nil {
			return err
		}
	} else {
		var err error
		node.rootID, err = strconv.Atoi(pageResult.Results[0].ID)
		if err != nil {
			return err
		}

		err = Client.UpdatePage(node.rootID, int64(pageResult.Results[0].Version.Number), newPageContents, node.rootID)
		if err != nil {
			return err
		}
	}

	return nil
}

func isVendorOrGit(name string) bool {
	if strings.Contains(name, "vendor") || strings.Contains(name, ".github") {
		return true
	}

	return false
}

// isFolder checks whether a file is a folder or not
func isFolder(name string) bool {
	file, err := os.Open(filepath.Clean(name))
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if fileInfo.IsDir() {
		return true
	}

	return false
}

// checkAll is where we will create or update the page, and upload or update attachments
func (node *Node) checkAll(path string) {
	node.checkreadme(node.path)
	node.checkpuml(node.path, path)
}

// checks to see if the file is within a subdirectory of the base path
func sub(base, path string) bool {
	return strings.Count(path, "/")-strings.Count(base, "/") == 1
}

// generateMaster function is to generate a master Node struct where we can append files
// to the folder node as well as subfolders.
func (node *Node) generateMaster() {
	root := strings.ReplaceAll(node.path, ".", "")
	masterpagecontents := markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": root,
		},
		Body: []byte("This is a placeholder page for the " + root + " folder in this repo"),
	}

	err := node.checkConfluencePages(&masterpagecontents)
	if err != nil {
		fmt.Println(err)
	}

	subNode := newFolder()
	subNode.index = node.index + 1
	subNode.path = node.path
	subNode.rootFolder = node
	subNode.iteratefiles()
}

// iteratefiles function is to iterate through the files in a folder.
// if it finds a file it will begin processing that file
func (node *Node) iteratefiles() {
	err := filepath.Walk(node.path, func(fpath string, info os.FileInfo, err error) error {
		if isVendorOrGit(fpath) {
			return filepath.SkipDir
		}
		if !isFolder(fpath) {
			if sub(node.path, fpath) {
				node.checkAll(fpath)
			}
		}
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	node.iteratefolders()
}

func (node *Node) verifyCreateNode(fpath string) {
	if node.path != fpath && sub(node.path, fpath) {
		subNode := newFolder()
		subNode.path = fpath

		if node.alive {
			subNode.rootFolder = node.rootFolder
		} else {
			subNode.rootFolder = node.rootFolder.rootFolder
		}

		subNode.index = node.index + 1
		subNode.generateMaster()
	}
}

// iteratefolders function is to iterate through the subfolders of a folder
// if it finds a folder, it will create a new Node
// and begin repeating the same process from that node
func (node *Node) iteratefolders() {
	err := filepath.Walk(node.path, func(fpath string, info os.FileInfo, err error) error {
		if isVendorOrGit(fpath) {
			return filepath.SkipDir
		}
		if isFolder(fpath) {
			node.verifyCreateNode(fpath)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	node.addToOverView()
}
