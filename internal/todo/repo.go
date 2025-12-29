package todo


type TaskRepo interface {
	Create(Task) (Task, error)
}