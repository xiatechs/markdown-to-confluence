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

func TestGenerateTitles(t *testing.T) {
	node := Node{}
	node.path = "./folder/subfolder"
	dir, fullDir := node.generateTitles()

	if dir != "subfolder" {
		t.Errorf("got %s want %s", dir, "subfolder")
	}

	if fullDir != "folder/subfolder" {
		t.Errorf("got %s want %s", fullDir, "folder/subfolder")
	}
}
