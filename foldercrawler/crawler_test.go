package foldercrawler

//notodo: ignore this page
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

	c.Start("../../")

}
