package node

import (
	"testing"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

func TestDelete(t *testing.T) {
	node := Node{}

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
