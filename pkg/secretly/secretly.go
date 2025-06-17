package secretly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL = "http://localhost:8080"
	defaultTimeout = 10 * time.Second
)

// Client represents a Secretly client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets the timeout for HTTP requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{}
		}
		c.HTTPClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = client
	}
}

// New creates a new Secretly client with the given options
func New(opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetEnv retrieves all environment variables from the Secretly server
/*
Example respose:
[
	{
		"id": 1,
		"name": "development",
		"values": [
			{
				"id": 5,
				"key": "PORT",
				"value": "8080"
			}
		]
	},
	{
		"id": 2,
		"name": "staging",
		"values": [
			{
				"id": 6,
				"key": "PORT",
				"value": "8080"
			}
		]
	},
]
*/
func (c *Client) getAllEnvironments() ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/env", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get env: %s", resp.Status)
	}

	var environments []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&environments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return environments, nil
}

/*
Query params:
- name: string

Example respose:
[

	{
		"id": 1,
		"name": "development",
		"values": [
			{
				"id": 5,
				"key": "PORT",
				"value": "8080"
			}
		]
	}

]
*/
func (c *Client) getEnvironmentByName(environmentName string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/env?name=%s", c.BaseURL, environmentName)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get env: %s", resp.Status)
	}

	var environments []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&environments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return environments, nil
}

// LoadToEnvironment loads the retrieved variables into the current process environment
func (c *Client) LoadToEnvironment(environmentName string) error {
	environments, err := c.getAllEnvironments()
	if err != nil {
		return err
	}

	for _, environment := range environments {
		if environment["name"] == environmentName {
			for _, value := range environment["values"].([]map[string]interface{}) {
				os.Setenv(value["key"].(string), value["value"].(string))
			}
		}
	}

	return nil
}

func (c *Client) GetEnvironmentByName(environmentName string) (map[string]interface{}, error) {
	environments, err := c.getEnvironmentByName(environmentName)
	if err != nil {
		return nil, err
	}

	if len(environments) == 0 {
		return nil, fmt.Errorf("environment %s not found", environmentName)
	}

	return environments[0], nil
}

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	return err != nil && err.Error() == "variable not found"
}

// IsUnauthorized checks if the error is an "unauthorized" error
func IsUnauthorized(err error) bool {
	return err != nil && err.Error() == "unauthorized"
}

func (c *Client) GetAll() ([]map[string]interface{}, error) {
	environments, err := c.getAllEnvironments()
	if err != nil {
		return nil, err
	}

	return environments, nil
}
