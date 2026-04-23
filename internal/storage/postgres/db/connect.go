package db

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    
	"github.com/Saintrad/todo-server-client/internal/domain/task"
	"github.com/Saintrad/todo-server-client/internal/domain/user"
)

func Connect(host, username, pass, name string, port int) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
        host, username, pass, name, port,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Auto-migrate domain models
    if err := db.AutoMigrate(&user.User{}, &task.Task{}); err != nil {
        return nil, err
    }

    return db, nil
}
