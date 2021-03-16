package markdown_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

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
