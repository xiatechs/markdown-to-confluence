package foldercrawler

import (
	"strings"

	"github.com/xiatechs/markdown-to-confluence/apihandler/confluence"
	apitest "github.com/xiatechs/markdown-to-confluence/apihandler/test"
	"github.com/xiatechs/markdown-to-confluence/filehandler/standard"
)

func New(api, fileConverter string) *Controller {
	c := &Controller{}

	/////////////////////// collect the API Handler ///////////////
	switch strings.ToLower(api) {
	case "test":
		c.API = &apitest.Local{}

	case "confluence":
		c.API = confluence.NewAPIClient()
	}

	/////////////////////// collect the File Handler ///////////////
	switch strings.ToLower(fileConverter) {
	case "standard":
		c.FH = &standard.Basic{}
	}

	return c
}
