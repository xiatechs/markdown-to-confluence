package markdown_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestParseMarkDown(t *testing.T) {
	link := `https://xiatech-markup.atlassian.net/wiki/download/attachments/0/node.png`
	testInputs := []struct {
		Name     string
		input    []byte
		expected *markdown.FileContents
	}{
		{
			Name: "title & no URL",
			input: []byte(`# Markdown to Confluence Action

This Action will trawl through a repository.`),
			expected: &markdown.FileContents{
				MetaData: map[string]interface{}{
					"title": "Markdown to Confluence Action",
				},
				Body: []byte(`<h1>Markdown to Confluence Action</h1>
<p>This Action will trawl through a repository.</p>`),
			},
		},
		{
			Name: "title & URL",
			input: []byte(`# Markdown to Confluence Action

![Diagram of action methodology](node.png)`),
			expected: &markdown.FileContents{
				MetaData: map[string]interface{}{
					"title": "Markdown to Confluence Action",
				},
				Body: []byte(`<h1>Markdown to Confluence Action</h1>
<p><span class="confluence-embedded-file-wrapped"><img src="">` + link + `</img></span></p>`),
			},
		},
	}

	for _, test := range testInputs {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			result, _ := markdown.ParseMarkdown(0, test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
