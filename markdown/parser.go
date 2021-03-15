// Package markdown provides a method for working with and parsing markdown documents
package markdown

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// FileContents contains information from a file after being parsed from markdown.
// `Metadata` in the format of a `map[string]interface{}` this can contain title, description, slug etc.
// `Body` a `[]byte` that contains the resulting HTML after parsing the markdown and converting to HTML using Goldmark.
type FileContents struct {
	MetaData map[string]interface{}
	Body     []byte
}

// ParseMarkdown reads the incoming data and splits out it's front matter into a
// metadata map and a converts its body into HTML
// The resulting information is put into a FileContents for use.
func ParseMarkdown(r io.Reader) (*FileContents, error) {
	// fmc is a shorthand for frontmatter and content, the 2 sections of a doc page.
	// we subsequently put this into a FileContents type.
	fmc, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing file contents and front matter: %w", err)
	}

	if len(fmc.FrontMatter) == 0 {
		return nil, fmt.Errorf("no frontmatter")
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)

	var buf bytes.Buffer

	if err := md.Convert(fmc.Content, &buf); err != nil {
		return nil, fmt.Errorf("error converting markdown to HTML: %w", err)
	}

	return &FileContents{
		MetaData: fmc.FrontMatter,
		Body:     buf.Bytes(),
	}, nil
}
