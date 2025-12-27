package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

const BaseURL = "https://api.host-palace.com"

// Client handles all communication with Host-Palace
type Client struct {
	Username string
	Password string
}

// NewClient creates a client using credentials stored in Viper config
func NewClient() *Client {
	return &Client{
		Username: viper.GetString("auth.username"),
		Password: viper.GetString("auth.password"),
	}
}

// Get makes a GET request to the API
func (c *Client) Get(endpoint string) ([]byte, error) {
	client := &http.Client{}

	// Construct the full URL (e.g., https://api.host-palace.com/details)
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add Basic Auth Header
	auth := c.Username + ":" + c.Password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	req.Header.Add("Content-Type", "application/json")

	// Send Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read Response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API Error: %s", resp.Status)
	}

	return body, nil
}
