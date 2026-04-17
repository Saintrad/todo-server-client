package todo

import "time"

type CreateTaskInput struct{
	Title string
	Category *string
	DueDate *time.Time
}


type UpdateTaskInput struct {
	id int
	Title *string
	Category *string
	DueDate *time.Time
	IsDone *bool
}