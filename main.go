package main

import (
	"log"
	"os"

	"github.com/xiatechs/markdown-to-confluence/node"
)

const projectPathEnv = "PROJECT_PATH"

func main() {
	root := node.Node{}

	projectPath, exists := os.LookupEnv(projectPathEnv)
	if !exists {
		log.Printf("Environment variable not set for %s, defaulting to `./`", projectPathEnv)

		projectPath = "/mnt/c/Users/conotox/Desktop/Go/markdown-to-confluence"
	}

	if root.Instantiate(projectPath) { // if project path is a folder

		root.Scrub() // delete pages on confluence that shouldn't exist anymore
	}
}
