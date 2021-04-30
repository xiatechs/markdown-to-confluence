package node

import (
	"testing"
)

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

func TestCheckreadme(t *testing.T) {
	// run function on it's own
	node := Node{}
	node.checkreadme("fakefolder")
}

func TestCheckpuml(t *testing.T) {
	// run function on it's own
	node := Node{}
	node.alive = true
	node.checkpuml("fake", "hello.png")
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
