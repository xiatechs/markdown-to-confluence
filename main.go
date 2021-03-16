package main

import (
	"fmt"
	"github.com/xiatechs/markdown-to-confluence/confluence"
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

// processFile is the function called on eligible files to handle uploads.
// API calls should be in here.
// Potentially this could hang off a struct type that contains an instance of API
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

	err = checkConfluenceFunc(contents)
	if err != nil {
		return err
	}
	return nil

	return nil
}

// checkConfluenceEnv is a placeholder function for checking the required env vars are set
func checkConfluenceEnv() {
	username, exists := os.LookupEnv(confluenceUsernameEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceUsernameEnv)
	} else {
		log.Printf("API KEY: %s", username)
	}

	space, exists := os.LookupEnv(confluenceSpaceEnv)
	if !exists {
		log.Printf("Environment variable not set for %s", confluenceSpaceEnv)
	} else {
		log.Printf("SPACE: %s", space)
	}

}

func checkConfluenceFunc(newPageContents *markdown.FileContents) error {
	// todo: search confluence for filename
	// some logic to see if content is accurate
	// update or create logic
	// push new data to page
	fmt.Println("running find page function: ")
	a, ok := confluence.NewAPIClient()
	if !ok {
		log.Println("error creating a new client")
		return nil
	}
	// return ID and also the version number bool
	id, version, exists, err := a.FindPage(newPageContents.MetaData["title"].(string))
	if err != nil {
		return err
	}

	if !exists {
		//create page
	} else {
		//do some check and update if required, possibly decode newPageContents.Body
		a.UpdatePage(id, version, newPageContents)
	}

	return nil
}
