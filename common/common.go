// Package common is for storing common constants/vars used in app
package common

import "os"

var (
	GitHubPrefix      = "/github/workspace/"
	ConfluenceSpace   = os.Getenv("CONFLUENCE_SPACE")
	ConfluenceBaseURL = os.Getenv("CONFLUENCE_BASE_URL")
)
