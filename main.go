package main

import (
	"log"
	"os"

	"github.com/xiatechs/markdown-to-confluence/node"
)

const projectPathEnv = "PROJECT_PATH" // need to change this

func main() {
	root := node.Node{}

	projectPath, exists := os.LookupEnv(projectPathEnv)
	if !exists {
		log.Printf("Environment variable not set for %s, defaulting to `./`", projectPathEnv)

		projectPath = "../testfolder"
	}

	if root.Instantiate(projectPath) {
		node.PrintOverview()
	}
}
