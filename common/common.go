// Package common is for storing common constants/vars used in app
package common

const (
	// ConfluenceBaseURL is the base URL for the confluence page you want the API to connect to
	ConfluenceBaseURL = "https://xiatech-markup.atlassian.net"	
	
	// EnvsNotSetError is template env arg not set error
	EnvsNotSetError = "args not set, please assign values for: "
)

var (
	// ConfluenceUsername is to collect external arg for confluence username
	ConfluenceUsername string

	// ConfluenceAPIKey is to collect external arg for api key
	ConfluenceAPIKey string

	// ConfluenceSpace is to collect external arg for confluence space
	ConfluenceSpace string

	// ProjectPathEnv is to collect external arg for project path
	ProjectPathEnv string
)
