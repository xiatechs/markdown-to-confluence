package confluence

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
)

//go:generate mockgen --source=api.go -package confluencemocks -destination=test/confluencemocks/api.go
var constantsHardCoded = true

const (
	confluenceUsernameEnv = "INPUT_CONFLUENCE_USERNAME"
	confluenceAPIKeyEnv   = "INPUT_CONFLUENCE_API_KEY"
	confluenceSpaceEnv    = "INPUT_CONFLUENCE_SPACE"
	envsNotSetError       = "environment variable not set, please assign values for: "
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
		BaseURL:  "https://xiatech-markup.atlassian.net",
		Space:    lookupEnv(confluenceSpaceEnv),
		Username: lookupEnv(confluenceUsernameEnv),
		Password: lookupEnv(confluenceAPIKeyEnv),
		Client:   httpClient,
	}
}

// lookupEnv checks the environment variables required for creating the client have been set
func lookupEnv(env string) string {
	if !constantsHardCoded {
		v, exists := os.LookupEnv(env)
		if !exists {
			log.Printf("Environment variable not set for: %s", env)
			return ""
		}

		return v
	}

	return env
}
