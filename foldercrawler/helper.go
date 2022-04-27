package foldercrawler

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// withinDirectory function checks to see if the file (base) is within the folder (path)
func withinDirectory(base, path string) bool {
	return strings.Count(path, "/")-strings.Count(base, "/") == 1
}

// isVendorOrGit function takes in name of folder and
// checks if it is a vendor or github folder
func isVendorOrGit(name string) bool {
	if strings.Contains(name, "vendor") || strings.Contains(name, ".github") || strings.Contains(name, ".git") {
		return true
	}

	return false
}

// isFolder function checks whether a file is a folder or not
func isFolder(name string) bool {
	file, err := os.Open(filepath.Clean(name))
	if err != nil {
		log.Println(err)
		return false
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return false
	}

	return fileInfo.IsDir()
}

func isReadMeFile(fpath string) bool {
	return strings.Contains(strings.ToLower(filepath.Base(fpath)), "readme")
}

func isMarkdownFile(fpath string) bool {
	return strings.HasSuffix(strings.ToLower(fpath), ".md")
}
