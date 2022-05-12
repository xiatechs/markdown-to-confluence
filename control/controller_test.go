package control

import (
	"os"
	"testing"

	"github.com/xiatechs/markdown-to-confluence/common"
)

func TestCrawl(t *testing.T) {
	_ = os.Setenv("CONFLUENCE_BASE_URL", "https://xiatech.atlassian.net")
	_ = os.Setenv("CONFLUENCE_SPACE", "MTW")
	_ = os.Setenv("CONFLUENCE_USERNAME", "apps.markdown@xiatech.co.uk")
	_ = os.Setenv("CONFLUENCE_API_KEY", "")

	common.Refresh()

	c := New("template", "standard")

	c.Start("./testRoot")
}

/* TODO: create mock test using interfaces
func TestFolderCrawler(t *testing.T) {
	ctrl := gomock.NewController(t)

	api := apihandler.NewMockApiController(ctrl)

	fh := filehandler.NewMockFileHandler(ctrl)

	c := NewDI(fh, api)

	fileContents := &filehandler.FileContents{}
	parentMap := make(map[string]interface{})
	var input string

	//gomock.InOrder(
	fh.EXPECT().ConvertMarkdown(input, input, parentMap).Return(fileContents, nil)
	//)

	c.Start("./testfolder")
}
*/
