package todo

import (
	"errors"
	"time"
)


type Service struct{
	repo TaskRepo
}

func (s Service)CreateTask(i CreateTaskInput) (Task, error){

	// Check title not to be empty
	if i.Title == ""{
		return Task{}, errors.New("task title is empty")
	}

	//Create a new task and initialize the attributes
	newTask := Task{
		Title: i.Title,
		Category: i.Category,
		DueDate: i.DueDate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDone: false,
	}

	return s.repo.Create(newTask)
}