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

func TestIterate(t *testing.T) {
	node := Node{}

	node.path = "./"

	boolean := node.iterate(false, false)
	if boolean != false {
		t.Errorf("got %t want %t", boolean, false)
	}

	boolean = node.iterate(true, true)
	if boolean != false {
		t.Errorf("got %t want %t", boolean, false)
	}

	boolean = node.iterate(true, false)
	if boolean != false {
		t.Errorf("got %t want %t", boolean, false)
	}

	boolean = node.iterate(false, true)
	if boolean != false {
		t.Errorf("got %t want %t", boolean, false)
	}
}
