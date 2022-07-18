// Package control is for handling the service
//nolint: gocritic // is fine
package control

import (
	"strings"

	"github.com/xiatechs/markdown-to-confluence/apihandler"
	"github.com/xiatechs/markdown-to-confluence/apihandler/confluence"
	"github.com/xiatechs/markdown-to-confluence/apihandler/template"
	"github.com/xiatechs/markdown-to-confluence/filehandler"
	"github.com/xiatechs/markdown-to-confluence/filehandler/standard"
)

// New - create a controller using strategy architecture
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

// NewDI - create a controller using dependency injection
func NewDI(fh filehandler.FileHandler, api apihandler.APIController) *Controller {
	c := &Controller{
		FH:  fh,
		API: api,
	}

	return c
}
