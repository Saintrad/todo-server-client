package apiclient

import (
	"time"
	"fmt"
)

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Category  *string    `json:"category,omitempty"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	IsDone    bool       `json:"is_done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateTaskRequest struct {
	Title    string     `json:"title"`
	Category *string    `json:"category,omitempty"`
	DueDate  *time.Time `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	Title    *string    `json:"title,omitempty"`
	Category *string    `json:"category,omitempty"`
	DueDate  *time.Time `json:"due_date,omitempty"`
	IsDone   *bool      `json:"is_done,omitempty"`
}


type APIError struct {
	Status int    `json:"-"`
	Code   string `json:"code"`
	Msg    string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s (%s)", e.Msg, e.Code)
	}
	return e.Msg
}

