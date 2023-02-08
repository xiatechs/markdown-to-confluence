package node

//notodo: ignore this page
import (
	"sync"
	"testing"

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

	t.Skip() // skip test as concurrency means it fails - only used locally for debugging

	if node.Start(0, "../node", false) {
		node.Delete()
	}

	m.Print()
}
