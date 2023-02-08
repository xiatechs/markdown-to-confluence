// Package node is to enable reading through a repo and create a tree of content on confluence
package node

//notodo: no need
import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/semaphore"
)

const (
	numberOfRoutines = 10 // limit number of goroutines (to balance load on confluence API)
	indexName        = "readme.md"
)

var (
	mapSem              = make(chan struct{}, 1)                   // for controlling access to Tree map
	wg                  = semaphore.NewSemaphore(numberOfRoutines) // for controlling number of goroutines
	numberOfFolders     float64                                    // for counting number of folders in repo
	foldersWithMarkdown float64                                    // for counting number of folders with markdown in repo
	rootDir             string                                     // will contain the root folderpath of the repo
	// NodeAPIClient is interface where a confluence API client can be placed
	nodeAPIClient APIClienter // api client will be stored here
	t             *Tree
)

// SetAPIClient sets a confluence.APIClient object into the node package
func SetAPIClient(client APIClienter) {
	nodeAPIClient = client
}

// Tree - capture what has been generated
type Tree struct {
	branches map[string]string
}

// Node struct enables creation of a page tree
type Node struct {
	treeLink  *Tree
	masterID  int    // the confluence page ID for the parent page the mtc tool should create the files in
	id        int    // when page is created, page ID will be stored here.
	alive     bool   // for tracking if the folder has any valid content within it asides more folders
	path      string // file / folderpath will be stored here
	hasIndex  bool
	root      *Node           // the parent page node will be linked here
	branches  []*Node         // any children page nodes will be stored here (for deleting)
	titles    []string        // titles of pages created by node (for deleting)
	images    map[string]bool // image files uploaded to prevent uploading them multiple times
	mu        *sync.RWMutex   // for locking/unlocking when multiple goroutines are working on same node
	indexPage bool
	indexName string
}

// Start method begins the generation of a tree of the repo for confluence
// first it validates whether the project path is a folder
// if yes then it sets the rootDir as the project path folder name
// then begins the recursive method generateMaster
// and returns bool - if true then it means pages have been created/updated/checked on confluence
// and there is markdown content in the folder
func (node *Node) Start(projectMasterID int, projectPath string, onlyDocs bool) bool {
	if t == nil {
		log.Println("instantiating TREE")

		t = func() *Tree { // t - Tree - will contain tree of pages created and their subsequent confluence URL
			return &Tree{
				branches: make(map[string]string),
			}
		}()
	}

	node.treeLink = t

	node.mu = &sync.RWMutex{}

	node.images = map[string]bool{}

	/*
		FOR RELATIVE FILE LINKS IN CONFLUENCE...

		Start needs to be called twice because confluence page ID's are captured _only_ when they are generated
		successfully. So, if there are local links inside pages, they cannot be determined until after the page
		is created - bizarre i know - so the logic has to be run twice so that we can first:

		- generate a tree of pages with their confluence ID's
		- re-generate a tree of pages with the confluence ID's known
	*/

	if isFolder(projectPath) {
		numberOfFolders++

		node.masterID = projectMasterID

		node.path = projectPath

		rootDir = strings.ReplaceAll(projectPath, `/github/workspace/`, "")

		rootDir = strings.ReplaceAll(rootDir, ".", "")

		rootDir = strings.ReplaceAll(rootDir, "/", "")

		err := node.generateFolderPage(false) // create the main page first
		if err != nil {
			return false
		}

		node.generateMaster() // contains concurrency

		log.Println("WAITING FOR GOROUTINES [1] - generating first run to capture page ID's")

		wg.Wait()

		node.generateMaster() // contains concurrency

		log.Println("WAITING FOR GOROUTINES [2] - now trying to match relative links to page ID's")

		wg.Wait()

		log.Println("Here are the pages I got today:")

		node.Tree()

		log.Println("FINISHED GOROUTINES - NOW CHECKING FOR DELETE")

		return true
	}

	return false
}

// Tree - print out what has been generated
func (node *Node) Tree() {
	for path, id := range t.branches {
		log.Println(path, "|", common.ConfluenceBaseURL+"/wiki/spaces/"+common.ConfluenceSpace+"/pages/"+id)
	}
}

// iterate method is to scan through the files or folders in a folder.
// and takes in two bools (justChecking, foldersOnly)
// if justChecking is true then it will only check whether
// there is a valid file in folder and return true if there is
// if foldersOnly is true then it will only iterate through folders
func (node *Node) iterate(justChecking, foldersOnly bool) bool {
	var thereIsAValidFile bool

	err := filepath.Walk(node.path, func(fpath string, info os.FileInfo, err error) error {
		if node.withinDirectory(node.path, fpath) {
			if strings.ToLower(filepath.Base(fpath)) == indexName {
				node.hasIndex = true
				node.alive = true
				if node.root != nil {
					node.root.indexName = filepath.Base(fpath)
				} else {
					node.indexName = filepath.Base(fpath)
				}
			}

			validFile := node.fileInDirectoryCheck(fpath, justChecking, foldersOnly)
			if validFile {
				thereIsAValidFile = true
			}
		}

		return nil
	})
	if err != nil {
		log.Println("iterate: ", err)
	}

	return thereIsAValidFile
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
