package todo

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
