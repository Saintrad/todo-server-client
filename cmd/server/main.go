package main

import (
	"log"
	"time"

	"github.com/Saintrad/todo-server-client/internal/auth"
	"github.com/Saintrad/todo-server-client/internal/httpapi"
	"github.com/Saintrad/todo-server-client/internal/task"
	"github.com/Saintrad/todo-server-client/internal/user"
	"github.com/Saintrad/todo-server-client/internal/validation"
)


func main() {
	taskRepo, err := task.NewFileTaskRepo("data/tasks.JSON")
	if err != nil {
		log.Fatal(err)
	}

	userRepo, err := user.NewFileUserRepo("data/users.JSON")
	if err != nil {
		log.Fatal(err)
	}

	t := task.NewTaskSvc(taskRepo)
	u := user.NewUserSvc(userRepo)
	v := validation.New()
	j := auth.NewJWTManager("secret123", time.Hour*24)
	api := httpapi.NewServer(u, t, v, j)

	log.Println("listening on :8080")
	log.Fatal(api.Start(":8080"))
}