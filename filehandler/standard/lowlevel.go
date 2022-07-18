package standard

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/xiatechs/markdown-to-confluence/common"
	fh "github.com/xiatechs/markdown-to-confluence/filehandler"
	m "gitlab.com/golang-commonmark/markdown"
)

// grabtitle function collects the title of a markdown file
// and returns it as a string
func grabTitle(content string) string {
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

func convertMarkdownToHTML(content []byte, parentFolderMetaData map[string]interface{}) (*fh.FileContents, error) {
	r := bytes.NewReader(content)

	f := fh.NewFileContents()

	fmc, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		log.Println("issue parsing frontmatter - (using # title instead): %w", err)

		f.MetaData["title"] = grabTitle(string(content))
	} else {
		if len(fmc.FrontMatter) != 0 {
			f.MetaData = fmc.FrontMatter
		} else {
			f.MetaData["title"] = grabTitle(string(content))
		}
	}

	value, ok := f.MetaData["title"]
	if !ok {
		return nil, fmt.Errorf("markdown page parsing error - page title is not assigned via toml or # section")
	}

	if value == "" {
		return nil, fmt.Errorf("markdown page parsing error - page title is empty")
	}

	md := m.New(
		m.HTML(true),
		m.Tables(true),
		m.Linkify(true),
		m.Typographer(false),
		m.XHTMLOutput(true),
	)

	preformatted := md.RenderToString(content)

	var rootID int

	if _, ok := parentFolderMetaData["id"]; ok {
		rootID = parentFolderMetaData["id"].(int)
	}

	f.Body = stripFrontmatterReplaceURL(rootID, preformatted)

	return f, nil
}

// for absolute links we don't need to do any fancy pants processing
func linkFilterLogic(item string) bool {
	if strings.Contains(item, "https://") {
		return true
	}

	if strings.Contains(item, "http://") {
		return true
	}

	if strings.Contains(item, "www") {
		return true
	}

	if strings.Contains(item, "mailto:") {
		return true
	}

	return false
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

		// temporary solution to local url path issue - try and identify them with fuzzy logic
		if strings.Contains(lines[index], "<a href=") && !linkFilterLogic(lines[index]) {
			lines[index] = "[please use absolute links]"
		}

		// can't use local links yet for images
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

	var c string

	if len(sliceOne) > 1 {
		sliceTwo := strings.Split(sliceOne[1], `"`)

		if len(sliceTwo) > 1 {
			attachmentFileName := sliceTwo[0]

			rootPageID := strconv.Itoa(rootID)

			a := `<p><span class="confluence-embedded-file-wrapped">`

			b := `<img src="` + common.ConfluenceBaseURL + `/wiki/download/attachments/`

			c = rootPageID + `/` + attachmentFileName + `"></img>`
			
			d := `</span></p>`

			return a + b + c + d
		}
	}

	return item
}
