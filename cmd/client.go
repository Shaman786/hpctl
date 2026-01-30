package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type APIClient struct {
	BaseURL    string
	AuthHeader string
	Client     *http.Client
}

func NewClient() *APIClient {
	// 1. URL from Env or Default (Fixes localhost issue)
	url := os.Getenv("HP_API_URL")
	if url == "" {
		url = "http://localhost:5000/api/v1"
	}

	// 2. Load Auth (Fixes security issue)
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
	req, _ := http.NewRequest("POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonBody))
	return c.do(req)
}

func (c *APIClient) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	if c.AuthHeader != "" {
		req.Header.Set("Authorization", c.AuthHeader)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
