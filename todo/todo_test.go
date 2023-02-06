package todo

//notodo: this is the TODO PACKAGE so we don't need to pick this up

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestGenerateTODO(t *testing.T) {
	expect := func() *markdown.FileContents {
		return &markdown.FileContents{
			MetaData: map[string]interface{}{
				"title": "More info on 'testname' repo",
			},
			Body: []byte("<h2></h2>\n"),
		}
	}()

	filecontents := GenerateTODO("testname", "")

	assert.Equal(t, expect, filecontents)
}

func TestGrabTODO(t *testing.T) {
	expect := "## Filename: name-of-file\n\nRow: <1> TODO: this is some todo text\n\n\n"

	filecontents := grabTODO("TODO: this is some todo text", "name-of-file")

	assert.Equal(t, expect, filecontents)
}
