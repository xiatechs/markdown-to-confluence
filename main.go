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

	err = checkConfluencePages(contents)
	if err != nil {
		return err
	}
	return nil

	return nil
}

// checkConfluencePages runs through the CRUD operations for confluence
func checkConfluencePages(newPageContents *markdown.FileContents) error {
	fmt.Println("running find page function: ") //todo remove
	a, ok := confluence.NewAPIClient()
	if !ok {
		log.Println("error creating a new client")
		return nil
	}

	pageTitle := newPageContents.MetaData["title"].(string)

	id, version, exists, err := a.FindPage(pageTitle)
	if err != nil {
		return err
	}

	if !exists {
		//todo: create page
	} else {
		//do some check and update if required
		a.UpdatePage(id, version, newPageContents)
	}

	return nil
}
