package todo

import (
	"errors"
	"testing"
)

type fakeRepo struct {
	tasks  []Task
	nextID int
}

// fakeRepo constructor
func NewFakeRepo() *fakeRepo {
	return &fakeRepo{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}
func (r *fakeRepo) Create(t Task) (Task, error) {
	t.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, t)

	return t, nil
}

func (r *fakeRepo) List() ([]Task, error) {
	tasks := r.tasks

	return tasks, nil
}

func (r *fakeRepo) GetByID(id int) (Task, error) {

	for _, task := range r.tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return Task{}, ErrTaskNotFound
}

func (r *fakeRepo) Update(t Task) (Task, error) {

	for idx, task := range r.tasks {
		if task.ID == t.ID {
			r.tasks[idx] = t

			return t, nil
		}
	}

	return Task{}, ErrTaskNotFound
}

func (r *fakeRepo) Delete(id int) (Task, error) {

	for idx, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:idx], r.tasks[idx+1:]...)

			return task, nil
		}
	}

	return Task{}, ErrTaskNotFound
}

func TestCreateTaskEmptyTitle(t *testing.T) {
	input := CreateTaskInput{
		Title:    "",
		Category: nil,
		DueDate:  nil,
	}

	r := NewFakeRepo()
	taskService := NewService(r)

	_, err := taskService.CreateTask(input)

	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected error %v, got %v", ErrEmptyTitle, err)
	}

}

func TestCreateTaskOkTask(t *testing.T) {
	input := CreateTaskInput{
		Title:    "task",
		Category: nil,
		DueDate:  nil,
	}

	r := NewFakeRepo()
	taskService := NewService(r)

	currID := r.nextID
	createdTask, err := taskService.CreateTask(input)

	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	// Check if ID is assigned currectly
	if createdTask.ID != currID {
		t.Fatalf("expected ID to be %d, got %d", currID, createdTask.ID)
	}

	// Check if title is assigned correctly
	if createdTask.Title != input.Title {
		t.Fatalf("expected title to be %s, got %s", input.Title, createdTask.Title)
	}

	// Check if category is assigned correctly
	if createdTask.Category != input.Category {
		t.Fatalf("expected category to be %v, got %v", input.Category, createdTask.Category)
	}

	found := false
	for _, task := range r.tasks {
		if task.ID == createdTask.ID {
			found = true
		}
	}

	// Check if the task is in repo tasks
	if !found {
		t.Fatalf("task is not in repo tasks")
	}

}

func TestListTask(t *testing.T) {
	r := NewFakeRepo()
	taskService := NewService(r)

	tasks, err := taskService.ListTask()

	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	// Check if function works with no tasks
	if len(tasks) != len(r.tasks) {
		t.Fatalf("the size of listed tasks doesn't match the size of tasks")
	}

	r.tasks = append(r.tasks, Task{
		Title: "first",
	})

	tasks, err = taskService.ListTask()

	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	// Check if function works with tasks
	if len(tasks) != len(r.tasks) {
		t.Fatalf("the size of listed tasks doesn't match the size of tasks")
	}
}

func TestGetByID(t *testing.T) {
	r := NewFakeRepo()
	taskService := NewService(r)
	r.Create(Task{
		Title: "first",
	})

	// Check missing ID
	_, err := taskService.GetByID(2)

	if err == nil {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected error %v, got %v", ErrTaskNotFound, err)
	}

	// Check existing ID
	task, err := taskService.GetByID(1)

	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if task.ID != 1 {
		t.Fatalf("returned task with wrong ID")
	}

	if task.Title != "first" {
		t.Fatalf("the returned task is incorrect")
	}
}

func TestUpdateTask(t *testing.T) {

	r := NewFakeRepo()
	s := NewService(r)
	done := true

	// Check missing task
	_, err := s.UpdateTask(1, UpdateTaskInput{IsDone: &done})

	if err == nil {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected %v, got %v", ErrTaskNotFound, err)
	}

	// Check existing task
	r.Create(Task{})

	
	input := UpdateTaskInput{
		Title:    strPtr("changed"),
		IsDone:   &done,
		Category: strPtr("changed"),
	}
	_, err = s.UpdateTask(1, input)
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	task, _ := s.GetByID(1)

	if !task.IsDone {
		t.Fatalf("task done status was not updated")
	}

	if task.Title != *input.Title {
		t.Fatalf("task title was not updated")
	}

	if task.Category != input.Category {
		t.Fatalf("task category was not updated")
	}

}

func TestDeleteTask(t *testing.T) {

	r := NewFakeRepo()
	s := NewService(r)

	// Check delete missing task
	_, err := s.Delete(1)

	if err == nil {
		t.Fatalf("expected error %v, got %v", ErrTaskNotFound, err)
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected error %v, got %v", ErrTaskNotFound, err)
	}

	//Check delete existing task
	r.Create(Task{
		Title: "test",
	})

	_, err = s.Delete(1)

	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	// Check if the task is deleted
	_, err = r.GetByID(1)

	if err == nil {
		t.Fatalf("expected error %v, got %v", ErrTaskNotFound, err)
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected error %v, got %v", ErrTaskNotFound, err)
	}

}
