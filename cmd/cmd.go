// Package cmd contains code necessary to start the github action
// taking in arguments & setting them in variables in common package
package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/node"
)

// setArgs function takes in cmd line arguments
// and sets common variables (api key / space / username / project path)
func setArgs() bool {
	var argLength = 4

	if len(os.Args) > 1 {
		vars := strings.Split(os.Args[1], "_")

		if len(vars) == argLength {
			common.ConfluenceAPIKey = vars[0]
			common.ConfluenceSpace = vars[1]
			common.ConfluenceUsername = vars[2]
			common.ProjectPathEnv = vars[3]

			return true
		}
	}

	log.Println("usage: app apikey_space_username_path")

	return false
}

// Start function sets argument inputs, creates confluence API client
// and begins the process of creating confluence pages via calling
// the node.Start method
// if node.Start returns true, then calls node.Delete method
func Start() {
	if setArgs() {
		root := node.Node{}

		if client, err := confluence.CreateAPIClient(); err != nil {
			log.Println(err)
		} else if root.Start(common.ProjectPathEnv, client) {
			root.Delete()
		}
	}
}
