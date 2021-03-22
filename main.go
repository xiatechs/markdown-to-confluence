package main

import (
	"fmt"
	"github.com/xiatechs/markdown-to-confluence/confluence"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
			if err := processFile(path); err != nil {
				log.Printf("error processing file: %s", err)
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

	err = checkConfluencePages(contents)
	if err != nil {
		log.Printf("error completing confluence operations: %s", err)
	}

	return nil
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

	pageResult, err := Client.FindPage(pageTitle)
	if err != nil {
		return err
	}

	if pageResult == nil {
		fmt.Println("page does not exists, creating it now...")

		err = Client.CreatePage(newPageContents)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("page exists, updating confluence now...")

		pageID, err := strconv.Atoi(pageResult.Results[0].ID)
		if err != nil {
			return err
		}

		err = Client.UpdatePage(pageID, int64(pageResult.Results[0].Version.Number), newPageContents)
		if err != nil {
			return err
		}
	}

	return nil
}
