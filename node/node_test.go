package node

import (
	"testing"
)

// TODO add more test coverage

func TestInstantiate(t *testing.T) {
	node := Node{}
	want := false

	got := node.Instantiate("./fakedirectory")
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}

	nodeOne := Node{}

	nodeTest := Node{
		root: &nodeOne,
	}
	want = false

	got = nodeTest.Instantiate("./fakedirectory")
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestCheckMarkDown(t *testing.T) {
	node := Node{}
	node.checkMarkDown(false, "fakefolder")
	node.checkMarkDown(true, "fakefolder")
}

func TestCheckOtherFiles(t *testing.T) {
	node := Node{}
	node.checkOtherFiles(true, "fake.png")
	node.checkOtherFiles(false, "fake.png")
}

func TestScrub(t *testing.T) {
	node := Node{}
	node.Scrub()
}

func TestCheckAll(t *testing.T) {
	node := Node{}
	node.checkAll(false, "fake")
	node.checkAll(true, "fake")
}

func TestGenerateMaster(t *testing.T) {
	node := Node{}
	node.generateMaster()
}

func TestVerifyCreateNode(t *testing.T) {
	node := Node{}
	node.alive = false
	node.verifyCreateNode("fake")

	node = Node{}
	node.alive = true
	node.root = nil
	node.verifyCreateNode("fake")

	node = Node{}
	node.alive = true
	node.root = &node
	node.verifyCreateNode("fake")

	node = Node{}
	node.alive = false
	node.root = &node
	node.verifyCreateNode("fake")
}

func TestRemoveFirstByte(t *testing.T) {
	output := removefirstbyte("")
	if output != "" {
		t.Errorf("got %s want %s", output, "")
	}

	output = removefirstbyte("1")
	if output != "1" {
		t.Errorf("got %s want %s", output, "1")
	}
	
	output = removefirstbyte("/1")
	if output != "1" {
		t.Errorf("got %s want %s", output, "1")
	}
}

func TestIsFolder(t *testing.T) {
	isFolder("hello")
}
