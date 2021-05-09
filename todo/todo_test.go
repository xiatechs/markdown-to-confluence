package todo

//nolint:all // exclude test from lint due to long lines
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/markdown"
)

func TestGenerateTODO(t *testing.T) {
	expect := &markdown.FileContents{
		MetaData: map[string]interface{}{
			"title": "TODO list for 'the-name-of-the-repo' repo",
		},
		Body: []byte{},
	}

	filecontents := GenerateTODO("the-name-of-the-repo")

	assert.Equal(t, expect, filecontents)
}

func TestGrabTODO(t *testing.T) {
	expect := "## Filename: name-of-file\n\nRow: <1> TODO: this is some todo text\n\n\n"

	filecontents := grabTODO("TODO: this is some todo text", "name-of-file")

	assert.Equal(t, expect, filecontents)
}
