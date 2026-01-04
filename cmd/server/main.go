package main

import (
	"log"
	"net/http"

	"github.com/Saintrad/todo-server-client/internal/httpapi"
	"github.com/Saintrad/todo-server-client/internal/storage"
	"github.com/Saintrad/todo-server-client/internal/todo"
)


func main() {
	repo, err := storage.NewFileTaskRepo("data/tasks.JSON")
	if err != nil {
		log.Fatal(err)
	}

	svc := todo.NewService(repo)
	api := httpapi.NewServer(svc)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.Routes()))
}