package confluence

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
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
}

// CreateAPIClient creates the API client with relevant login details for confluence's API
func CreateAPIClient() (*APIClient, error) {
	apiClient := APIClientWithAuths(retryablehttp.NewClient())
	if apiClient.Password == "" ||
		apiClient.Username == "" ||
		apiClient.Space == "" {
		return nil, fmt.Errorf("%s", "one or more arguments are not set - please ensure they are before running this action")
	}

	return apiClient, nil
}

// APIClientWithAuths returns an APIClient with dependencies defaulted to sane values
func APIClientWithAuths(httpClient HTTPClient) *APIClient {
	return &APIClient{
		BaseURL:  os.Getenv("CONFLUENCE_BASE_URL"),
		Space:    os.Getenv("CONFLUENCE_SPACE"),
		Username: os.Getenv("CONFLUENCE_USERNAME"),
		Password: os.Getenv("CONFLUENCE_API_KEY"),
		Client:   httpClient,
	}
}
