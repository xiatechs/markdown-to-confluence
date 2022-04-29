package standard

import (
	"bytes"
	"fmt"
	"log"
	"sort"
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
