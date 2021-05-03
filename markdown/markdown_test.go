package markdown_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestParseMarkDown(t *testing.T) {
	input := []byte(`# Markdown to Confluence Action

	This Action will trawl through a repository`)

	expected := &markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "Markdown to Confluence Action",
		},
		Body: []byte(`<h1>Markdown to Confluence Action</h1>
<pre><code>This Action will trawl through a repository</code></pre>

`),
	}

	result, _ := markdown.ParseMarkdown(0, input)

	assert.Equal(t, expected, result)
}
