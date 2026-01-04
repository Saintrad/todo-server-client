package httpapi

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

// POST /v1/tasks
type CreateTaskRequest struct {
	Title    string     `json:"title"`
	Category *string    `json:"category,omitempty"`
	DueDate  *time.Time `json:"due_date,omitempty"`
}

// PATCH /v1/tasks/{id}
type UpdateTaskRequest struct {
	Title    *string    `json:"title,omitempty"`
	Category *string    `json:"category,omitempty"`
	DueDate  *time.Time `json:"due_date,omitempty"`
	IsDone   *bool      `json:"is_done,omitempty"`
}

type TaskResponse struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Category  *string    `json:"category,omitempty"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	IsDone    bool       `json:"is_done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// ---------- Mapping helpers (DTO -> domain) ----------

func (r CreateTaskRequest) ToDomain() todo.CreateTaskInput {
	return todo.CreateTaskInput{
		Title:    r.Title,
		Category: r.Category,
		DueDate:  r.DueDate,
	}
}

func (r UpdateTaskRequest) ToDomain() todo.UpdateTaskInput {
	return todo.UpdateTaskInput{
		Title:    r.Title,
		Category: r.Category,
		DueDate:  r.DueDate,
		IsDone:   r.IsDone,
	}
}

// ---------- Mapping helper (domain -> DTO) ----------

func ToTaskResponse(t todo.Task) TaskResponse {
	return TaskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Category:  t.Category,
		DueDate:   t.DueDate,
		IsDone:    t.IsDone,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
