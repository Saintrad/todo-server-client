package todo

import (
	"time"
)

type Task struct{
	ID int
	Title string
	Category *string
	DueDate *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDone bool
}
