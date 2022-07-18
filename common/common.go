// Package common is for storing common constants/vars used in app
package common

import "os"

var (
	// GitHubPrefix is prefix for where files are in github
	GitHubPrefix      = "/github/workspace/"

	// ConfluenceSpace is the Space in confluence we want to work with
	ConfluenceSpace   = os.Getenv("CONFLUENCE_SPACE")

	// ConfluenceBaseURL is the URL of confluence we want to work with
	ConfluenceBaseURL = os.Getenv("CONFLUENCE_BASE_URL")
)

// Refresh - in case we want to update the vars above
func Refresh() {
	ConfluenceSpace = os.Getenv("CONFLUENCE_SPACE")
	ConfluenceBaseURL = os.Getenv("CONFLUENCE_BASE_URL")
}
