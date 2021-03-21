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

func main() {
	projectPath, exists := os.LookupEnv(projectPathEnv)
	if !exists {
		log.Printf("Environment variable not set for %s, defaulting to `./`", projectPathEnv)

		projectPath = "./"
	}

	err := filepath.WalkDir(projectPath, func(path string, info os.DirEntry, err error) error {
		if strings.Contains(path, "vendor") || strings.Contains(path, ".github") {
			return filepath.SkipDir
		}

		if strings.HasSuffix(info.Name(), ".md") {
			pageContent, err := processFile(path)
			if err != nil {
				log.Printf("error processing file: %s", err)
			}

			err = checkConfluencePages(pageContent)
			if err != nil {
				log.Printf("error completing confluence operations: %s", err)
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
func processFile(path string) (*markdown.FileContents, error) {
	log.Println("Processing:", filepath.Clean(path))

	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		log.Printf("error opening file (%s): %v", path, err)
		return nil, err
	}

	contents, err := markdown.ParseMarkdown(f)
	if err != nil {
		return nil, err
	}

	log.Printf("%T, %+v", contents, contents)

	return contents, nil
}

// checkConfluencePages runs through the CRUD operations for confluence
func checkConfluencePages(newPageContents *markdown.FileContents) error {
	fmt.Println("running find page function: ") // todo remove

	Client, err := confluence.CreateAPIClient()
	if err != nil {
		log.Printf("error creating APIClient: %s", err)
		return nil
	}

	pageTitle := newPageContents.MetaData["title"].(string)

	id, version, exists, err := Client.FindPage(pageTitle)
	if err != nil {
		return err
	}

	if !exists {

		err = Client.CreatePage(newPageContents)
		if err != nil {
			return err
		}
	} else {
		// do some check and update if required
		err = Client.UpdatePage(id, version, newPageContents)
		if err != nil {
			return err
		}
	}

	return nil
}
