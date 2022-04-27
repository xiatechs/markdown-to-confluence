package foldercrawler

import (
	"strings"

	apitest "github.com/xiatechs/markdown-to-confluence/apihandler/test"
	"github.com/xiatechs/markdown-to-confluence/filehandler/standard"
)

func New(api, fileConverter string) *Controller {
	c := &Controller{}
	switch strings.ToLower(api) {
	case "test":
		c.API = &apitest.Local{}
	}

	switch strings.ToLower(fileConverter) {
	case "standard":
		c.FH = &standard.Basic{}
	}

	return c
}
