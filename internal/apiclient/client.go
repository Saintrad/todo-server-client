package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
    baseURL string
    client  *http.Client

    token   string       // JWT
    tokenFile string     // path to ~/.todo/token
}


func New(baseURL string) *Client {
    home, _ := os.UserHomeDir()
    tokenPath := filepath.Join(home, ".todo", "token")

    c := &Client{
        baseURL:   baseURL,
        client:    &http.Client{},
        tokenFile: tokenPath,
    }
    c.LoadToken()
    return c
}

func (c *Client) LoadToken() {
    data, err := os.ReadFile(c.tokenFile)
    if err == nil {
        c.token = strings.TrimSpace(string(data))
    }
}

func (c *Client) SaveToken() error {
    os.MkdirAll(filepath.Dir(c.tokenFile), 0700)
    return os.WriteFile(c.tokenFile, []byte(c.token), 0600)
}

func (c *Client) ClearToken() error {
    c.token = ""
    return os.Remove(c.tokenFile)
}

func (c *Client) SetToken(t string) error {
    c.token = t
    return c.SaveToken()
}


// do sends an HTTP request and decodes the JSON response.
// It returns the status code and either a typed APIError or nil.
func (c *Client) do(method, path string, body any, respBody any) (*http.Response, error) {
    var buf io.ReadWriter
    if body != nil {
        buf = new(bytes.Buffer)
        if err := json.NewEncoder(buf).Encode(body); err != nil {
            return nil, err
        }
    }

    req, err := http.NewRequest(method, c.baseURL+path, buf)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    // Add token if available
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        data, _ := io.ReadAll(resp.Body)
        return resp, fmt.Errorf("api error: %s", data)
    }

    if respBody != nil {
        return resp, json.NewDecoder(resp.Body).Decode(&respBody)
    }

    return resp, nil
}
