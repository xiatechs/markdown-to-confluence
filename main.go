package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/markdown"
)

const projectPathEnv = "PROJECT_PATH"
const confluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
const confluenceAPIKeyEnv = "INPUT_CONFLUENCE_API_KEY"
const confluenceSpaceEnv = "INPUT_CONFLUENCE_SPACE"

func main() {
	projectPath, exists := os.LookupEnv(projectPathEnv)
	if !exists {
		log.Printf("Environment variable not set for %s, defaulting to `./`", projectPathEnv)

		projectPath = "./"
	}

	checkConfluenceEnv()

	err := filepath.WalkDir(projectPath, func(path string, info os.DirEntry, err error) error {
		if strings.Contains(path, "vendor") || strings.Contains(path, ".github") {
			return filepath.SkipDir
		}

		if strings.HasSuffix(info.Name(), ".md") {
			if err := processFile(path); err != nil {
				log.Println(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func processFile(path string) error {
	log.Println("Processing:", filepath.Clean(path))

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Printf("error opening file (%s): %v", path, err)
		return err
	}

	contents, err := markdown.ParseMarkdown(f)
	if err != nil {
		return err
	}

	log.Printf("%+v", contents)

	return nil
}

func checkConfluenceEnv() {
	username, exists := os.LookupEnv(confluenceUsernameEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceUsernameEnv)
	} else {
		log.Printf("API KEY: %s", username)
	}

	apiKey, exists := os.LookupEnv(confluenceAPIKeyEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceAPIKeyEnv)
	} else {
		log.Printf("API KEY: %s", apiKey)
	}

	space, exists := os.LookupEnv(confluenceSpaceEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceSpaceEnv)
	} else {
		log.Printf("SPACE: %s", space)
	}
}
