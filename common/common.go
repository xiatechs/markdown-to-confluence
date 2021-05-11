// Package common is for storing common constants/vars used in app
//nolint:all // text in here could be long but on purpose i.e a long URL
package common

var (
	// ConfluenceBaseURL is the base URL for the confluence page you want to work with
	ConfluenceBaseURL = "https://xiatech-markup.atlassian.net"

	// ConfluenceUsername is to collect external arg for confluence username
	ConfluenceUsername = "INPUT_CONFLUENCE_USERNAME"

	// ConfluenceAPIKey is to collect external arg for api key
	ConfluenceAPIKey = "INPUT_CONFLUENCE_API_KEY"

	// ConfluenceSpace is to collect external arg for confluence space
	ConfluenceSpace = "INPUT_CONFLUENCE_SPACE"

	// EnvsNotSetError is template env arg not set error
	EnvsNotSetError = "args not set, please assign values for: "

	// ProjectPathEnv is to collect external arg for project path
	ProjectPathEnv = "./"
)
