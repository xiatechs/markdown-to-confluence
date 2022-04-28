package cmd

import (
	"os"

	"github.com/xiatechs/markdown-to-confluence/control"
)

func Start() {
	c := control.New("test", "standard")

	c.Start(os.Getenv("PROJECT_PATH"))
}
