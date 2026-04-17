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

// do sends an HTTP request and decodes the JSON response.
// It returns the status code and either a typed APIError or nil.
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

	// Standard headers
	req.Header.Set("Accept", "application/json")
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Non‑2xx → structured error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if json.Unmarshal(raw, &apiErr) == nil && apiErr.Msg != "" {
			apiErr.Status = resp.StatusCode
			return resp.StatusCode, &apiErr
		}

		return resp.StatusCode, fmt.Errorf(
			"http %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(raw)),
		)
	}

	// No body expected
	if respBody == nil || len(raw) == 0 {
		return resp.StatusCode, nil
	}

	// Normal JSON decode
	if err := json.Unmarshal(raw, respBody); err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}
