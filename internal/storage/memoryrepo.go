package storage

import (
	"sync"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

type MemoryTaskRepo struct {
	mu     sync.Mutex
	tasks  []todo.Task
	nextID int
}

func NewMemoryTaskRepo() *MemoryTaskRepo {
	return &MemoryTaskRepo{
		tasks:  make([]todo.Task, 0),
		nextID: 1,
	}
}

func (r *MemoryTaskRepo) Create(t todo.Task) (todo.Task, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	// Assign ID
	t.ID = r.nextID
	r.nextID++

	// Save task
	r.tasks = append(r.tasks, t)

	return t, nil
}
