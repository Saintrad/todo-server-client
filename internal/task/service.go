package task

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/richerror"
)

type TaskSvc struct {
	repo TaskRepo
}

func NewTaskSvc(r TaskRepo) *TaskSvc {
	return &TaskSvc{repo: r}
}

func (s TaskSvc) CreateTask(i CreateTaskInput) (Task, *richerror.AppError) {

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

func (s TaskSvc) ListTask() ([]Task, *richerror.AppError) {

	return s.repo.List()
}

func (s TaskSvc) GetByID(id int) (Task, *richerror.AppError) {

	return s.repo.GetByID(id)
}

func (s TaskSvc) UpdateTask(id int,i UpdateTaskInput) (Task, *richerror.AppError) {

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

func (s TaskSvc) Delete(id int) (Task, *richerror.AppError) {

	_, err := s.repo.GetByID(id)

	if err != nil {

		return Task{}, err
	}

	return s.repo.Delete(id)
}
