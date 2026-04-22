package user

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/richerror"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc struct {
	repo UserRepo
}

func NewUserSvc(r UserRepo) *UserSvc {
	return &UserSvc{repo: r}
}

func (s UserSvc) Register(i RegisterInput) (User, *richerror.AppError) {

	exists, err := s.repo.EmailExists(i.Email)
	if err != nil {
		return User{}, richerror.Internal("internal error", err)
	}

	if exists {
		return User{}, richerror.InvalidInput("this email has already been used", nil)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, richerror.Internal("failed to hash password", err)
	}

	u := User{
		Name:         i.Name,
		Email:        i.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	id, err := s.repo.Create(u)
	if err != nil {
		return User{}, richerror.Internal("failed to create user", err)
	}
	u.ID = id

	return u, nil
}

func (s UserSvc) Login(i LoginInput) (User, *richerror.AppError) {

	u, err := s.repo.GetByEmail(i.Email)
	if err != nil {
		// We do NOT reveal whether the email exists.
		return User{}, richerror.InvalidInput("email or password incorrect", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(i.Password))
	if err != nil {
		return User{}, richerror.InvalidInput("email or password incorrect", err)
	}

	return u, nil
}
