package main

import (
	"log"
	"time"

	"github.com/Saintrad/todo-server-client/internal/auth"
	"github.com/Saintrad/todo-server-client/internal/httpapi"
	taskservice "github.com/Saintrad/todo-server-client/internal/service/task"
	userservice "github.com/Saintrad/todo-server-client/internal/service/user"
	"github.com/Saintrad/todo-server-client/internal/storage/postgres/db"
	taskstorage "github.com/Saintrad/todo-server-client/internal/storage/postgres/task"
	userstorage "github.com/Saintrad/todo-server-client/internal/storage/postgres/user"
	"github.com/Saintrad/todo-server-client/internal/validation"
)

func main() {
	// --------------------------
	// Connect to PostgreSQL
	// --------------------------
	db, err := db.Connect(
		"localhost",
		"todo",
		"todo",
		"todo",
		5432,
	)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// --------------------------
	// Create repositories
	// --------------------------
	userRepo := userstorage.NewUserDBRepo(db)
	taskRepo := taskstorage.NewTaskDBRepo(db)

	u := userservice.NewUserSvc(userRepo)
	t := taskservice.NewTaskSvc(taskRepo)

	v := validation.New()
	j := auth.NewJWTManager("secret123", time.Hour*24)
	api := httpapi.NewServer(u, t, v, j)

	log.Println("listening on :8080")
	log.Fatal(api.Start(":8080"))
}
