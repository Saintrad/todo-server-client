package todo


type TaskRepo interface {
	Create(Task) (Task, error)
	List() []Task
	GetByID(int) (Task, error)
}