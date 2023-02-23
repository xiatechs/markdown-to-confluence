// Package markdown provides a method for working with and parsing markdown documents
//
//nolint:wsl // is fine
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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// grabtitle function collects the filename of a markdown file
// and returns it as a string
//
//nolint:deadcode,unused // not used anymore
func grabtitle(path string) string {
	return strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
}

// newFileContents function creates a new filecontents object
func newFileContents() *FileContents {
	f := FileContents{}
	f.MetaData = make(map[string]interface{})

	return &f
}

// Paragraphify takes in a .puml file contents and returns
// a formatted HTML page as a string
// currently unused
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

// use git to capture authors by username & email & commits
//
//nolint:gosec // is fine
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
func ParseMarkdown(rootID int, content []byte, isIndex bool,
	pages map[string]string, path, abs, fileName string) (*FileContents, error) {
	r := bytes.NewReader(content)
	f := newFileContents()
	fmc, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		log.Println("issue parsing frontmatter - (using # title instead): %w", err)
	} else if len(fmc.FrontMatter) != 0 {
		f.MetaData = fmc.FrontMatter
	}

	// if the file name is readme.md then then the space should be named after the final folder
	if strings.ToLower(fileName) == "readme.md" {
		fileName = strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	}

	f.MetaData["title"] = fileName

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
	f.Body = stripFrontmatterReplaceURL(preformatted, isIndex, pages, abs, fileName)

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
func stripFrontmatterReplaceURL(content string,
	isIndex bool, pages map[string]string, abs, fileName string) []byte {
	var pre string

	var frontmatter bool

	lines := strings.Split(content, "\n")

	for index := range lines {
		if strings.Contains(lines[index], "+++") {
			frontmatter = flip(frontmatter)
			continue
		}

		// header lines are converted to ProperCase for confluence local linking
		htmlHeaders := []string{"<h1>", "<h2>", "<h3>", "<h4>", "<h5>", "<h6>"}

		for _, header := range htmlHeaders {
			if strings.Contains(lines[index], header) {
				lines[index] = updateHeaderToProperCase(lines[index])
			}
		}

		// correct the local url paths to be absolute paths
		if strings.Contains(lines[index], "<a href=") && !linkFilterLogic(lines[index]) {
			lines[index] = relativeURLdetector(lines[index], pages, abs, fileName)
		}

		// set up the local url image links
		if strings.Contains(lines[index], "<img src=") {
			lines[index] = URLConverter(pages, lines[index], isIndex, abs)
		}

		if !frontmatter {
			pre += lines[index] + "\n"
		}
	}

	pre = strings.TrimSpace(pre)

	return []byte(pre)
}

// updateHeaderToProperCase makes all headers be in Proper Case so local links work
func updateHeaderToProperCase(line string) string {
	splitOnLinks := strings.Split(line, `<a href="`)
	caser := cases.Title(language.English)

	if len(splitOnLinks) == 1 { // means there are no links in the header
		line = updateHeaderWithNoLinks(line, caser)
	} else { // means there are links in the header - don't want to alter the links
		line = updateHeaderWithLinks(splitOnLinks, caser)
	}

	return line
}

func updateHeaderWithNoLinks(line string, caser cases.Caser) string {
	line = caser.String(line)

	// this changes the html tags to be capitalized so reverse them back
	line = updateHTMLHeaders(line)

	// drop brackets from the title
	line = strings.ReplaceAll(line, "(", "")
	line = strings.ReplaceAll(line, ")", "")

	return line
}

