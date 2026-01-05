package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type APIError struct {
	Error string `json:"error"`
}

func (c *Client) do(method, path string, reqBody any, respBody any) (int, error) {
	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return 0, err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return 0, err
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Non-2xx: try to decode JSON error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(resp.Body)
		var apiErr APIError
		if json.Unmarshal(raw, &apiErr) == nil && apiErr.Error != "" {
			return resp.StatusCode, fmt.Errorf(apiErr.Error)
		}
		return resp.StatusCode, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	if respBody != nil {
		return resp.StatusCode, json.NewDecoder(resp.Body).Decode(respBody)
	}
	return resp.StatusCode, nil
}
