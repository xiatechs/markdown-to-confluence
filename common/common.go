// Package common is for storing common constants/vars used in app
package common

var (
	// ConfluenceBaseURL is the base URL for the confluence page you want the API to connect to
	// by default it is https://xiatech-markup.atlassian.net but can be changed below
	ConfluenceBaseURL = "https://xiatech-markup.atlassian.net"
	// ConfluenceUsername is to collect external arg for confluence username
	ConfluenceUsername string

	// ConfluenceAPIKey is to collect external arg for api key
	ConfluenceAPIKey string

	// ConfluenceSpace is to collect external arg for confluence space
	ConfluenceSpace string

	// ProjectPathEnv is to collect external arg for project path
	ProjectPathEnv string
)
