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
<p><span class="confluence-embedded-file-wrapped"><img src="` + link + `"></img></span></p>`),
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


func TestParseMarkdown_HappyPath(t *testing.T) {
	testContent := strings.NewReader(`
+++
categories = ["Development", "Github Actions"]
date = "2021-03-10"
description = "A guide on how to use the markdown to confluence action"
slug = "markdown-to-confluence-guide"
title = "Markdown to Confluence Action Guide"
+++

# Test Content 
test description`)

	expectOutput := &markdown.FileContents{
		MetaData: map[string]interface{}{
			"categories":  []interface{}{"Development", "Github Actions"},
			"date":        "2021-03-10",
			"description": "A guide on how to use the markdown to confluence action",
			"slug":        "markdown-to-confluence-guide",
			"title":       "Markdown to Confluence Action Guide",
		},
		Body: []byte(`<h1>Test Content</h1>
<p>test description</p>
`),
	}

	out, err := markdown.ParseMarkdown(testContent)
	assert.Nil(t, err)
	assert.Equal(t, out, expectOutput)
}

func TestParseMarkdown_MalformedFrontMatter(t *testing.T) {
	testContent := strings.NewReader(`
	+++
	badFrontMatter = 253svsasrg
	categories = ["Development", "Github Actions"]
	date = "2021-03-10"
	description = "A guide on how to use the markdown to confluence action"
	slug = "markdown-to-confluence-guide"
	title = "Markdown to Confluence Action Guide"
	+++
	
	# Test Content 
	test description`)

	_, err := markdown.ParseMarkdown(testContent)
	assert.NotNil(t, err)
}

