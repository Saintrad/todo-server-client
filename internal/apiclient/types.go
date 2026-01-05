package apiclient

import "time"

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
