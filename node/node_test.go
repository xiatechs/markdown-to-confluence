package node

//notodo: ignore this page
import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	markdown "github.com/xiatechs/markdown-to-confluence/markdown"
)

// this test lets you see visually how all the content is generated in case you want to debug the output locally
// basically run it against any path you want and you'll see the pages generated at the end - after logging
func TestStartDebugEverything(t *testing.T) {
	markdown.GrabAuthors = false

	node := Node{
		mu: &sync.RWMutex{},
	}

	m := mockclient{
		i: &iterator{},
	}

	SetAPIClient(m)

	if node.Start("../node") {
		node.Delete()
	}

	m.Print()

	expectedOutput := []Page{
		{
			title:  "file within readme!-downhere-file-node",
			isroot: false, body: "<h1>file within readme!</h1>\n<p>this is the file within the readme</p>",
			root: 3, id: 4,
		},
		{
			title:  "downhere-file-node",
			isroot: false,
			//nolint: lll // test data
			body: "<h1>a deeply nested readme file</h1>\n<p><span class=\"confluence-embedded-file-wrapped\"><img src=\"https://xiatech-markup.atlassian.net/wiki/download/attachments/3/picture.jpg\"></img></span></p>\n<p>this readme should be the index of a folder with another file in\nit called 'file-within-readme' or something similar</p>",
			root: 2,
			id:   3,
		},
		{
			title:  "file-testfolder2-node",
			isroot: false,
			//nolint: lll // test data
			body: "<h1>INDEX readme</h1>\n<p>here i am testing out using relative local links in github to product\nconfluence absolute links - as if by magic.</p>\n<p><a class=\"confluence-link\" href=\"/wiki/spaces//pages/5\" data-linked-resource-id=\"5\" data-base-url=\"https://xiatech-markup.atlassian.net/wiki\">lower level relative link</a></p>\n<p><a class=\"confluence-link\" href=\"/wiki/spaces//pages/1\" data-linked-resource-id=\"1\" data-base-url=\"https://xiatech-markup.atlassian.net/wiki\">test markdown folder</a></p>\n<p>[please start your links with https://]</p>\n<p>[please start your links with https://]</p>\n<p>[please start your links with https://]</p>\n<p>using fuzzy logic the links above to re-jigged to align with their correct source.\nit's not 100% foolproof but hey it's better than nothing, right?</p>\n<ul>\n<li>Tom Balcombe</li>\n</ul>",
			root: 0,
			id:   2,
		},
		{
			title:  "mtc-testpage-testfolder-node",
			isroot: false,
			//nolint: lll // test data
			body: "<h1>mtc-testpage</h1>\n<h2>This is a H2 line of text</h2>\n<h3>This is H3</h3>\n<p>This is some standard text</p>\n<pre><code>Here is some code formatted text\n</code></pre>\n<ul>\n<li><em>Sir Emailington</em> <a href=\"mailto:email@siremailington.com\">email@siremailington.com</a> - Chief Email Office</li>\n</ul>",
			root: 1,
			id:   1,
		},
		{
			title:  "testfolder-node",
			isroot: true,
			//nolint: lll // test data
			body: "<p>Welcome to the '<b>testfolder-node</b>' folder of this Xiatech code repo.</p>\n\t\t<p>This folder full path in the repo is: node/testfolder</p>\n<p>You will find attachments/images for this folder via the ellipsis at the top right.</p>\n<p>Any markdown or subfolders is available in children pages under this page.</p>",
			root: 0,
			id:   0,
		},
	}

	assert.Equal(t, expectedOutput, m.GetPages())
}
