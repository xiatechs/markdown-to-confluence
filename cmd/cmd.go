// Package cmd contains code necessary to start the github action
// taking in arguments & setting them in variables in common package
package cmd

import (
	"log"
	"os"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/node"
)

// setArgs function takes in cmd line arguments
// and sets common variables (api key / space / username / project path / confluenceURL)
func setArgs() bool {
	var argLength = 6

	if len(os.Args) < argLength-1 {
		log.Println("usage: app [key space username repopath confluenceURL(optional)]")
		return false
	}

	vars := os.Args[1:]

	if len(vars) == argLength-1 {
		common.ConfluenceAPIKey = vars[0]
		common.ConfluenceSpace = vars[1]
		common.ConfluenceUsername = vars[2]
		common.ProjectPathEnv = vars[3]
		if vars[4] != "" {
			common.ConfluenceBaseURL = vars[4]
		}
	}

	return true
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
			return
		}

		node.SetAPIClient(client)

		if root.Start(common.ProjectPathEnv) {
			root.Delete()
		}
	}
}
