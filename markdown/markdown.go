// Package markdown provides a method for working with and parsing markdown documents
package markdown

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	m "gitlab.com/golang-commonmark/markdown"
)

// FileContents contains information from a file after being parsed from markdown.
// `Metadata` in the format of a `map[string]interface{}` this can contain title, description, slug etc.
// `Body` a `[]byte` that contains the resulting HTML after parsing the markdown and converting to HTML using Goldmark.
type FileContents struct {
	MetaData map[string]interface{}
	Body     []byte
}

// grabtitle collects the title of a markdown file
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

func newFileContents() *FileContents {
	f := FileContents{}
	f.MetaData = make(map[string]interface{})

	return &f
}

// ParseMarkdown is a function that uses external parsing library to grab markdown contents
// and return a filecontents object
func ParseMarkdown(rootID int, content []byte) (*FileContents, error) {
	r := bytes.NewReader(content)
	f := newFileContents()

	fmc, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		log.Println("error parsing front matter (using # title instead): %w", err)

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
		return nil, fmt.Errorf("page title is empty")
	}

	return f, nil
}

func stripFrontmatterReplaceURL(rootID int, content string) []byte {
	var pre string

	var frontmatter bool

	lines := strings.Split(content, "\n")

	for index := range lines {
		if strings.Contains(lines[index], "+++") {
			frontmatter = flip(frontmatter)
			continue
		}

		if strings.Contains(lines[index], "<img src=") {
			lines[index] = urlConverter(rootID, lines[index])
		}

		if !frontmatter {
			pre += lines[index] + "\n"
		}
	}

	pre = strings.TrimSpace(pre)

	return []byte(pre)
}

func flip(b bool) bool {
	return !b
}

// for images to be loaded in to confluence page, they must be in same directory as markdown to work
// this function replaces local url paths in html img links with a confluence path for folder page attachments
func urlConverter(rootID int, item string) string {
	sliceOne := strings.Split(item, `<p><img src="`)

	if len(sliceOne) > 1 {
		sliceTwo := strings.Split(sliceOne[1], `"`)

		if len(sliceTwo) > 1 {
			attachmentFileName := sliceTwo[0]
			rootPageID := strconv.Itoa(rootID)
			a := `<p><span class="confluence-embedded-file-wrapped">`
			b := `<img src="https://xiatech-markup.atlassian.net/wiki/download/attachments/`
			c := rootPageID + `/` + attachmentFileName + `"></img>`
			d := `</span></p>`

			return a + b + c + d
		}
	}

	return item
}
