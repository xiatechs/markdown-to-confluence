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
		rootFolder: &nodeOne,
	}
	want = false

	got = nodeTest.Instantiate("./fakedirectory")
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestPrintOverview(t *testing.T) {
	// run function on it's own
	PrintOverview()
}

func TestCheckMarkDown(t *testing.T) {
	// run function on it's own
	node := Node{}
	node.checkMarkDown("fakefolder")
}

func TestCheckOtherFiles(t *testing.T) {
	// run function on it's own
	node := Node{}
	node.alive = true
	node.checkOtherFiles("fake.png")
}

func TestCheckAll(t *testing.T) {
	node := Node{}
	node.checkAll("fake")
}

func TestGenerateMaster(t *testing.T) {
	node := Node{}
	node.generateMaster()
}

func TestVerifyCreateNode(t *testing.T) {
	node := Node{}
	node.verifyCreateNode("fake")
}
