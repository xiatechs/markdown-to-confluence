// Package todo is to gather all TODO's in a repo and store them in one file for people to look through
package todo

import (
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/markdown"
	m "gitlab.com/golang-commonmark/markdown"
)

var collator string

// MainTODOPage will be the todo page that will be uploaded to root folder
var MainTODOPage markdown.FileContents

func newFileContents() *markdown.FileContents {
	f := markdown.FileContents{}
	f.MetaData = make(map[string]interface{})

	return &f
}

func grabTODO(content, filename string) string {
	var output string

	var containsTODO bool

	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")

	for index := range lines {
		if strings.Contains(lines[index], `TODO`) {
			output += "## Filename: " + filename + "\n"
			output += "\n"
			containsTODO = true

			break
		}
	}

	if !containsTODO {
		return ""
	}

	for index := range lines {
		if strings.Contains(lines[index], `TODO`) {
			rownumber := strconv.Itoa(index + 1)
			output += "Row: <" + rownumber + "> " + lines[index]
			output += "\n\n"
		}
	}

	output += "\n"

	return output
}

// ParseGo is to parse a .go file for todo rows
func ParseGo(content []byte, filename string) {
	collator += grabTODO(string(content), filename)
}

// GenerateTODO lint
func GenerateTODO(rootDir string) *markdown.FileContents {
	f := newFileContents()
	f.MetaData["title"] = "TODO list for '" + rootDir + "' repo"

	md := m.New(
		m.HTML(true),
		m.Tables(true),
		m.Linkify(true),
		m.Typographer(false),
		m.XHTMLOutput(true),
	)

	var preformatted string

	if preformatted = md.RenderToString([]byte(collator)); preformatted == "" {
		return nil
	}

	f.Body = []byte(preformatted)

	return f
}
