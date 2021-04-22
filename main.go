package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/markdown"
)

const (
	// add your confluence username / api key / space here before building app
	confluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
	confluenceAPIKeyEnv   = "INPUT_CONFLUENCE_API_KEY"
	confluenceSpaceEnv    = "INPUT_CONFLUENCE_SPACE"
)

// grab 1 argument (filepath) when calling app
func grabargs() (valid bool, projectPath string) {
	if len(os.Args) == 2 {
		projectPath = os.Args[1]
	} else {
		log.Println("usage: app [folder/.]")
		return false, ""
	}
	return true, projectPath
}

//thefilepath is the relative filepath of the app, localpath is the folder you want to run this app through
func iterate(localpath string) {
	//Go 1.15 doesn't have the WalkDir method for filepath package so adjusted it below
	filepath.Walk(localpath, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		path := info.Name()
		if strings.Contains(path, "vendor") || strings.Contains(path, ".github") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(info.Name(), ".md") {
			if err := processFile(fpath); err != nil {
				log.Println(err)
			}
		}
		return nil
	})
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

func main() {
	if ok, projectPath := grabargs(); ok {
		checkConfluenceEnv()
		iterate(projectPath) //pass the project path
	}
}
