// Package markdown provides a method for working with and parsing markdown documents
//nolint: wsl // is fine
package markdown

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/xiatechs/markdown-to-confluence/common"
	m "gitlab.com/golang-commonmark/markdown"
)

// GrabAuthors - do we want to collect authors?
var (
	GrabAuthors bool
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

type author struct {
	name    string
	howmany int
}

type authors []author // not using a map so the order of authors can be maintained

func (a *authors) append(item string) {
	au := *a
	var exists bool
	for index := range au {
		if au[index].name == item {
			au[index].howmany++
			exists = true
			break
		}
	}

	if !exists {
		au = append(au, author{
			name:    item,
			howmany: 1,
		})
	}

	*a = au
}

func (a *authors) sort() {
	au := *a

	sort.Slice(au, func(i, j int) bool {
		return au[i].howmany > au[j].howmany
	})

	*a = au
}

//nolint: gosec // is fine
// use git to capture authors by username & email & commits
func capGit(path string) string {
	here, _ := os.Getwd()
	log.Println("collecting authorship for ", path)
	git := exec.Command("git", "log", "--all", `--format='%an | %ae'`, "--", here)

	out, err := git.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	a := authors{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		l := strings.ReplaceAll(line, `'`, "")
		a.append(l)
	}

	a.sort()

	// to let the output be displayed in confluence - wrapping it in code block
	output := `<pre><code>
[authors | email addresses | how many commits]
`

	index := 0
	for _, author := range a {
		if author.name == "" {
			continue
		}

		no := strconv.Itoa(author.howmany)

		if index == 0 {
			output += author.name + " - total commits: " + no
		} else {
			output += `
` + author.name + " - total commits: " + no
		}

		index++
	}

	output += `
</code></pre>`

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

		// can't use local links yet for images
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

type fpage struct {
	distance       int
	sim            int
	confluencepage string
	url            string
}

type pages []fpage

func (p pages) filter() fpage {
	max := p[0].sim
	i := 0
	for index := range p {
		if p[index].distance == p[0].distance {
			if p[index].sim > max {
				i = index
				max = p[index].sim
			}
		}
	}

	return p[i]
}

// fuzzy logic for local links - it'll try and match the link to a generated confluence page
// if this fails, it will just return a template
func fuzzyLogicURLdetector(item string, page map[string]string) string {
	const fail = `<p>[please start your links with https://]</p>`

	urlLink := strings.Split(item, `</a>`)

	originalURLslice := strings.Split(strings.ReplaceAll(urlLink[0], "<p>", ""), `>`)
	if len(originalURLslice) <= 1 {
		log.Println(item, "failed to be processed during fuzzy logic")
		return fail
	}

	originalURL := originalURLslice[1]

	sliceOne := strings.Split(item, `<a href="`)
	if len(sliceOne) <= 1 {
		log.Println(item, "failed to be processed during fuzzy logic")
		return fail
	}

	url := strings.Split(sliceOne[1], `"`)[0] // the local/relative URL for the page
	simMinimum := 0

	p := pages{}
	for localURL, confluencepage := range page {
		url := strings.ReplaceAll(url, "%20", " ")
		similarity := exists(localURL, url) // how many similar fields are in the two links we are looking at

		if similarity != 0 && similarity >= simMinimum {
			simMinimum = similarity

			// if a page link is more similar than previous page link, let's use that page
			// and use levenshtein algorithm to determine which is most similar out of a group
			check := levenshtein([]rune(url), []rune(localURL))

			p = append(p, fpage{
				distance:       check,
				confluencepage: confluencepage,
				url:            localURL,
				sim:            similarity,
			})
		}
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].distance < p[j].distance
	})

	if len(p) == 0 {
		return fail
	}

	thepage := p.filter()

	likelyURL := thepage.url
	likelypage := thepage.confluencepage

	log.Println("relative link -> ", url, "is LIKELY to be this page:", likelyURL, likelypage)

	// to format this in confluence we must follow how confluence formats its content in the web frontend
	a := `<p><a class="confluence-link" href="`

	b := "/wiki/spaces/" + common.ConfluenceSpace + "/pages/" + likelypage + `"`

	c := ` data-linked-resource-id="` + likelypage + `" data-base-url="` + common.ConfluenceBaseURL + `/wiki">`

	d := originalURL

	e := `</a></p>`

	return a + b + c + d + e
}

type fielditem struct {
	item   string
	index1 int
	index2 int
}

type fielditems []fielditem

func (f fielditems) validate() bool {
	if len(f) == 1 {
		return true
	}

	for index := len(f) - 1; index > 0; index-- {
		if f[index].index1-f[index-1].index1 != 1 {
			return false
		}

		if f[index].index2-f[index-1].index2 != 1 {
			return false
		}
	}

	return true
}

// check how many fields exist in two strings (split by '/')
func exists(a, b string) int {
	var f fielditems
	similarity := 0
	a = strings.ReplaceAll(a, "../", "")
	b = strings.ReplaceAll(b, "../", "")
	a = strings.ReplaceAll(a, ".", "")
	b = strings.ReplaceAll(b, ".", "")
	aa := strings.Split(a, "/")
	bb := strings.Split(b, "/")
	for index, line := range aa {
		for index2, line2 := range bb {
			if line == line2 {
				similarity++
				f = append(f, fielditem{
					item:   line,
					index1: index,
					index2: index2,
				})
			}
		}
	}

	sort.Slice(f, func(i, j int) bool {
		return f[i].index1 < f[j].index1
	})

	/*
		INVALID
		a) ../node/testfolder/folder
		b ../node/folder/testfolder

		VALID
		a) ../../node/testfolder/folder
		b) node/testfolder/folder
	*/

	if f.validate() {
		return similarity
	}

	return 0
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
