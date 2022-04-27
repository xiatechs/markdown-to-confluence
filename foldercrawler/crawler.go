//Package foldercrawler is the object that handles the iterating through files in a github repository
package foldercrawler

import (
	"strings"
	"sync"

	"github.com/xiatechs/markdown-to-confluence/apihandler"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
)

var rootDir string // will contain the root folderpath of the repo

type Controller struct {
	Root   *Node
	FH     filehandler.FileHandler
	API    apihandler.ApiController
	errors []error
}

func (c *Controller) ingestError(err error) {
	c.errors = append(c.errors, err)
}

func (c *Controller) Start(projectPath string) {
	rootDir = strings.ReplaceAll(projectPath, `/github/workspace/`, "")

	rootDir = strings.ReplaceAll(rootDir, ".", "")

	rootDir = strings.ReplaceAll(rootDir, "/", "")

	Root := &Node{
		mu:       &sync.RWMutex{},
		filePath: projectPath,
		isFolder: true,
	}

	if Root.validate(c) {
		Root.hasMarkDown = true
		Root.checkFolderPageGeneration(c)
		Root.checkForFiles(c)
	}
}
