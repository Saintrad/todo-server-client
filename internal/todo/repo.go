package todo

type TaskRepo interface {
	Create(Task) (Task, *AppError)
	List() ([]Task, *AppError)
	GetByID(int) (Task, *AppError)
	Update(Task) (Task, *AppError)
	Delete(int) (Task, *AppError)
}
