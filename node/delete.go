package node

// delete - methods regarding deleting pages in confluence wiki

import (
	"log"
	"strconv"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

// used to verify whether pages need to be deleted or not
var masterTitles []string

// findPagesToDelete method grabs results of page to begin deleting
func (node *Node) findPagesToDelete(id string) {
	findParentPageAndChildren := true

	children, err := nodeAPIClient.FindPage(id, findParentPageAndChildren)
	if err != nil {
		log.Printf("error finding page: %s", err)
	}

	if children != nil {
		node.deletePages(children)
	}
}

// deletePages method is to find a page to delete
// and any children pages that might need to be deleted
func (node *Node) deletePages(children *confluence.PageResults) {
	for index := range children.Results {
		var noDelete bool

		for index2 := range masterTitles {
			if children.Results[index].Title == masterTitles[index2] {
				noDelete = true
				break
			}
		}

		if !noDelete {
			node.findPagesToDelete(children.Results[index].ID)
			node.deletePage(children.Results[index].ID)
		}
	}
}

// convert id to integer to pass to the API method DeletePage
func (node *Node) deletePage(id string) {
	convert, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error getting page ID: %s", err)
		return
	}

	err = nodeAPIClient.DeletePage(convert)
	if err != nil {
		log.Printf("error deleting page: %s", err)
	}
}
