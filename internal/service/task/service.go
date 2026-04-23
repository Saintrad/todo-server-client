package taskservice

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/domain/task"
	"github.com/Saintrad/todo-server-client/internal/richerror"
)

type TaskSvc struct {
	repo task.TaskRepo
}

func NewTaskSvc(r task.TaskRepo) *TaskSvc {
	return &TaskSvc{repo: r}
}

func (s TaskSvc) CreateTask(userID int, i task.CreateTaskInput) (*task.Task, *richerror.AppError) {

	// Check title not to be empty
	if i.Title == "" {
		return nil, richerror.InvalidInput("title must not be empty", nil)
	}

	//Create a new task and initialize the attributes
	newTask := task.Task{
		UserID:    userID,
		Title:     i.Title,
		Category:  i.Category,
		DueDate:   i.DueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDone:    false,
	}

	task, err := s.repo.Create(newTask)
	if err != nil {
		return nil, richerror.Internal("failed to create task", err)
	}

	return &task, nil
}

func (s TaskSvc) ListTask(userID int) ([]task.Task, *richerror.AppError) {

	tasks, err := s.repo.List(userID)
	if err != nil {
		return nil, richerror.Internal("failed to fetch tasks", err)
	}

	return tasks, nil
}

func (s TaskSvc) GetByID(userID, id int) (*task.Task, *richerror.AppError) {

	exists, err := s.repo.IDExists(userID, id)
	if err != nil {
		return nil, richerror.Internal("failed to check existence", err)
	}

	if !exists {
		return nil, richerror.NotFound("task Not Found", nil)
	}

	task, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, richerror.Internal("failed to fetch task", err)
	}

	return &task, nil
}

func (s TaskSvc) UpdateTask(userID, id int, i task.UpdateTaskInput) (*task.Task, *richerror.AppError) {

	exists, err := s.repo.IDExists(userID, id)

	if err != nil {
		return nil, richerror.Internal("failed to check existence", err)
	}

	if !exists {
		return nil, richerror.NotFound("task not found", nil)
	}

	task, err := s.repo.GetByID(userID, id)

	if err != nil {
		return nil, richerror.Internal("failed to fetch task", err)
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

	updatedTask, err := s.repo.Update(task)
	if err != nil {
		return nil, richerror.Internal("failed to update task", err)
	}

	return &updatedTask, nil
}

func (s TaskSvc) Delete(userID, id int) *richerror.AppError {

	exists, err := s.repo.IDExists(userID, id)

	if err != nil {

		return richerror.Internal("failed to check existence", err)
	}

	if !exists {
		return richerror.NotFound("task not found", nil)
	}

	err = s.repo.Delete(userID, id)
	if err != nil {
		return richerror.Internal("failed to delete task", err)
	}

	return nil
}
