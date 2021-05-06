package main

import (
	"log"
	"os"

	"github.com/xiatechs/markdown-to-confluence/confluence"
	"github.com/xiatechs/markdown-to-confluence/node"
)

const projectPathEnv = "PROJECT_PATH"
const defaultProjectPath = "./"

func main() {
	root := node.Node{}

	if client, err := confluence.CreateAPIClient(); err != nil {
		log.Println(err)
	} else {
		projectPath, exists := os.LookupEnv(projectPathEnv)
		if !exists {
			log.Printf("Environment variable not set for %s", projectPathEnv)

			projectPath = defaultProjectPath
		}

		if root.Instantiate(projectPath, client) { // if project path is a folder
			root.Scrub() // delete pages on confluence that shouldn't exist anymore
		}
	}
}
