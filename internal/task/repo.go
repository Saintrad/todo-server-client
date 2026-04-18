package task

import "github.com/Saintrad/todo-server-client/internal/richerror"

type TaskRepo interface {
	Create(Task) (Task, *richerror.AppError)
	List() ([]Task, *richerror.AppError)
	GetByID(int) (Task, *richerror.AppError)
	Update(Task) (Task, *richerror.AppError)
	Delete(int) (Task, *richerror.AppError)
}
