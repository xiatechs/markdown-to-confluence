// Package markdown provides a method for working with and parsing markdown documents
package markdown

import (
	"strings"

	m "gitlab.com/golang-commonmark/markdown"
)

// FileContents contains information from a file after being parsed from markdown.
// `Metadata` in the format of a `map[string]interface{}` this can contain title, description, slug etc.
// `Body` a `[]byte` that contains the resulting HTML after parsing the markdown and converting to HTML using Goldmark.
type FileContents struct {
	MetaData map[string]interface{}
	Body     []byte
}

func grabtitle(content string) string {
	lines := strings.Split(content, "\n")
	for index := range lines {
		if len(lines[index]) != 0 {
			if lines[index][0] == '#' && len(lines[index]) > 1 {
				return strings.ReplaceAll(lines[index], "#", "")
			}
		}
	}

	return "no title found"
}

// ParseMarkdown is a function that uses external parsing library to grab markdown contents
func ParseMarkdown(rootID int, content []byte) (*FileContents, error) {
	f := FileContents{
		MetaData: map[string]interface{}{
			"title": "",
		},
	}
	md := m.New(m.XHTMLOutput(true))
	f.MetaData["title"] = grabtitle(string(content))
	preformatted := md.RenderToString(content)
	f.Body = []byte(preformatted)

	return &f, nil
}

// TODO - Figure out how to embed images into a conference page.

/*
func convertImageURLtoRootLink(rootID int, content string) []byte {
	var pre string
	lines := strings.Split(content, "\n")
	for index := range lines {
		if strings.Contains(lines[index], "<img src=") {
			lines[index] = switcheroo(rootID, lines[index])
		}
		pre += lines[index] + "\n"
	}
	return []byte(pre)
}

func switcheroo(rootID int, item string) string {
	item1 := strings.Split(item, `<p><img src="`)
	if len(item1) > 1 {
		item2 := strings.Split(item1[1], `"`)[0]
		number := strconv.Itoa(rootID)
		return `<p><span class="confluence-embedded-file-wrapped">
		<img src="/download/attachments/` + number + `/` + item2 + `></img></span></p>`
	}
	return "<p>img could not be parsed</p>"

	// https://xiatech-markup.atlassian.net/wiki/spaces/~802150356
}
*/
