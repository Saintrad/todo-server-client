package userstorage

import (

	"github.com/Saintrad/todo-server-client/internal/domain/user"
	"gorm.io/gorm"
)

type UserDBRepo struct {
    db *gorm.DB
}

func NewUserDBRepo(db *gorm.DB) *UserDBRepo {
    return &UserDBRepo{db: db}
}

func (r *UserDBRepo) Create(u user.User) (user.User, error) {
    if err := r.db.Create(&u).Error; err != nil {
        return user.User{}, err
    }
    return u, nil
}

func (r *UserDBRepo) GetByEmail(email string) (user.User, error) {
    var u user.User
    if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
        return user.User{}, err
    }
    return u, nil
}

func (r *UserDBRepo) EmailExists(email string) (bool, error) {
    var count int64

    err := r.db.Model(&user.User{}).
        Where("email = ?", email).
        Count(&count).Error

    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func (r *UserDBRepo) FindByID(id int) (user.User, error) {
    var u user.User
    if err := r.db.First(&u, id).Error; err != nil {
        return user.User{}, err
    }
    return u, nil
}
