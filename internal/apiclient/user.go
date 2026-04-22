package apiclient

import "net/http"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Register(req RegisterRequest) error {
	_, err := c.do(http.MethodPost, "/v1/users", req, nil)
	return err
}

func (c *Client) Login(req LoginRequest) (LoginResponse, error) {
	var out LoginResponse
	_, err := c.do(http.MethodPost, "/v1/users/login", req, &out)
	return out, err
}
