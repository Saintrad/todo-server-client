package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

type fileState struct {
	NextID int        `json:"next_id"`
	Tasks  []todo.Task `json:"tasks"`
}

type FileTaskRepo struct {
	mu       sync.Mutex
	filePath string
	state    fileState
}

var _ todo.TaskRepo = (*FileTaskRepo)(nil)

// NewFileTaskRepo loads state from file if present, otherwise starts empty.
func NewFileTaskRepo(path string) (*FileTaskRepo, error) {
	r := &FileTaskRepo{
		filePath: path,
		state: fileState{
			NextID: 1,
			Tasks:  make([]todo.Task, 0),
		},
	}

	// Ensure parent dir exists (e.g., data/)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Missing file is not an error: start empty
		if errors.Is(err, os.ErrNotExist) {
			return r, nil
		}
		return nil, err
	}

	// Empty file: treat as empty state
	if len(data) == 0 {
		return r, nil
	}

	var st fileState
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	// Defensive defaults
	if st.NextID <= 0 {
		st.NextID = computeNextID(st.Tasks)
	}
	if st.Tasks == nil {
		st.Tasks = make([]todo.Task, 0)
	}

	r.state = st
	return r, nil
}

func computeNextID(tasks []todo.Task) int {
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

// Create assigns an ID, stores the task, and persists to disk.
func (r *FileTaskRepo) Create(task todo.Task) (todo.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.state.NextID
	r.state.NextID++

	r.state.Tasks = append(r.state.Tasks, task)

	if err := r.saveLocked(); err != nil {
		return todo.Task{}, err
	}
	return task, nil
}

// saveLocked persists r.state to disk atomically.
// Call only while holding r.mu.
func (r *FileTaskRepo) saveLocked() error {
	b, err := json.MarshalIndent(r.state, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	dir := filepath.Dir(r.filePath)
	tmp, err := os.CreateTemp(dir, "tasks-*.tmp")
	if err != nil {
		return err
	}

	tmpName := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName) // no-op if rename succeeded
	}()

	if _, err := tmp.Write(b); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// Atomic replace on most OS/filesystems when same directory
	if err := os.Rename(tmpName, r.filePath); err != nil {
		return err
	}
	return nil
}
