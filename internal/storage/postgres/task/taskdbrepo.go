package taskstorage

import (
	"errors"

	"github.com/Saintrad/todo-server-client/internal/domain/task"
	"gorm.io/gorm"
)

type TaskDBRepo struct {
	db *gorm.DB
}

func NewTaskDBRepo(db *gorm.DB) *TaskDBRepo {
	return &TaskDBRepo{db: db}
}

func (r *TaskDBRepo) Create(t task.Task) (task.Task, error) {
	
	if err := r.db.Create(&t).Error; err != nil {
		return task.Task{}, err
	}
	return t, nil
}

func (r *TaskDBRepo) List(userID int) ([]task.Task, error) {
	var tasks []task.Task
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskDBRepo) GetByID(userID, id int) (task.Task, error) {
	var t task.Task
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return task.Task{}, &task.ErrNotFound{}
	}
	return t, err
}

func (r *TaskDBRepo) IDExists(userID, id int) (bool, error) {

	err := r.db.Where("id = ? AND user_id = ?", id, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (r *TaskDBRepo) Update(t task.Task) (task.Task, error) {
	// Make sure user-owned record exists
	var existing task.Task
	err := r.db.Where("id = ? AND user_id = ?", t.ID, t.UserID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return task.Task{}, &task.ErrNotFound{}
	} else if err != nil {
		return task.Task{}, err
	}

	// Keep same primary/owner
	t.UserID = existing.UserID
	if err := r.db.Save(&t).Error; err != nil {
		return task.Task{}, err
	}
	return t, nil
}

func (r *TaskDBRepo) Delete(userID, id int) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&task.Task{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return &task.ErrNotFound{}
	}
	return nil
}
