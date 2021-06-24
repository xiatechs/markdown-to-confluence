package node

import (
	"testing"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

func TestDelete(t *testing.T) {
	node := Node{}
	node.id = 1
	node.root = &node
	node.Delete()

	node.id = 0
	node.root = &node
	node.Delete()
}

func TestDeletePages(t *testing.T) {
	node := Node{}

	c := confluence.PageResults{}

	node.deletePages(&c)
}

func TestDeletePage(t *testing.T) {
	node := Node{}

	node.deletePage("")
}
