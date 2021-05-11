package main

import (
	"log"
	"os"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/node"
)

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

func main() {
	if setArgs() {
		root := node.Node{}

		if client, err := confluence.CreateAPIClient(); err != nil {
			log.Println(err)
		} else if root.Start(common.ProjectPathEnv, client) {
			root.Delete()
		}
	}
}
