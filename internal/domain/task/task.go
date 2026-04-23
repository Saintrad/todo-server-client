package task

import (
	"time"
)

type Task struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	UserID    int `gorm:"index"`
	Title     string
	Category  *string
	DueDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDone    bool
}
