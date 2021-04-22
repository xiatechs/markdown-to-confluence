package confluence

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/object"
)

//go:generate mockgen --source=api.go -package confluencemocks -destination=test/confluencemocks/api.go
var confluenceobject = object.ConfluenceObject //the struct containing the username/api/space vars from object

// here we are grabbing the variables in the confluence object
var (
	confluenceUsernameEnv = confluenceobject.ConfluenceUsernameEnv
	confluenceAPIKeyEnv   = confluenceobject.ConfluenceAPIKeyEnv
	confluenceSpaceEnv    = confluenceobject.ConfluenceSpaceEnv
	envsNotSetError       = "environment variable(s) not set"
)

// APIClient struct for interacting with confluence
type APIClient struct {
	BaseURL  string
	Space    string
	Username string
	Password string
	Client   HTTPClient
}

// HTTPClient is required to mock the http requests
type HTTPClient interface {
	Do(
		req *retryablehttp.Request,
	) (*http.Response, error)
}

// CreateAPIClient creates the API client with relevant login details for confluence's API
func CreateAPIClient() (*APIClient, error) {
	apiClient := APIClientWithAuths(retryablehttp.NewClient())
	if apiClient.Password == "" ||
		apiClient.Username == "" ||
		apiClient.Space == "" {
		return nil, fmt.Errorf("%s %s, %s, %s",
			envsNotSetError,
			confluenceAPIKeyEnv,
			confluenceSpaceEnv,
			confluenceUsernameEnv)
	}

	return apiClient, nil
}

// APIClientWithAuths returns an APIClient with dependencies defaulted to sane values
func APIClientWithAuths(httpClient HTTPClient) *APIClient {
	return &APIClient{
		BaseURL:  "https://xiatech.atlassian.net",
		Space:    confluenceSpaceEnv,
		Username: confluenceUsernameEnv,
		Password: confluenceAPIKeyEnv,
		Client:   httpClient,
	}
}
