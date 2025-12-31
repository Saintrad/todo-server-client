package todo

import (
	"time"
)

type Service struct {
	repo TaskRepo
}

func NewService(r TaskRepo) Service{
	return Service{repo: r}
}

func (s Service) CreateTask(i CreateTaskInput) (Task, error) {

	// Check title not to be empty
	if i.Title == "" {
		return Task{}, ErrEmptyTitle
	}

	//Create a new task and initialize the attributes
	newTask := Task{
		Title:     i.Title,
		Category:  i.Category,
		DueDate:   i.DueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDone:    false,
	}

	return s.repo.Create(newTask)
}

func (s Service) ListTask(i ListTaskInput) []Task {

	return s.repo.List()
}

func (s Service) GetByID(id int) (Task, error) {
	
	return s.repo.GetByID(id)
}
