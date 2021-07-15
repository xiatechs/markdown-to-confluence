package node

//notodo: ignore this page
import (
	"testing"

	"github.com/golang/mock/gomock"
	confluence "github.com/xiatechs/markdown-to-confluence/confluence"
	markdown "github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestStartAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockAPIClienter(ctrl)

	results := confluence.PageResults{}

	page := markdown.FileContents{}

	node := Node{}

	SetAPIClient(client)

	gomock.InOrder(
		client.EXPECT().FindPage("node", false).Times(1).Return(&results, nil),
		client.EXPECT().FindPage("markdown-to-confluence/node+readme", false).Times(1).Return(&results, nil),
		client.EXPECT().CreatePage(0, &page, true).Times(1).Return(0, nil),
	)

	node.Start("../node")
}

func TestStartBrandNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockAPIClienter(ctrl)

	nodePage := markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "node",
		},
		Body: []byte(`<p>Welcome to the '<b>node</b>' folder of this Xiatech code repo.</p>
		<p>This folder full path in the repo is: node</p>
<p>You will find attachments/images for this folder via the ellipsis at the top right.</p>
<p>Any markdown or subfolders is available in children pages under this page.</p>`),
	}
	readmePage := markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "markdown-to-confluence/node readme",
		},
		Body: []byte(`<h1>markdown-to-confluence/node readme</h1>
<h2>the node package is to enable reading through a repo and create a tree of content on confluence</h2>
<h3>The package contains one exported struct:</h3>
<pre><code>// Node struct enables creation of a page tree
Node{}
</code></pre>
<h3>The package contains two exported methods:</h3>
<pre><code>Node{}

// this method begins the generation of a page tree in confluence for a repo project path.
// it ruturns a boolean confirming 'is projectPath a valid folder path'.
Start(projectPath string, client *confluence.APIClient) bool

// this method begins the deletion of pages in confluence that do not exist in
// local repository project path - it can be called after Instantiate method is called and returns true.
Delete()
</code></pre>`),
	}

	node := Node{}

	SetAPIClient(client)

	gomock.InOrder(
		client.EXPECT().FindPage("node", false).Times(1).Return(nil, nil),
		client.EXPECT().CreatePage(0, &nodePage, true).Times(1).Return(0, nil),
		client.EXPECT().FindPage("markdown-to-confluence/node+readme", false).Times(1).Return(nil, nil),
		client.EXPECT().CreatePage(0, &readmePage, false).Times(1).Return(0, nil),
	)

	if node.Start("../node") {
		node.Delete()
	}
}

func TestStart(t *testing.T) {
	node := Node{}
	want := false

	got := node.Start("./fakedirectory")
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}

	nodeOne := Node{}

	nodeTest := Node{
		root: &nodeOne,
	}
	want = false

	got = nodeTest.Start("./fakedirectory")
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
