package main

import (
	"log"
	"time"

	"github.com/Saintrad/todo-server-client/internal/auth"
	"github.com/Saintrad/todo-server-client/internal/config"
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
	// Load Configurations
	// --------------------------
	cfg := config.Load()
	
	// --------------------------
	// Connect to PostgreSQL
	// --------------------------
	db, err := db.Connect(
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
	)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// --------------------------
	// Create repositories
	// --------------------------
	userRepo := userstorage.NewUserDBRepo(db)
	taskRepo := taskstorage.NewTaskDBRepo(db)

	// --------------------------
	// Create services
	// --------------------------
	u := userservice.NewUserSvc(userRepo)
	t := taskservice.NewTaskSvc(taskRepo)

	// --------------------------
	// Create validator
	// --------------------------
	v := validation.New()

	// --------------------------
	// Create authenticator
	// --------------------------
	j := auth.NewJWTManager(cfg.JWTSecret, time.Hour*24)
	
	api := httpapi.NewServer(u, t, v, j)

	addr := ":" + cfg.Port

	log.Println("Listening on", addr)

	log.Fatal(api.Start(addr))
}
