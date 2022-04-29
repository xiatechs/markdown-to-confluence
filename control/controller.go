//Package control is the object that handles the iterating through files in a github repository
package control

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/xiatechs/markdown-to-confluence/apihandler"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
	"github.com/xiatechs/markdown-to-confluence/sem"
)

const (
	numberOfRoutines = 4 // limit number of goroutines (to balance load on confluence API)
)

var (
	rootDir string                               // will contain the root folderpath of the repo
	wg      = sem.NewSemaphore(numberOfRoutines) // waitgroup with limiter
)

type Controller struct {
	mu     *sync.RWMutex // for locking/unlocking when multiple goroutines are working on same node
	Root   *Node
	FH     filehandler.FileHandler
	API    apihandler.ApiController
	titles map[string]struct{}
}

func (c *Controller) Start(projectPath string) {
	c.mu = &sync.RWMutex{}

	c.titles = make(map[string]struct{})

	rootDir = strings.ReplaceAll(projectPath, `/github/workspace/`, "")

	rootDir = strings.ReplaceAll(rootDir, ".", "")

	rootDir = strings.ReplaceAll(rootDir, "/", "")

	Root := &Node{
		mu:       &sync.RWMutex{},
		filePath: rootDir,
		isFolder: true,
	}

	if Root.validate(c) {
		Root.hasMarkDown = true
		Root.checkFolderPageGeneration(c)
		Root.checkForFiles(c)
	}

	wg.Wait()

	Root.Delete(c)
}

func (node *Node) Delete(c *Controller) {
	if node.responseMetaData == nil {
		log.Println("NIL")
		return
	}

	if _, ok := node.responseMetaData["id"]; !ok {
		log.Println("NO ID")
		return
	}

	convert := strconv.Itoa(node.responseMetaData["id"].(int))

	c.API.Delete(c.titles, convert)
}
