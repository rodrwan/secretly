package secretly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const defaultBaseURL = "http://localhost:8080/api/v1"

type Client struct {
	BaseURL string

	HTTPClient *http.Client

	variables map[string]interface{}
}

type ClientOption func(*Client)

func OptionBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

func New(opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) GetEnv() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/env", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get env: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&c.variables); err != nil {
		return nil, err
	}

	return c.variables, nil
}

func (c *Client) LoadToEnvitonemt() error {
	if c.variables == nil {
		return fmt.Errorf("no variables to load")
	}

	for key, value := range c.variables {
		os.Setenv(key, value.(string))
	}

	return nil
}
