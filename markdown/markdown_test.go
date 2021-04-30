package markdown_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestParseMarkDown(t *testing.T) {
	input := []byte(`
# Title

This is some content.
	`)

	expected := &markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "Title",
		},
		Body: []byte(`<h1>Title</h1>
<p>This is some content.</p>
`),
	}

	result, _ := markdown.ParseMarkdown(0, input)

	fmt.Println(string(expected.Body))
	fmt.Println(string(result.Body))
	assert.Equal(t, expected, result)
}
