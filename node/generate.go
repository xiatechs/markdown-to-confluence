package node

// generate - methods where pages/content/nodes are being created

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	goplantuml "github.com/jfeliu007/goplantuml/parser"
	"github.com/xiatechs/markdown-to-confluence/markdown"
	"github.com/xiatechs/markdown-to-confluence/todo"
)

// generateMaster method checks whether the folder is alive (has markdown files in it)
// and if it is, creates a page for the folder. Also, then checks whether the folder
// has subfolders in it, and then begins the process of checking those folders (recursively)
func (node *Node) generateMaster() {
	// these constants are to aid navigation of iterate method lower down
	const checking = true

	const processing = false

	const folders = true

	const files = false

	subNode := newNode()
	subNode.path = node.path
	subNode.root = node
	node.branches = append(node.branches, subNode)

	thereAreValidFiles := subNode.iterate(checking, files)
	if thereAreValidFiles {
		node.alive = true
		node.generateFolderPage()
		subNode.generatePlantuml(node.path) // generate plantuml in folders with markdown in it only
		subNode.iterate(processing, files)
	}

	subNode.iterate(processing, folders)
}

// generateFolderPage method creates a folder page in confluence for a folder
func (node *Node) generateFolderPage() {
	dir, fullDir := node.generateTitles()

	masterpagecontents := markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": dir,
		},
		Body: []byte(`<p>Welcome to the '<b>` + dir + `</b>' folder of this Xiatech code repo.</p>
		<p>This folder full path in the repo is: ` + fullDir + `</p>
<p>You will find attachments/images for this folder via the ellipsis at the top right.</p>
<p>Any markdown or subfolders is available in children pages under this page.</p>`),
	}

	err := node.checkConfluencePages(&masterpagecontents)
	if err != nil {
		log.Println(err)
	}
}

// generateTODOPage method creates a page in parent folder
// that contains todo's for a codebase
func (node *Node) generateTODOPage(percentage string) {
	todonode := Node{}
	todonode.root = node

	page := todo.GenerateTODO(rootDir, percentage)

	if page != nil {
		err := todonode.checkConfluencePages(page)
		if err != nil {
			log.Println(err)
		}
	}
}

// generateTitles returns two strings
// string 1 - folder of the node
// string 2 - the absolute filepath to the node dir from root dir
func (node *Node) generateTitles() (string, string) {
	const nestedDepth = 2

	fullDir := strings.ReplaceAll(node.path, "/github/workspace/", "")
	fullDir = strings.ReplaceAll(fullDir, ".", "")
	fullDir = strings.TrimPrefix(fullDir, "/")
	dirList := strings.Split(fullDir, "/")
	dir := dirList[len(dirList)-1]

	if len(dirList) > nestedDepth {
		dir += "-"
		dir += dirList[len(dirList)-2]
	}

	if node.root != nil {
		dir += "-"
		dir += rootDir
	}

	return dir, fullDir
}

// generatePlantuml takes in a folder path and
// generates a .puml file of the go code in the folder
// then calls generatePlantumlImage method to create a picture
// then creates a page for the image to be uploaded to and displayed
func (node *Node) generatePlantuml(fpath string) {
	const minimumPageSize = 20 // plantuml that is generated <= 20 chars long is too small

	const iterateThroughSubFolders = false

	path, _ := node.generateTitles()

	if node.root.root == nil {
		path = rootDir
	}

	result, err := goplantuml.NewClassDiagram([]string{fpath}, []string{}, iterateThroughSubFolders)
	if err != nil {
		log.Println("plantuml file generation error: %w", err)
		return
	}

	rendered := result.Render()
	if len(rendered) > minimumPageSize {
		var filename = path + "-pumldiagram"

		var buf bytes.Buffer

		var headerstring = `<p><img src="` + filename + ".png" + `"/></img></p>`

		headerstring = markdown.URLConverter(node.root.id, headerstring)

		var writer io.Writer

		writer, err = os.Create(node.path + "/" + filename + ".puml")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}

		fmt.Fprint(&buf, headerstring)

		fmt.Fprint(writer, rendered)

		node.generatePlantumlImage(node.path + "/" + filename + ".puml")

		masterpagecontents := markdown.FileContents{
			MetaData: map[string]interface{}{
				"title": "plantuml-" + path,
			},
			Body: buf.Bytes(),
		}

		err = node.checkConfluencePages(&masterpagecontents)
		if err != nil {
			log.Println(err)
		}
	}
}

// generatePlantumlImage method calls external application (plantuml.jar)
// in the docker container to generate the plantuml image (as a .png)
func (node *Node) generatePlantumlImage(fpath string) {
	convertPlantuml := exec.Command("java", "-jar", "/app/plantuml.jar", "-tpng", fpath) // #nosec - pumlimage
	convertPlantuml.Stdout = os.Stdout

	err := convertPlantuml.Run()
	if err != nil {
		log.Println(err)
	}
}

// generatePage method creates a new page for node.
// and sets the parent page as the node root id.
// unless the node.root is nil in which case it is the root page
func (node *Node) generatePage(newPageContents *markdown.FileContents) error {
	var isParentPage = true

	var err error

	if nodeAPIClient != nil {
		if node.root == nil {
			node.id, err = nodeAPIClient.CreatePage(0, newPageContents, isParentPage)
		} else {
			node.id, err = nodeAPIClient.CreatePage(node.root.id, newPageContents, !isParentPage)
		}
	}

	return err
}
