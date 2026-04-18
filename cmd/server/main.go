package main

import (
	"log"

	"github.com/Saintrad/todo-server-client/internal/httpapi"
	"github.com/Saintrad/todo-server-client/internal/task"
)


func main() {
	repo, err := task.NewFileTaskRepo("data/tasks.JSON")
	if err != nil {
		log.Fatal(err)
	}

	svc := task.NewService(repo)
	api := httpapi.NewServer(svc)

	log.Println("listening on :8080")
	log.Fatal(api.Start(":8080"))
}