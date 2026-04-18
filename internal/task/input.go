package task

import "time"

type CreateTaskInput struct{
	Title string
	Category *string
	DueDate *time.Time
}


type UpdateTaskInput struct {
	Title *string
	Category *string
	DueDate *time.Time
	IsDone *bool
}