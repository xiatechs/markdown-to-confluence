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

// newFileContents function creates a new 'filecontents' object
func newFileContents() *markdown.FileContents {
	f := markdown.FileContents{}
	f.MetaData = make(map[string]interface{})

	return &f
}

// grabTODO function takes in go code (content) and a filename
// and returns a string containing a list of all the TODO's in a piece of code
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

// ParseGo function is to parse a .go file for todo rows
// and appends output to the collator string variable
func ParseGo(content []byte, filename string) {
	collator += grabTODO(string(content), filename)
}

// GenerateTODO function takes in rootDir foldername and
// percentage of markdown content in folder as input
// and returns a page
func GenerateTODO(rootDir, percentage string) *markdown.FileContents {
	f := newFileContents()
	f.MetaData["title"] = "More info on '" + rootDir + "' repo"

	md := m.New(
		m.HTML(true),
		m.Tables(true),
		m.Linkify(true),
		m.Typographer(false),
		m.XHTMLOutput(true),
	)

	collator = "## " + percentage + "\n" + collator
	preformatted := md.RenderToString([]byte(collator))

	f.Body = []byte(preformatted)

	return f
}
