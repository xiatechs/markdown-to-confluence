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

	node := Node{}

	SetAPIClient(client)

	gomock.InOrder(
		client.EXPECT().FindPage("testfolder", false).Times(1).Return(&results, nil),
		client.EXPECT().FindPage("mtc-testpage-testfolder-nodetestfolder", false).Times(1).Return(&results, nil),
	)

	if node.Start("../node/testfolder") {
		node.Delete()
	}
}

func TestStartBrandNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := NewMockAPIClienter(ctrl)

	testfolderFolderpage := markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "testfolder",
		},
		Body: []byte(`<p>Welcome to the '<b>testfolder</b>' folder of this Xiatech code repo.</p>
		<p>This folder full path in the repo is: node/testfolder</p>
<p>You will find attachments/images for this folder via the ellipsis at the top right.</p>
<p>Any markdown or subfolders is available in children pages under this page.</p>`),
	}

	testPage := markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "mtc-testpage-testfolder-nodetestfolder",
		},
		Body: []byte(`<h1>mtc-testpage</h1>
<h2>This is a H2 line of text</h2>
<h3>This is H3</h3>
<p>This is some standard text</p>
<pre><code>Here is some code formatted text
</code></pre>
<ul>
<li><em>Sir Emailington</em> <a href="mailto:email@siremailington.com">email@siremailington.com</a> - Chief Email Office</li>
</ul>`),
	}

	node := Node{}

	SetAPIClient(client)

	gomock.InOrder(
		client.EXPECT().FindPage("testfolder", false).Times(1).Return(nil, nil),
		client.EXPECT().CreatePage(0, &testfolderFolderpage, true).Times(1).Return(0, nil),
		client.EXPECT().FindPage("mtc-testpage-testfolder-nodetestfolder", false).Times(1).Return(nil, nil),
		client.EXPECT().CreatePage(0, &testPage, false).Times(1).Return(0, nil),
	)

	if node.Start("../node/testfolder") {
		node.Delete()
	}
}
