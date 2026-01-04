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
	NextID int         `json:"next_id"`
	Tasks  []todo.Task `json:"tasks"`
}

type FileTaskRepo struct {
	mu       sync.Mutex
	filePath string
	state    fileState
}

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

func (r *FileTaskRepo) List() ([]todo.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := make([]todo.Task, 0)

	for _, task := range r.state.Tasks {

		list = append(list, task)
	}

	return list, nil
}

func (r *FileTaskRepo) GetByID(id int) (todo.Task, error) {

	tasks := r.state.Tasks

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return todo.Task{}, todo.ErrTaskNotFound
}

func (r *FileTaskRepo) Update(t todo.Task) (todo.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var oldIdx int
	var oldTask todo.Task
	found := false

	// Find the task by ID
	for idx, task := range r.state.Tasks {
		if task.ID == t.ID {
			oldTask = task
			oldIdx = idx
			found = true
			r.state.Tasks[idx] = t

			break
		}
	}

	if !found {

		return todo.Task{}, todo.ErrTaskNotFound
	}

	// Save the state
	err := r.saveLocked()

	// If save fails, revert the update
	if err != nil {
		r.state.Tasks[oldIdx] = oldTask

		return todo.Task{}, err
	}

	return t, nil
}

func (r *FileTaskRepo) Delete(id int) (todo.Task, error) {

	var oldTask todo.Task
	found := false

	for idx, task := range r.state.Tasks {
		if task.ID == id {
			oldTask = task
			r.state.Tasks = append(r.state.Tasks[:idx], r.state.Tasks[idx+1:]...)
			found = true

			break
		}
	}

	if !found {

		return todo.Task{}, todo.ErrTaskNotFound
	}

	err := r.saveLocked()

	if err != nil {
		r.state.Tasks = append(r.state.Tasks, oldTask)

		return todo.Task{}, err
	}

	return oldTask, nil
}
