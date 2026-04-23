package httpapi

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/domain/task"
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

// ---------- Mapping helpers (DTO -> domain) ----------

func (r CreateTaskRequest) ToDomain() task.CreateTaskInput {
	return task.CreateTaskInput{
		Title:    r.Title,
		Category: r.Category,
		DueDate:  r.DueDate,
	}
}

func (r UpdateTaskRequest) ToDomain() task.UpdateTaskInput {
	return task.UpdateTaskInput{
		Title:    r.Title,
		Category: r.Category,
		DueDate:  r.DueDate,
		IsDone:   r.IsDone,
	}
}

// ---------- Mapping helper (domain -> DTO) ----------

func ToTaskResponse(t task.Task) TaskResponse {
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
