package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

func TestFileRepo_NoFile_StartsEmptyAndCreatesFile(t *testing.T) {

	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")
	now := time.Now()

	repo, err := NewFileTaskRepo(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Call Create twice and assert IDs are 1 and 2.
	task1, err := repo.Create(todo.Task{
		Title:     "first",
		DueDate:   &now,
		CreatedAt: now,
		UpdatedAt: now,
		IsDone:    false,
	})

	if err != nil{
		t.Fatalf("expected no error, got %v", err)
	}

	if task1.ID != 1{
		t.Fatalf("expected ID 1, got %d", task1.ID)
	}

	task2, err := repo.Create(todo.Task{
		Title:     "second",
		DueDate:   &now,
		CreatedAt: now,
		UpdatedAt: now,
		IsDone:    false,
	})

	if task2.ID != 2 {
		t.Fatalf("expected ID 2, got %d", task2.ID)
	}

	// Assert file exists
	_, statErr := os.Stat(path)

	if statErr != nil {
		t.Fatalf("expected file to exist, got stat error %v", statErr)
	}

	}



