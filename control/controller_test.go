package control

import (
	"os"
	"testing"
)

func TestCrawl(t *testing.T) {
	_ = os.Setenv("CONFLUENCE_BASE_URL", ":)")
	_ = os.Setenv("CONFLUENCE_SPACE", ":0")
	_ = os.Setenv("CONFLUENCE_USERNAME", ":|")
	_ = os.Setenv("CONFLUENCE_API_KEY", ":<")

	c := New("test", "standard")

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
