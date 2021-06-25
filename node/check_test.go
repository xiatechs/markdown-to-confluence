package node

import (
	"log"
	"testing"

	"github.com/xiatechs/markdown-to-confluence/confluence"
)

func TestCheckIfMarkDown(t *testing.T) {
	node := Node{}

	b := node.checkIfMarkDown("fakefolder", false)
	if b {
		t.Errorf("got %t want %t", b, false)
	}

	b = node.checkIfMarkDown("fakefolder", true)
	if b {
		t.Errorf("got %t want %t", b, false)
	}

	if node.alive {
		t.Errorf("got %t want %t", true, false)
	}
}

func TestCheckIfMarkDownFile(t *testing.T) {
	node := Node{}

	output := node.checkIfMarkDownFile(true, "test.md")
	if output != true {
		t.Errorf("got %t want %t", output, true)
	}

	output = node.checkIfMarkDownFile(false, "test.md")
	if output != true {
		t.Errorf("got %t want %t", output, true)
	}

	output = node.checkIfMarkDownFile(true, "test.notmd")
	if output != false {
		t.Errorf("got %t want %t", output, false)
	}

	output = node.checkIfMarkDownFile(false, "test.notmd")
	if output != false {
		t.Errorf("got %t want %t", output, false)
	}
}

func TestCheckForImages(t *testing.T) {
	node := Node{}
	node.root = &node
	node.alive = true
	node.checkForImages("test.gif")
}

func TestCheckIfGoFile(t *testing.T) {
	node := Node{}
	node.root = &node
	node.alive = true
	node.checkIfGoFile("test.go", false)
}

func TestProcessGoFile(t *testing.T) {
	node := Node{}
	node.root = &node
	node.alive = true
	node.processGoFile("test.go")
}

func TestCheckPageID(t *testing.T) {
	node := Node{}

	data := confluence.PageResults{}

	err := node.checkPageID(data)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}

func TestFileInDirectoryCheck(t *testing.T) {
	node := Node{}
	node.fileInDirectoryCheck("fake", true, true)
	node.fileInDirectoryCheck("fake", true, false)
	node.fileInDirectoryCheck("fake", false, true)
	node.fileInDirectoryCheck("fake", false, false)
}

func TestCheckIfRootAlive(t *testing.T) {
	node := Node{}
	node.alive = false
	node.root = &node
	node.checkIfRootAlive("fake")

	node = Node{}
	node.alive = true
	node.root = &node
	node.checkIfRootAlive("fake")

	node = Node{}
	node.alive = false
	node.root = nil
	node.checkIfRootAlive("fake")

	node = Node{}
	node.alive = true
	node.root = nil
	node.checkIfRootAlive("fake")
}
