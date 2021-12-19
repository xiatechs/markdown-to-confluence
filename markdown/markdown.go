// Package markdown provides a method for working with and parsing markdown documents
//nolint: wsl // is fine
package markdown

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/xiatechs/markdown-to-confluence/common"
	m "gitlab.com/golang-commonmark/markdown"
)

// GrabAuthors - do we want to collect authors?
var GrabAuthors bool

var sem = make(chan bool, 1)

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

// Paragraphify takes in a .puml file contents and returns
// a formatted HTML page as a string
func Paragraphify(content string) string {
	var pre string
	pre += "### To view this try copy&paste to this site: [PlainText UML Editor](https://www.planttext.com/) \n"
	pre += "### Alternatively please install a _PlantUML Visualizer plugin_ for Chrome or Firefox \n"
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

//nolint: gosec // is fine
func capGit(path string) string {
	sem <- true // race block
	log.Println("collecting authorship for ", path)
	git := exec.Command("git", "log", `--format='%ae'`, path)

	out, err := git.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	authors := make(map[string]int)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		l := strings.ReplaceAll(line, `'`, "")
		authors[l]++
	}

	output := "```" + "\n"

	index := 0
	for author, number := range authors {
		if author == "" {
			continue
		}

		no := strconv.Itoa(number)

		if index == 0 {
			output += author + " total commits: " + no
		} else {
			output += "\n" + author + " total commits: " + no
		}
	}

	output += "```"

	<-sem
	return output
}

// ParseMarkdown function uses external parsing library to grab markdown contents
// and return a filecontents object
func ParseMarkdown(rootID int, content []byte, isIndex bool, id int,
	pages map[string]string, path string) (*FileContents, error) {
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
	f.Body = stripFrontmatterReplaceURL(rootID, preformatted, isIndex, id, pages)

	if GrabAuthors {
		f.Body = append(f.Body, []byte(capGit(path))...)
	}

	return f, nil
}

func linkFilterLogic(item string) bool {
	if strings.Contains(item, "https://") {
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
func stripFrontmatterReplaceURL(rootID int, content string,
	isIndex bool, id int, pages map[string]string) []byte {
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
			lines[index] = fuzzyLogicURLdetector(lines[index], pages)
		}

		if strings.Contains(lines[index], "<img src=") {
			lines[index] = URLConverter(rootID, lines[index], isIndex, id)
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

// fuzzy logic for local links - it'll try and match the link to a generated confluence page
// if this fails, it will just return a template
func fuzzyLogicURLdetector(item string, page map[string]string) string {
	const fail = `<p>[please start your links with https://]</p>`

	urlLink := strings.Split(item, `</a>`)

	originalURLslice := strings.Split(strings.ReplaceAll(urlLink[0], "<p>", ""), `>`)
	if len(originalURLslice) <= 1 {
		return fail
	}

	originalURL := originalURLslice[1]

	sliceOne := strings.Split(item, `<p><a href="`)
	if len(sliceOne) <= 1 {
		return fail
	}

	url := strings.Split(sliceOne[1], `"`)[0]
	minimum := 0
	simMinimum := 0
	likelypage := ""
	likelyURL := ""
	first := true
	for localURL, confluencepage := range page {
		similarity := exists(localURL, url)

		if similarity != 0 && similarity > simMinimum {
			simMinimum = similarity

			// if a page link is more similar than previous page link, let's use that page
			check := levenshtein([]rune(url), []rune(localURL))

			if first {
				first = false
				minimum = check
				likelypage = confluencepage
				likelyURL = localURL
			}

			if check < minimum {
				minimum = check
				likelypage = confluencepage
				likelyURL = localURL
			}
		}
	}

	if likelypage == "" {
		return fail
	}

	log.Println("relative link -> ", url, "is LIKELY to be this page:", likelyURL, likelypage)

	a := `<p><a class="confluence-link" href="`

	b := "/wiki/spaces/" + common.ConfluenceSpace + "/pages/" + likelypage + `"`

	c := ` data-linked-resource-id="` + likelypage + `" data-base-url="` + common.ConfluenceBaseURL + `/wiki">`

	d := originalURL

	e := `</a></p>`

	return a + b + c + d + e
}

// check how many fields exist in two strings (split by '/')
func exists(a, b string) int {
	similarity := 0
	aa := strings.Split(a, "/")
	bb := strings.Split(b, "/")
	for _, line := range aa {
		for _, line2 := range bb {
			if line == line2 {
				similarity++
			}
		}
	}

	return similarity
}

// levenshtein fuzzy logic algorithm to determine similarity of two strings
func levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

// URLConverter function for images to be loaded in to confluence page
// (they must be in same directory as markdown to work)
// this function replaces local url paths in html img links
// with a confluence path for folder page attachments on parent page
func URLConverter(rootID int, item string, isindex bool, id int) string {
	sliceOne := strings.Split(item, `<p><img src="`)
	var c string
	if len(sliceOne) > 1 {
		sliceTwo := strings.Split(sliceOne[1], `"`)

		if len(sliceTwo) > 1 {
			attachmentFileName := sliceTwo[0]
			rootPageID := strconv.Itoa(rootID)
			ID := strconv.Itoa(id)
			a := `<p><span class="confluence-embedded-file-wrapped">`
			b := `<img src="` + common.ConfluenceBaseURL + `/wiki/download/attachments/`
			if isindex {
				c = ID + `/` + attachmentFileName + `"></img>`
			} else {
				c = rootPageID + `/` + attachmentFileName + `"></img>`
			}
			d := `</span></p>`

			return a + b + c + d
		}
	}

	return item
}
