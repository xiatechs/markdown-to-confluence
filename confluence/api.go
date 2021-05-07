package confluence

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xiatechs/markdown-to-confluence/common"
)

//go:generate mockgen --source=api.go -package confluencemocks -destination=test/confluencemocks/api.go

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
			common.EnvsNotSetError,
			common.ConfluenceAPIKeyEnv,
			common.ConfluenceSpaceEnv,
			common.ConfluenceUsernameEnv)
	}

	return apiClient, nil
}

// APIClientWithAuths returns an APIClient with dependencies defaulted to sane values
func APIClientWithAuths(httpClient HTTPClient) *APIClient {
	return &APIClient{
		BaseURL:  common.ConfluenceBaseURL,
		Space:    lookupEnv(common.ConfluenceSpaceEnv),
		Username: lookupEnv(common.ConfluenceUsernameEnv),
		Password: lookupEnv(common.ConfluenceAPIKeyEnv),
		Client:   httpClient,
	}
}

// lookupEnv checks the environment variables required for creating the client have been set
func lookupEnv(env string) string {
	if !common.ConstantsHardCoded {
		v, exists := os.LookupEnv(env)
		if !exists {
			log.Printf("Environment variable not set for: %s", env)
			return ""
		}

		return v
	}

	return env
}