func updateHeaderWithLinks(splitOnLinks []string, caser cases.Caser) string {
	line := splitOnLinks[0]
	line = caser.String(line)

	for i := 1; i < len(splitOnLinks); i++ {
		line += `<a href="` // add the '<a href="' back in

		extraParts := strings.SplitN(splitOnLinks[i], ">", 2) //nolint:gomnd // split the final part on the first >

		for ii := range extraParts {
			if ii == 0 { // add the link and the '>' back in
				line += extraParts[ii]
				line += `>`
			}

			if ii > 0 { // ProperCase the header and add it to the line
				extraParts[ii] = caser.String(extraParts[ii])
				line += extraParts[ii]
			}
		}

		// this changes the html tags to be capitalized so reverse them back
		line = updateHTMLHeaders(line)

		// drop brackets from the title
		line = strings.ReplaceAll(line, "(", "")
		line = strings.ReplaceAll(line, ")", "")
	}

	return line
}

// uncapitalize the 'h' in the html header tag
func updateHTMLHeaders(line string) string {
	line = strings.ReplaceAll(line, "H1>", "h1>")
	line = strings.ReplaceAll(line, "H2>", "h2>")
	line = strings.ReplaceAll(line, "H3>", "h3>")
	line = strings.ReplaceAll(line, "H4>", "h4>")
	line = strings.ReplaceAll(line, "H5>", "h5>")
	line = strings.ReplaceAll(line, "H6>", "h6>")

	return line
}

// flip function returns the opposite of bool
func flip(b bool) bool {
	return !b
}

//nolint:unused // not used anymore
type fpage struct {
	distance       int
	sim            int
	confluencepage string
	url            string
}

//nolint:unused // not used anymore
type pages []fpage

//nolint:unused // not used anymore
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

// takes in the absolute URL and will match the relative link to a generated confluence page
// if this fails, it will just return a template
func relativeURLdetector(item string, page map[string]string, abs, fileName string) string {
	const fail = `<p>[failed during relativeURLdetector]</p>`

	urlLink := strings.Split(item, `</a>`)

	originalURLslice := strings.Split(strings.ReplaceAll(urlLink[0], "<p>", ""), `>`)
	if len(originalURLslice) <= 1 {
		return fail
	}

	sliceOne := strings.Split(item, `<a href="`)
	if len(sliceOne) <= 1 {
		return fail
	}

	url := strings.Split(sliceOne[1], `"`)[0] // the local/relative URL for the page

	// create the absolute url
	updatedURL, localLink := convertRelativeToAbsoluteURL(abs, url)

	// confluence is case sensitive - headers are saved using proper case (i.e. So The Title Is Always Like This)
	caser := cases.Title(language.English)
	localLink = caser.String(localLink)

	var link string

	// if there were any local links (identified by #) then create the link for confluence
	if localLink != "" {
		// as there is a local link the path must be in a .md file
		// if the updated url does not contain a .md then it must be a local link for the current .md file so add it
		if !strings.Contains(updatedURL, ".md") {
			updatedURL += "/" + fileName
		}

		fileName = strings.ReplaceAll(fileName, " ", "+")
		localLinkAbs := strings.ReplaceAll(abs, "/", "+")

		link = "/" + fileName
		if !strings.HasPrefix(localLinkAbs, "+") {
			link += "+"
		}

		link += localLinkAbs + "#" + localLink
	}

	// replace the relative url in the item with the absolute url
	splitItem := strings.Split(item, "<a href=")

	stringToReturn := splitItem[0]

	for i := 1; i < len(splitItem); i++ {
		link := generateLink(page, updatedURL, link)
		extraParts := strings.SplitN(splitItem[i], ">", 2) //nolint:gomnd // only want to split the final part on the first >

		stringToReturn += link

		for ii := range extraParts {
			if ii != 0 {
				stringToReturn += extraParts[ii]
			}
		}
	}

	return stringToReturn
}

func generateLink(page map[string]string, updatedURL string, localLink string) string {
	// to format this in confluence we must follow how confluence formats its content in the web frontend
	a := `<a href="/wiki/spaces/` + common.ConfluenceSpace + `/pages/` + page[updatedURL] + localLink + `" `

	b := `title="/wiki/spaces/` + common.ConfluenceSpace + `/pages/` + page[updatedURL] + localLink + `" `

	c := `data-linked-resource-id="` + page[updatedURL] + `" `

	d := `data-linked-resource-type="page" `

	e := `data-renderer-mark="true" `

	f := `class="css-tgp101">`

	return a + b + c + d + e + f
}

