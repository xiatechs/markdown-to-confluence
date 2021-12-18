package node

// delete - methods regarding deleting pages in confluence wiki

import (
	"log"
	"strconv"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/confluence"
)

// findPagesToDelete method grabs results of page to begin deleting
func (node *Node) findPagesToDelete(id string) {
	findParentPageAndChildren := true

	if nodeAPIClient != nil {
		children, err := nodeAPIClient.FindPage(id, findParentPageAndChildren)
		if err != nil {
			log.Printf("error finding page: %s", err)
		}

		if children != nil {
			node.deletePages(children)
		}
	}
}

// deletePages method is to find a page to delete
// and any children pages that might need to be deleted
func (node *Node) deletePages(children *confluence.PageResults) {
	for index := range children.Results {
		var noDelete bool

		for index2 := range node.titles {
			if children.Results[index].Title == node.titles[index2] {
				noDelete = true
				break
			}
		}

		if !noDelete {
			node.findPagesToDelete(children.Results[index].ID)

			go node.deletePage(children.Results[index].ID)
		}
	}

	log.Println("Here are the pages:")

	for path, id := range t.branches {
		log.Println(path, "|", common.ConfluenceBaseURL+"/wiki/spaces/"+common.ConfluenceSpace+"/pages/"+id)
	}
}

// deletePage method converts id to integer to pass to the API method DeletePage
// this method can be run concurrently with no wait needed:
// the pages are deleted by ID and don't need parent page reference
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
