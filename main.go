package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/yuin/goldmark"
)

const projectPathEnv = "PROJECT_PATH"
const confluenceAPIKeyEnv = "INPUT_CONFLUENCE_API_KEY"
const confluenceSpaceEnv = "INPUT_CONFLUENCE_SPACE"

func main() {
	projectPath, exists := os.LookupEnv(projectPathEnv)
	if !exists {
		log.Printf("Environment variable not set for %s, defaulting to `./`", projectPathEnv)
		projectPath = "./"
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

	err := filepath.WalkDir(projectPath, func(path string, info os.DirEntry, err error) error {
		if strings.Contains(path, "vendor") || strings.Contains(path, ".github") {
			return filepath.SkipDir
		}

		if strings.HasSuffix(info.Name(), ".md") {
			if err := parseContent(path); err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		log.Fatal(err)
	}
}

func parseContent(filename string) error {
	r, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}

	cfm, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		return err
	}

	log.Println(filename)
	log.Println(cfm.FrontMatter)

	if len(cfm.FrontMatter) == 0 {
		return fmt.Errorf("no frontmatter")
	}

	var buf bytes.Buffer

	if err := goldmark.Convert(cfm.Content, &buf); err != nil {
		return err
	}

	log.Println(buf.String())

	return nil
}
