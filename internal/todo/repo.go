package todo

type TaskRepo interface {
	Create(Task) (Task, error)
	List() ([]Task, error)
	GetByID(int) (Task, error)
	Update(Task) (Task, error)
	Delete(int) (Task, error)
}
