package node

import (
	"testing"

	"github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestGenerateMaster(t *testing.T) {
	node := Node{}
	node.generateMaster()
}

func TestGeneratePage(t *testing.T) {
	node := Node{}
	newPageContents := markdown.FileContents{}

	node.generatePage(&newPageContents)
}
