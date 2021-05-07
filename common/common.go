// Package common is for storing common constants/vars used in app
package common

var ConstantsHardCoded = false

const (
  // ConfluenceBaseURL is the base URL for the confluence page you want to work with
  	ConfluenceBaseURL     = "https://xiatech-markup.atlassian.net"
	ConfluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
	ConfluenceAPIKeyEnv   = "INPUT_CONFLUENCE_API_KEY"
	ConfluenceSpaceEnv    = "INPUT_CONFLUENCE_SPACE"
	EnvsNotSetError       = "environment variable not set, please assign values for: "
)
