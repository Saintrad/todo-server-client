package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	go func() {
		if err := api.Start(addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("HTTP server stopped")

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("failed to get sql.DB: %v", err)
	} else {
		log.Println("Closing database...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
		log.Println("Database closed")
	}

	log.Println("Shutdown complete")

	}
