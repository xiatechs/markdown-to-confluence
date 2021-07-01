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
// and sets common variables (api key / space / username / project path / confluenceURL)
func setArgs() bool {
	var argLength = 5

	if len(os.Args) > 1 {
		vars := strings.Split(os.Args[1], "_")

		if len(vars) == argLength-1 {
			common.ConfluenceAPIKey = vars[0]
			common.ConfluenceSpace = vars[1]
			common.ConfluenceUsername = vars[2]
			common.ProjectPathEnv = vars[3]

			return true
		}

		if len(vars) == argLength {
			common.ConfluenceAPIKey = vars[0]
			common.ConfluenceSpace = vars[1]
			common.ConfluenceUsername = vars[2]
			common.ProjectPathEnv = vars[3]
			common.ConfluenceBaseURL = vars[4]

			return true
		}
	}

	log.Println("usage: app apikey_space_username_path_confluenceURL")

	return false
}

// Start function sets argument inputs, creates confluence API client
// and begins the process of creating confluence pages via calling
// the node.Start method
// if node.Start returns true, then calls node.Delete method
func Start() {
	if setArgs() {
		root := node.Node{}

		client, err := confluence.CreateAPIClient()
		if err != nil {
			log.Println(err)
		}

		node.NodeAPIClient = client

		if root.Start(common.ProjectPathEnv) {
			root.Delete()
		}
	}
}
