package node

// generate - methods where pages/content/nodes are being created

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	goplantuml "github.com/jfeliu007/goplantuml/parser"
	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/markdown"
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
	subNode.treeLink = node.treeLink
	node.branches = append(node.branches, subNode)

	thereAreValidFiles := subNode.iterate(checking, files)

	if !thereAreValidFiles {
		log.Printf("skipping this folder [%s] because no valid files here", node.path)
		subNode.iterate(processing, folders)

		return
	}

	log.Printf("generating a folder page here [%s] this page is alive", node.path)

	err := node.generateFolderPage(subNode.hasIndex, subNode.id)
	if err != nil {
		log.Println(fmt.Errorf("generate folder page error: %w", err))
		return
	}

	node.alive = true

	subNode.iterate(processing, folders)

	subNode.generateChildPages()
}

// generateChildPages method generates all children pages for all parent pages
// can be run concurrently as they all have a parent page to attach to
// so there's no need to order their generation
func (node *Node) generateChildPages() {
	const processing = false

	const files = false

	wg.Add()

	go func() {
		defer wg.Done()
		node.generatePlantuml(node.path) // generate plantuml in folders with markdown in it only
		node.iterate(processing, files)  // generate child pages for any valid files in parent page
	}()
}

// generateFolderPage method creates a folder page in confluence for a folder
func (node *Node) generateFolderPage(hasIndex bool, subindex int) error {
	dir, fullDir := node.generateTitles()
	log.Printf("START processing file [%s]", fullDir)

	if hasIndex {
		log.Printf("this location [%s] has a [%s] file so will use that as index at this location",
			node.path, node.indexName)

		node.indexPage = true

		masterpagecontents, err := node.processMarkDownIndex(filepath.Join(node.path, node.indexName), subindex)
		if err != nil {
			return err
		}

		err = node.checkConfluencePages(masterpagecontents, node.path)
		if err != nil {
			log.Printf("[generate folderpage] generation error for path [%s]: %v", node.path, err)
			return err
		}
		// have to do it twice...

		masterpagecontents, err = node.processMarkDownIndex(filepath.Join(node.path,
			node.indexName), node.id)
		if err != nil {
			return err
		}

		err = node.checkConfluencePages(masterpagecontents, filepath.Join(node.path,
			node.indexName))
		if err != nil {
			log.Printf("[generate folderpage] generation error for path [%s]: %v", node.path, err)
			return err
		}

		log.Printf("processed bespoke index file - id: [%d]", node.id)

		return nil
	}

	node.indexPage = false
	log.Printf("no [%s] located here [%s], will generate a generic folderpage",
		indexName, node.path)

	masterpagecontents := &markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": dir,
		},
		Body: []byte(`<p>Welcome to the '<b>` + dir + `</b>' folder of this Xiatech code repo.</p>
		<p>This folder full path in the repo is: ` + fullDir + `</p>
<p>You will find attachments/images for this folder via the ellipsis at the top right.</p>
<p>Any markdown or subfolders is available in children pages under this page.</p>`),
	}

	err := node.checkConfluencePages(masterpagecontents, node.path)
	if err != nil {
		log.Printf("[generate folderpage] generation error for path [%s]: %v", node.path, err)
		return err
	}

	log.Printf("processed generic index file - id: [%d]", node.id)

	return nil
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

	const iterateThroughSubFolders = false // we want to just generate plantuml for the folder

	path, abs := node.generateTitles()

	if node.root.root == nil {
		path = rootDir
	}

	log.Printf("generating plantuml text for %s", path)

	result, err := goplantuml.NewClassDiagram([]string{fpath}, []string{}, iterateThroughSubFolders)
	if err != nil {
		log.Printf("[generate diagram] plantuml file generation error for path [%s]: %v", abs, err)
		return
	}

	rendered := result.Render()
	if len(rendered) > minimumPageSize {
		var filename = path + "-pumldiagram"

		var buf bytes.Buffer

		var headerstring = `<p><img src="` + filename + ".png" + `"/></img></p>`

		headerstring = markdown.URLConverter(node.root.id, headerstring, node.indexPage, node.id)

		var writer io.Writer

		writer, err = os.Create(node.path + "/" + filename + ".puml")
		if err != nil {
			log.Printf("[create file] plantuml file generation error for path [%s]: %v", abs, err)
			return
		}

		fmt.Fprint(&buf, headerstring)

		fmt.Fprint(writer, rendered)

		err := node.generatePlantumlImage(node.path + "/" + filename + ".puml")
		if err != nil {
			log.Printf("generatePlantumlImage error for path [%s]: %v", abs, err)
			return
		}

		log.Printf("uploading generated png file [%s] to page id of [%d]",
			node.path+"/"+filename+".png", node.id)

		node.uploadFile(node.path+"/"+filename+".png", node.indexPage)

		masterpagecontents := markdown.FileContents{
			MetaData: map[string]interface{}{
				"title": "plantuml-" + path,
			},
			Body: buf.Bytes(),
		}

		err = node.checkConfluencePages(&masterpagecontents, node.path+"/"+filename+".png")
		if err != nil {
			log.Printf("check confluence page error for path [%s]: %v", abs, err)
		}

		url := common.ConfluenceBaseURL + "/wiki/spaces/" +
			common.ConfluenceSpace + "/pages/" + func() string {
			return strconv.Itoa(node.id)
		}()

		log.Printf("a plantuml image was generated for location [%s] & is available at [%s]", fpath, url)
	}
}

// generatePlantumlImage method calls external application (plantuml.jar)
// in the docker container to generate the plantuml image (as a .png)
func (node *Node) generatePlantumlImage(fpath string) error {
	log.Printf("generating plantuml png from plantuml context provided by go code...")

	convertPlantuml := exec.Command("java", "-jar", "/app/plantuml.jar", "-tpng", fpath) // #nosec - pumlimage
	convertPlantuml.Stdout = os.Stdout

	err := convertPlantuml.Run()
	if err != nil {
		return fmt.Errorf("generatePlantumlImage error: %w", err)
	}

	return nil
}

// generatePage method creates a new page for node.
// and sets the parent page as the node root id.
// unless the node.root is nil in which case it is the root page
func (node *Node) generatePage(newPageContents *markdown.FileContents) error {
	_, abs := node.generateTitles()

	var isParentPage = true

	if nodeAPIClient == nil {
		return fmt.Errorf("error: confluence API client is nil")
	}

	var err error

	if node.root == nil {
		node.id, err = nodeAPIClient.CreatePage(0, newPageContents, isParentPage)
		if err != nil {
			return fmt.Errorf("create page error for folder path [%s]: %w", abs, err)
		}

		return nil
	}

	node.id, err = nodeAPIClient.CreatePage(node.root.id, newPageContents, !isParentPage)
	if err != nil {
		return fmt.Errorf("create page error for folder path [%s]: %w", abs, err)
	}

	return nil
}
