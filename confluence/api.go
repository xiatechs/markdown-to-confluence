package confluence

import (
	"fmt"
	"net/http"

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

// HTTPClient interface will allow mock Do request
type HTTPClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
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
			common.ConfluenceAPIKey,
			common.ConfluenceSpace,
			common.ConfluenceUsername)
	}

	return apiClient, nil
}

// APIClientWithAuths returns an APIClient with dependencies defaulted to sane values
func APIClientWithAuths(httpClient HTTPClient) *APIClient {
	return &APIClient{
		BaseURL:  common.ConfluenceBaseURL,
		Space:    common.ConfluenceSpace,
		Username: common.ConfluenceUsername,
		Password: common.ConfluenceAPIKey,
		Client:   httpClient,
	}
}
