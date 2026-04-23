package task

type TaskRepo interface {
	Create(Task) (Task, error)
	List(int) ([]Task, error)
	GetByID(int, int) (Task, error)
	IDExists(int, int) (bool, error)
	Update(Task) (Task, error)
	Delete(int, int) error
}
