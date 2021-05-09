package node

import (
	"testing"
)

// TODO add more test coverage

func TestStart(t *testing.T) {
	node := Node{}
	want := false

	got := node.Start("./fakedirectory", nil)
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}

	nodeOne := Node{}

	nodeTest := Node{
		root: &nodeOne,
	}
	want = false

	got = nodeTest.Start("./fakedirectory", nil)
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestIsFolder(t *testing.T) {
	isFolder("hello")
}

func TestUploadFile(t *testing.T) {
	node := Node{}

	node.uploadFile("")
}
