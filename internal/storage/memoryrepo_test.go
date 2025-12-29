package storage

import (
	"testing"
	"time"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

func TestMemoryRepo(t *testing.T) {
	repo := NewMemoryTaskRepo()

	now := time.Now()
	task1, err := repo.Create(todo.Task{
		Title:     "first",
		DueDate:   &now,
		CreatedAt: now,
		UpdatedAt: now,
		IsDone:    false,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task1.ID != 1 {
		t.Fatalf("expected ID 1, got %d", task1.ID)
	}

	task2, err := repo.Create(todo.Task{
		Title:     "second",
		DueDate:   &now,
		CreatedAt: now,
		UpdatedAt: now,
		IsDone:    false,
	})

	if err != nil {
		t.Fatalf("excpected no error, got %v", err)
	}

	if task2.ID != 2 {
		t.Fatalf("expected ID 2, got %d", task2.ID)
	}
}
