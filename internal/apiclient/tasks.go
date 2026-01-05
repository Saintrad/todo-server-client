package apiclient

import (
	"net/http"
	"strconv"
)

func (c *Client) ListTasks() ([]Task, error) {
	var out []Task
	_, err := c.do(http.MethodGet, "/v1/tasks", nil, &out)
	return out, err
}

func (c *Client) CreateTask(req CreateTaskRequest) (Task, error) {
	var out Task
	_, err := c.do(http.MethodPost, "/v1/tasks", req, &out)
	return out, err
}

func (c *Client) GetTask(id int) (Task, error) {
	var out Task
	_, err := c.do(http.MethodGet, "/v1/tasks/"+itoa(id), nil, &out)
	return out, err
}

func (c *Client) UpdateTask(id int, req UpdateTaskRequest) (Task, error) {
	var out Task
	_, err := c.do(http.MethodPatch, "/v1/tasks/"+itoa(id), req, &out)
	return out, err
}

func (c *Client) DeleteTask(id int) error {
	_, err := c.do(http.MethodDelete, "/v1/tasks/"+itoa(id), nil, nil)
	return err
}

// small helper to avoid fmt.Sprintf in hot paths
func itoa(n int) string {
	return strconv.Itoa(n)
}
