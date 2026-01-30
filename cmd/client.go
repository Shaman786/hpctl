package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// APIClient handles all communication with the backend
type APIClient struct {
	BaseURL    string
	AuthHeader string
	Client     *http.Client
}

// NewClient creates a client and automatically loads credentials
func NewClient() *APIClient {
	url := os.Getenv("HP_API_URL")
	if url == "" {
		url = "http://localhost:5000/api/v1"
	}

	auth, _ := LoadAuthHeader()

	return &APIClient{
		BaseURL:    url,
		AuthHeader: auth,
		Client:     &http.Client{},
	}
}

func (c *APIClient) Post(endpoint string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *APIClient) Get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *APIClient) Delete(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.BaseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

func (c *APIClient) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	if c.AuthHeader != "" {
		req.Header.Set("Authorization", c.AuthHeader)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	// LINT FIX: Explicitly ignore the close error to satisfy errcheck
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}
