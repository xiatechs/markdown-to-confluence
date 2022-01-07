package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/common"
)

func TestFuzzyLogic(t *testing.T) {
	testInputs := []struct {
		name           string
		arg1           string
		arg2           map[string]string
		expectedoutput string
	}{
		{
			name: "local most likely",
			arg1: `<p><a href="../hello/there">a_page</a></p>`,
			arg2: map[string]string{
				"./there/hello/item.txt":                                      "1",
				"./hello/there/readme.md":                                     "2",
				"../../markdown/hello/there/readme.md":                        "3",
				"../../../anotherfolder/hello/markdown/hello/there/readme.md": "4",
			},
			//nolint:lll /// test data
			expectedoutput: "<p><a class=\"confluence-link\" href=\"/wiki/spaces//pages/2\" data-linked-resource-id=\"2\" data-base-url=\"https://xiatech-markup.atlassian.net/wiki\">a_page</a></p>",
		},
		{
			name: "distant markdown link most likely",
			arg1: `<p><a href="./node">a_page</a></p>`,
			arg2: map[string]string{
				"./there/hello/item.txt":                                      "1",
				"./hello/there/readme.md":                                     "2",
				"../../markdown/hello/there/readme.md":                        "3",
				"../../../anotherfolder/hello/markdown/hello/there/readme.md": "4",
				"../../../node":                                               "5",
			},
			//nolint:lll /// test data
			expectedoutput: "<p><a class=\"confluence-link\" href=\"/wiki/spaces//pages/5\" data-linked-resource-id=\"5\" data-base-url=\"https://xiatech-markup.atlassian.net/wiki\">a_page</a></p>",
		},
		{
			name: "can't locate a valid link",
			arg1: `<p><a href="../hello/there/item.txt">a_page</a></p>`,
			arg2: map[string]string{
				"./there/hello/item.txt":                                "1",
				"./ayy/there/hello/readme.md":                           "2",
				"../../there/hello/mark/readme.md":                      "3",
				"../../../anotherfolder/hello/markdown/there/readme.md": "4",
			},
			expectedoutput: "<p>[please start your links with https://]</p>",
		},
	}

	for _, test := range testInputs {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedoutput, fuzzyLogicURLdetector(test.arg1, test.arg2))
		})
	}
}

func TestParagraphify(t *testing.T) {
	input := `code line a
code line b
code line c`

	expected := `<h3>To view this try copy&amp;paste to this site: <a href="https://www.planttext.com/">PlainText UML Editor</a></h3>` + //nolint:lll // it's long test string
		`
<h3>Alternatively please install a <em>PlantUML Visualizer plugin</em> for Chrome or Firefox</h3>
<pre><code class="language-+">code line a
code line b
code line c
</code></pre>
`

	output := Paragraphify(input)
	assert.Equal(t, expected, output)
}

func TestParseMarkDown(t *testing.T) {
	link := common.ConfluenceBaseURL + `/wiki/download/attachments/0/node.png`
	testInputs := []struct {
		Name     string
		input    []byte
		expected *FileContents
	}{
		{
			Name: "title & no URL",
			input: []byte(`# Markdown to Confluence Action

This Action will trawl through a repository.`),
			expected: &FileContents{
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
			expected: &FileContents{
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
			result, _ := ParseMarkdown(0, test.input, false, 0, map[string]string{}, ".")
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestParseMarkdown_HappyPath(t *testing.T) {
	testContent := []byte(`
+++
categories = ["Development", "Github Actions"]
date = "2021-03-10"
description = "A guide on how to use the markdown to confluence action"
slug = "markdown-to-confluence-guide"
title = "Markdown to Confluence Action Guide"
+++

# Test Content 
test description`)

	expectOutput := &FileContents{
		MetaData: map[string]interface{}{
			"categories":  []interface{}{"Development", "Github Actions"},
			"date":        "2021-03-10",
			"description": "A guide on how to use the markdown to confluence action",
			"slug":        "markdown-to-confluence-guide",
			"title":       "Markdown to Confluence Action Guide",
		},
		Body: []byte(`<h1>Test Content</h1>
<p>test description</p>`),
	}

	out, err := ParseMarkdown(0, testContent, false, 0, map[string]string{}, ".")
	assert.Nil(t, err)
	assert.Equal(t, out, expectOutput)
}

func TestParseMarkdown_MalformedFrontMatter(t *testing.T) {
	testContent := []byte(`
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

	_, err := ParseMarkdown(0, testContent, false, 0, map[string]string{}, "markdown_test.go")
	assert.NotNil(t, err)
}
