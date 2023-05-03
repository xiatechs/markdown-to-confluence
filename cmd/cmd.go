// Package cmd contains code necessary to start the github action
// taking in arguments & setting them in variables in common package
package cmd

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/common"
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/markdown"
	"github.com/xiatechs/markdown-to-confluence/node"
)

func setFlags() bool {
	var err error
	apiKey := flag.String("key", "", "the confluence API Key")
	if apiKey != nil {
		common.ConfluenceAPIKey = *apiKey
	}
	apiSpace := flag.String("space", "", "the confluence API Space")
	if apiSpace != nil {
		common.ConfluenceSpace = *apiSpace
	}
	userName := flag.String("username", "", "the confluence API username")
	if userName != nil {
		common.ConfluenceUsername = *userName
	}
	folderPath := flag.String("folderpath", "", "the source of the documentation")
	if folderPath != nil {
		common.ProjectPathEnv = *folderPath
		common.ProjectPathEnv = strings.ReplaceAll(common.ProjectPathEnv, " ", "-") // replace spaces with -
	}
	masterPageId := flag.String("master-page-id", "0", "the id of the master page - default is 0 (root)")
	if masterPageId != nil {
		common.ProjectMasterID, err = strconv.Atoi(*masterPageId)
		if err != nil {
			log.Println("masterpageID should be an int. If MTC is to be the root enter 0")
			return false
		}
	}
	url := flag.String("confluence-URL", "https://xiatech.atlassian.net", "the url for confluence")
	if url != nil {
		common.ConfluenceBaseURL = *url
	}
	onlyDocs := flag.Bool("docs", true, "parse only the /docs folder")
	if onlyDocs != nil {
		common.OnlyDocs = *onlyDocs
	}
	flag.Parse()

	return true
}

// Start function sets argument inputs, creates confluence API client
// and begins the process of creating confluence pages via calling
// the node.Start method
// if node.Start returns true, then calls node.Delete method
func Start() {
	markdown.GrabAuthors = false

	if setFlags() {
		root := node.Node{}

		client, err := confluence.CreateAPIClient()
		if err != nil {
			log.Println(err)
			return
		}

		node.SetAPIClient(client)

		if root.Start(common.ProjectMasterID, common.ProjectPathEnv, common.OnlyDocs) {
			root.Delete()
		}
	}
}
