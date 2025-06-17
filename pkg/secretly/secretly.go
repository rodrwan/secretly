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
	variables  map[string]interface{}
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
func (c *Client) GetEnv() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/env", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get env: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&c.variables); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return c.variables, nil
}

// LoadToEnvironment loads the retrieved variables into the current process environment
func (c *Client) LoadToEnvironment() error {
	if c.variables == nil {
		return fmt.Errorf("no variables to load")
	}

	for key, value := range c.variables {
		os.Setenv(key, value.(string))
	}

	return nil
}

// Get retrieves a specific environment variable
func (c *Client) Get(key string) (string, error) {
	if c.variables == nil {
		if _, err := c.GetEnv(); err != nil {
			return "", err
		}
	}

	if value, ok := c.variables[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue, nil
		}
		return "", fmt.Errorf("variable %s is not a string", key)
	}

	return "", fmt.Errorf("variable %s not found", key)
}

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	return err != nil && err.Error() == "variable not found"
}

// IsUnauthorized checks if the error is an "unauthorized" error
func IsUnauthorized(err error) bool {
	return err != nil && err.Error() == "unauthorized"
}

func (c *Client) GetAll() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/env", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get env: %s", resp.Status)
	}

	var env map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}

	return env, nil
}
