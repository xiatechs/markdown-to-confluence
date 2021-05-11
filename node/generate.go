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

// method for checking if folder is 'alive' i.e contains markdown files.
// if it is, create a folder page for that folder
// then create a subnode, and that node will iterate
func (node *Node) generateMaster() {
	// these constants are to aid navigation of iterate method lower down
	const checking = true

	const processing = false

	const folders = true

	const files = false

	subNode := newNode()
	subNode.path = node.path
	subNode.root = node
	subNode.children = newPageResults()
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

// generateFolderPage method
// if called, this node is a master node for a folder which has content in it
// a page for that folder will be created in confluence
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

func (node *Node) generateTODOPage() {
	todonode := Node{}
	todonode.root = node

	page := todo.GenerateTODO(rootDir)

	err := todonode.checkConfluencePages(page)
	if err != nil {
		log.Println(err)
	}
}

func (node *Node) generateTitles() (string, string) {
	const nestedDepth = 2

	fullDir := strings.ReplaceAll(node.path, ".", "")
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

func (node *Node) generatePlantumlImage(fpath string) {
	convertPlantuml := exec.Command("java", "-jar", "/app/plantuml.jar", "-tpng", fpath)
	convertPlantuml.Stdout = os.Stdout

	err := convertPlantuml.Run()
	if err != nil {
		log.Println(err)
	}
}

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