//nolint:unused // not used anymore
type fielditem struct {
	item   string
	index1 int
	index2 int
}

//nolint:unused // not used anymore
type fielditems []fielditem

//nolint:unused // not used anymore
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
//
//nolint:deadcode,unused // not used anymore
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
//
//nolint:deadcode,unused // not used anymore
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

//nolint:unused // not used anymore
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
func URLConverter(page map[string]string, item string, isindex bool, abs string) string {
	sliceOne := strings.Split(item, `<img src="`)

	if len(sliceOne) > 1 {
		sliceTwo := strings.Split(sliceOne[1], `"`)

		if len(sliceTwo) > 1 {
			attachmentFileName := sliceTwo[0]

			// create the absolute url
			updatedURL, _ := convertRelativeToAbsoluteURL(abs, attachmentFileName)

			// split the path so we can rip out the file name
			splitURL := strings.Split(updatedURL, "/")
			urlWithoutFile := strings.Join(splitURL[:len(splitURL)-1], "/")

			stringToReturn := sliceOne[0]

			stringToReturn += `<span class="confluence-embedded-file-wrapped">`

			stringToReturn += `<img class="confluence-embedded-image" loading="lazy" `

			//nolint:lll /// set text
			stringToReturn += `src="` + common.ConfluenceBaseURL + `/wiki/download/attachments/` + page[urlWithoutFile] + `/` + splitURL[len(splitURL)-1] + `" `

			//nolint:lll /// set text
			stringToReturn += `data-image-src="` + common.ConfluenceBaseURL + `/wiki/download/attachments/` + page[urlWithoutFile] + `/` + splitURL[len(splitURL)-1] + `" `

			stringToReturn += `data-linked-resource-id="` + page[urlWithoutFile] + `" `

			stringToReturn += `data-linked-resource-type="attachment"></img></span>`

			stringToReturn += strings.Replace(sliceTwo[len(sliceTwo)-1], " />", "", 1)

			return stringToReturn
		}
	}

	return item
}

// convertRelativeToAbsoluteURL function takes in a relative url, and generates
// the correct absolute url based on the file path passed in
func convertRelativeToAbsoluteURL(abs, url string) (string, string) {
	var localLink string

	if strings.Contains(url, "warehouse") && !strings.Contains(url, "diagrams") {
		log.Printf("here")
	}

	// split on #
	// length 1 means no local links
	// length 2 means local links so save everything after the #
	if len(strings.Split(url, "#")) == 2 { //nolint: gomnd // magic length 2
		localLink = strings.Split(url, "#")[1]
	}

	url = strings.ReplaceAll(url, "%20", " ")
	abs = strings.ReplaceAll(abs, "%20", " ")

	splitRelativeURL := strings.Split(url, "/")
	splitAbsoluteURL := strings.Split(abs, "/")

	var firstFolder int

	for i := range splitRelativeURL {
		switch splitRelativeURL[i] {
		case ".":
			firstFolder = i + 1
		case "..": // need to remove the final folder for each '..'
			splitAbsoluteURL = splitAbsoluteURL[:len(splitAbsoluteURL)-1]

			firstFolder = i + 1
		}
	}

	// remove any trailing # links on the relative url
	splitRelativeURL[len(splitRelativeURL)-1] = strings.Split(splitRelativeURL[len(splitRelativeURL)-1], "#")[0]

	// if the last element is blank remove it so it doesn't add a trailing slash
	if splitRelativeURL[len(splitRelativeURL)-1] == "" {
		splitRelativeURL = splitRelativeURL[:len(splitRelativeURL)-1]
	}

	// append to the end of the absolute url
	return strings.Join(append(splitAbsoluteURL, splitRelativeURL[firstFolder:]...), "/"), localLink
}
