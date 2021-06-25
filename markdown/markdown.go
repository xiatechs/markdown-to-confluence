// Package markdown provides a method for working with and parsing markdown documents
package markdown

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/xiatechs/markdown-to-confluence/common"
	m "gitlab.com/golang-commonmark/markdown"
)

// FileContents contains information from a file after being parsed from markdown.
// `Metadata` in the format of a `map[string]interface{}` this can contain title, description, slug etc.
// `Body` a `[]byte` that contains the resulting HTML after parsing the markdown and converting to HTML using Goldmark.
type FileContents struct {
	MetaData map[string]interface{}
	Body     []byte
}

// grabtitle function collects the title of a markdown file
// and returns it as a string
func grabtitle(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")

	for index := range lines {
		if len(lines[index]) > 1 {
			if lines[index][0] == '#' {
				return strings.TrimSpace(strings.ReplaceAll(lines[index], "#", ""))
			}
		}
	}

	return ""
}

// newFileContents function creates a new filecontents object
func newFileContents() *FileContents {
	f := FileContents{}
	f.MetaData = make(map[string]interface{})

	return &f
}

// Paragraphify takes in a file contents and returns
// a formatted HTML page as a string
func Paragraphify(content string) string {
	var pre string
	pre += "``` + \n"
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")

	for index := range lines {
		pre += lines[index] + "\n"
	}

	pre += "```"

	md := m.New(
		m.HTML(true),
		m.Tables(true),
		m.Linkify(true),
		m.Typographer(false),
		m.XHTMLOutput(true),
	)

	preformatted := md.RenderToString([]byte(pre))

	return preformatted
}

// ParseMarkdown function uses external parsing library to grab markdown contents
// and return a filecontents object
func ParseMarkdown(rootID int, content []byte) (*FileContents, error) {
	r := bytes.NewReader(content)
	f := newFileContents()

	fmc, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		log.Println("issue parsing frontmatter - (using # title instead): %w", err)

		f.MetaData["title"] = grabtitle(string(content))
	} else {
		if len(fmc.FrontMatter) != 0 {
			f.MetaData = fmc.FrontMatter
		} else {
			f.MetaData["title"] = grabtitle(string(content))
		}
	}

	md := m.New(
		m.HTML(true),
		m.Tables(true),
		m.Linkify(true),
		m.Typographer(false),
		m.XHTMLOutput(true),
	)

	preformatted := md.RenderToString(content)
	f.Body = stripFrontmatterReplaceURL(rootID, preformatted)

	if f.MetaData["title"] == "" {
		return nil, fmt.Errorf("markdown parsing error - page title is empty")
	}

	return f, nil
}

// stripFrontmatterReplaceURL function takes in parent page ID and
// markdown file contents and removes TOML frontmatter, and replaces
// local URL with relative confluence URL
func stripFrontmatterReplaceURL(rootID int, content string) []byte {
	var pre string

	var frontmatter bool

	lines := strings.Split(content, "\n")

	for index := range lines {
		if strings.Contains(lines[index], "+++") {
			frontmatter = flip(frontmatter)
			continue
		}

		// temporary solution to local url path issue - remove them
		if strings.Contains(lines[index], "<a href=") && !strings.Contains(lines[index], "https://") {
			lines[index] = "<p></p>"
		}

		if strings.Contains(lines[index], "<img src=") {
			lines[index] = URLConverter(rootID, lines[index])
		}

		if !frontmatter {
			pre += lines[index] + "\n"
		}
	}

	pre = strings.TrimSpace(pre)

	return []byte(pre)
}

// flip function returns the opposite of bool
func flip(b bool) bool {
	return !b
}

// URLConverter function for images to be loaded in to confluence page
// (they must be in same directory as markdown to work)
// this function replaces local url paths in html img links
// with a confluence path for folder page attachments on parent page
func URLConverter(rootID int, item string) string {
	sliceOne := strings.Split(item, `<p><img src="`)

	if len(sliceOne) > 1 {
		sliceTwo := strings.Split(sliceOne[1], `"`)

		if len(sliceTwo) > 1 {
			attachmentFileName := sliceTwo[0]
			rootPageID := strconv.Itoa(rootID)
			a := `<p><span class="confluence-embedded-file-wrapped">`
			b := `<img src="` + common.ConfluenceBaseURL + `/wiki/download/attachments/`
			c := rootPageID + `/` + attachmentFileName + `"></img>`
			d := `</span></p>`

			return a + b + c + d
		}
	}

	return item
}
