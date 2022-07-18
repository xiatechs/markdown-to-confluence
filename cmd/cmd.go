// Package cmd is for running the entire service
package cmd

import (
	"os"

	"github.com/xiatechs/markdown-to-confluence/control"
)

// Start - start the service
func Start() {
	c := control.New("test", "standard")

	c.Start(os.Getenv("PROJECT_PATH"))
}
