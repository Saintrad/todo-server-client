package user

import "time"

type User struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	Name         string
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
