package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is an HTTP client for talking to the xbridge server.
type Client struct {
	BaseURL    string
	Key        string
	HTTPClient *http.Client
}

// New creates a new xbridge client.
func New(host, key string) *Client {
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	return &Client{
		BaseURL:    host,
		Key:        key,
		HTTPClient: &http.Client{},
	}
}

// Get performs a GET request.
func (c *Client) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	c.setAuth(req)
	return c.HTTPClient.Do(req)
}

// Post performs a POST request with JSON body.
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuth(req)
	return c.HTTPClient.Do(req)
}

// ReadJSON decodes a JSON response body into v.
func ReadJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) setAuth(req *http.Request) {
	if c.Key != "" {
		req.Header.Set("X-API-Key", c.Key)
	}
}
