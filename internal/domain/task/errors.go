package task

type ErrNotFound struct {}

func (e *ErrNotFound) Error() string {
	return "task not found"
}