package foldercrawler

//notodo: ignore this page
import (
	"testing"
)

func TestCrawl(t *testing.T) {
	c := New("test", "standard")
	c.Start("../../")
}
