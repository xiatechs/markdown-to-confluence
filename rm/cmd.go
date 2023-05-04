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
	apiSpace := flag.String("space", "", "the confluence API Space")
	userName := flag.String("username", "", "the confluence API username")
	folderPath := flag.String("folderpath", "", "the source of the documentation")
	masterPageId := flag.String("id", "0", "the id of the master page - default is 0 (root)")
	url := flag.String("url", "https://xiatech.atlassian.net", "the url for confluence")
	onlyDocs := flag.Bool("docs", true, "parse only the /docs folder")
	flag.Parse()
	if onlyDocs != nil {
		common.OnlyDocs = *onlyDocs
	} else {
		log.Println("docs flag is missing")
	}

	if url != nil {
		common.ConfluenceBaseURL = *url
	} else {
		log.Println("url flag is missing")
	}

	if masterPageId != nil {
		common.ProjectMasterID, err = strconv.Atoi(*masterPageId)
		if err != nil {
			log.Println("masterpageID should be an int. If MTC is to be the root enter 0")
			return false
		}
	} else {
		log.Println("id flag is missing")
	}

	if folderPath != nil {
		common.ProjectPathEnv = *folderPath
		common.ProjectPathEnv = strings.ReplaceAll(common.ProjectPathEnv, " ", "-") // replace spaces with -
	} else {
		log.Println("folderpath flag is missing")
	}

	if userName != nil {
		common.ConfluenceUsername = *userName
	} else {
		log.Println("username flag is missing")
	}

	if apiSpace != nil {
		common.ConfluenceSpace = *apiSpace
	} else {
		log.Println("space flag is missing")
	}

	if apiKey != nil {
		common.ConfluenceAPIKey = *apiKey
	} else {
		log.Println("key flag is missing")
	}

	// if lengths are 0 then it's being passed in wrong - good for debug
	log.Println("secrets - lengths")
	log.Println("username", len(common.ConfluenceUsername))
	log.Println("key", len(common.ConfluenceAPIKey))
	log.Println("space", len(common.ConfluenceSpace))
	log.Println("contents - string")
	log.Println("folderpath", common.ProjectPathEnv)
	log.Println("url", common.ConfluenceBaseURL)
	log.Println("docs", common.OnlyDocs)
	log.Println("rootid", common.ProjectMasterID)

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
