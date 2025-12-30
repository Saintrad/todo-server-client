package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

func strPtr(s string) *string {
	return &s
}

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

	if task2.ID != 2 {
		t.Fatalf("expected ID 2, got %d", task2.ID)
	}

	// Assert file exists
	_, statErr := os.Stat(path)

	if statErr != nil {
		t.Fatalf("expected file to exist, got stat error %v", statErr)
	}

}

func TestFileRepo_ExistingFile_LoadsStateAndContinuesIDs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")

	now := time.Now()
	later := now.Add(24 * time.Hour)

	state := fileState{
		NextID: 6,
		Tasks: []todo.Task{
			{
				ID:        1,
				Title:     "Buy milk",
				Category:  strPtr("personal"),
				IsDone:    false,
				DueDate:   &later,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				Title:     "Write unit tests",
				Category:  strPtr("work"),
				IsDone:    true,
				DueDate:   nil,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        3,
				Title:     "Refactor storage layer",
				Category:  strPtr("work"),
				IsDone:    false,
				DueDate:   nil,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        4,
				Title:     "Go for a run",
				Category:  strPtr("health"),
				IsDone:    true,
				DueDate:   &later,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        5,
				Title:     "Plan v2 auth",
				Category:  strPtr(""),
				IsDone:    false,
				DueDate:   nil,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	data, jErr := json.Marshal(state)
	if jErr != nil {
		t.Fatalf("expected no error, got %v", jErr)
	}
	os.WriteFile(path, data, 0644)

	repo, nErr := NewFileTaskRepo(path)
	if nErr != nil {
		t.Fatalf("expected no error, got %v", nErr)
	}

	task, cErr := repo.Create(todo.Task{
		Title: "task",
	})
	if cErr != nil {
		t.Fatalf("expected no error, got %v", cErr)
	}

	if task.ID != 6 {
		t.Fatalf("expected ID 6, got %d", task.ID)
	}

	if repo.state.NextID != 7 {
		t.Fatalf("expected next ID 7, got %d", repo.state.NextID)
	}
}
