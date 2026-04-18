package task

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/richerror"
)

type Service struct {
	repo TaskRepo
}

func NewService(r TaskRepo) *Service {
	return &Service{repo: r}
}

func (s Service) CreateTask(i CreateTaskInput) (Task, *richerror.AppError) {

	// Check title not to be empty
	if i.Title == "" {
		return Task{}, richerror.InvalidInput("title must not be empty", nil)
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

func (s Service) ListTask() ([]Task, *richerror.AppError) {

	return s.repo.List()
}

func (s Service) GetByID(id int) (Task, *richerror.AppError) {

	return s.repo.GetByID(id)
}

func (s Service) UpdateTask(id int,i UpdateTaskInput) (Task, *richerror.AppError) {

	task, err := s.repo.GetByID(id)

	if err != nil {
		return Task{}, err
	}

	if i.Title != nil {
		task.Title = *i.Title
	}
	if i.DueDate != nil {
		task.DueDate = i.DueDate
	}
	if i.Category != nil {
		task.Category = i.Category
	}

	if i.IsDone != nil {
		task.IsDone = *i.IsDone
	}

	task.UpdatedAt = time.Now()

	return s.repo.Update(task)
}

func (s Service) Delete(id int) (Task, *richerror.AppError) {

	_, err := s.repo.GetByID(id)

	if err != nil {

		return Task{}, err
	}

	return s.repo.Delete(id)
}
