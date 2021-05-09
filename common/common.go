// Package common is for storing common constants/vars used in app
//nolint:all // text in here could be long but on purpose i.e a long URL
package common

// ConstantsHardCoded - set to true if you are testing this in
// UAT environment & want to edit constants to be actual API parameters
var ConstantsHardCoded = false

const (
	// ConfluenceBaseURL is the base URL for the confluence page you want to work with
	ConfluenceBaseURL = "https://xiatech-markup.atlassian.net"

	// ConfluenceUsernameEnv is to collect external env var INPUT_CONFLUENCE_USERNAME
	ConfluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"

	// ConfluenceAPIKeyEnv is to collect external env var INPUT_CONFLUENCE_API_KEY
	ConfluenceAPIKeyEnv = "INPUT_CONFLUENCE_API_KEY"

	// ConfluenceSpaceEnv is to collect external env var INPUT_CONFLUENCE_SPACE
	ConfluenceSpaceEnv = "INPUT_CONFLUENCE_SPACE"

	// EnvsNotSetError is template env var not set error
	EnvsNotSetError = "environment variable not set, please assign values for: "
)
