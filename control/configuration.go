package control

import (
	"strings"

	"github.com/xiatechs/markdown-to-confluence/apihandler"
	"github.com/xiatechs/markdown-to-confluence/apihandler/confluence"
	"github.com/xiatechs/markdown-to-confluence/apihandler/template"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
	"github.com/xiatechs/markdown-to-confluence/filehandler/standard"
)

func New(api, fileConverter string) *Controller {
	c := &Controller{}

	/////////////////////// collect the API Handler ///////////////
	switch strings.ToLower(api) {
	case "template":
		c.API = &template.Example{}

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

func NewDI(fh filehandler.FileHandler, api apihandler.ApiController) *Controller {
	c := &Controller{
		FH:  fh,
		API: api,
	}

	return c
}
