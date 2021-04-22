package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiatechs/markdown-to-confluence/markdown"
	"github.com/xiatechs/markdown-to-confluence/object"
)

var confluenceobject = object.ConfluenceObject

// grab 1 argument (filepath) when calling app
func grabargs() (valid bool, projectPath string) {
	if len(os.Args) > 1 {
		projectPath = os.Args[1]
	} else {
		log.Println("usage: app [folder]")

		return false, ""
	}

	return true, projectPath
}

// check to see if the name of the file ends with .md i.e it's a markdown file
func checkinfoname(fpath, name string) {
	if strings.HasSuffix(name, ".md") {
		if err := processFile(fpath); err != nil {
			log.Println(err)
		}
	}
}

// iterates through files in a filepath. localpath is the folder you want to run this app through
func iterate(localpath string) {
	// Go 1.15 doesn't have the WalkDir method for filepath package so adjusted it below
	err := filepath.Walk(localpath, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if strings.Contains(info.Name(), "vendor") || strings.Contains(info.Name(), ".github") {
			return filepath.SkipDir
		}
		checkinfoname(fpath, info.Name())
		return nil
	})
	if err != nil {
		log.Println(err)
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

	return nil
}

/* Since this is a local binary that would be ran, do we need this?
// checkConfluenceEnv is a placeholder function for checking the required env vars are set
func (c confluenceVars) checkConfluenceEnv() bool {
	var somethingWrong bool
	username, exists := os.LookupEnv(c.ConfluenceUsernameEnv)
	if !exists {
		log.Println("Environment variable Username not set")
		somethingWrong = true
	} else {
		log.Printf("Username: %s", username)
	}

	apiKey, exists := os.LookupEnv(c.ConfluenceAPIKeyEnv)
	if !exists {
		log.Println("Environment variable ConfluenceAPIKeyEnv not set")
		somethingWrong = true
	} else {
		log.Printf("API KEY: %s", apiKey)
	}

	space, exists := os.LookupEnv(c.ConfluenceSpaceEnv)
	if !exists {
		log.Println("Environment variable ConfluenceSpaceEnv not set")
		somethingWrong = true
	} else {
		log.Printf("SPACE: %s", space)
	}
	if somethingWrong {
		log.Println("Please update the confobject.json located where this application is located")
		return false
	}
	return true
}*/

func main() {
	if ok := confluenceobject.Load(); ok {
		if ok, projectPath := grabargs(); ok {
			iterate(projectPath) // pass the project path
		}
	}
}
