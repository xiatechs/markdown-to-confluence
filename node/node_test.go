package node

import (
	"log"
	"testing"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

// TODO add more test coverage

func TestInstantiate(t *testing.T) {
	node := Node{}
	want := false

	got := node.Instantiate("./fakedirectory", nil)
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}

	nodeOne := Node{}

	nodeTest := Node{
		root: &nodeOne,
	}
	want = false

	got = nodeTest.Instantiate("./fakedirectory", nil)
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestCheckMarkDown(t *testing.T) {
	node := Node{}

	b := node.checkMarkDown(false, "fakefolder")
	if b != false {
		t.Errorf("got %t want %t", b, false)
	}

	b = node.checkMarkDown(true, "fakefolder")
	if b != false {
		t.Errorf("got %t want %t", b, false)
	}
}

func TestGrabPageData(t *testing.T) {
	node := Node{}

	data := confluence.PageResults{}

	err := node.grabpagedata(data)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
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
	node.root = &node
	node.verifyCreateNode("fake")

	node = Node{}
	node.alive = true
	node.root = &node
	node.verifyCreateNode("fake")

	node = Node{}
	node.alive = false
	node.root = nil
	node.verifyCreateNode("fake")

	node = Node{}
	node.alive = true
	node.root = nil
	node.verifyCreateNode("fake")
}

func TestIsFolder(t *testing.T) {
	isFolder("hello")
}

func TestCheckConfluencePages(t *testing.T) {
	node := Node{}

	newPageContents := markdown.FileContents{}

	node.checkConfluencePages(&newPageContents)
}

func TestDeletePage(t *testing.T) {
	node := Node{}

	node.deletePage("")
}

func TestGeneratePage(t *testing.T) {
	node := Node{}
	newPageContents := markdown.FileContents{}

	node.generatePage(&newPageContents)
}

func TestUploadFile(t *testing.T) {
	node := Node{}

	node.uploadFile("")
}

func TestDeletePages(t *testing.T) {
	node := Node{}

	c := confluence.PageResults{}

	node.deletePages(&c)
}
