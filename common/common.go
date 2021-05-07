// Package common is for storing common constants/vars used in app
package common

var ConstantsHardCoded = false

const (
  // ConfluenceBaseURL is the base URL for the confluence page you want to work with
  ConfluenceBaseURL = "https://xiatech-markup.atlassian.net"
	confluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
	confluenceAPIKeyEnv   = "INPUT_CONFLUENCE_API_KEY"
	confluenceSpaceEnv    = "INPUT_CONFLUENCE_SPACE"
	CnvsNotSetError       = "environment variable not set, please assign values for: "
)
