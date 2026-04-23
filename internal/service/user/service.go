package userservice

import (
	"time"

	"github.com/Saintrad/todo-server-client/internal/domain/user"
	"github.com/Saintrad/todo-server-client/internal/richerror"
	"golang.org/x/crypto/bcrypt"
)

type UserSvc struct {
	repo user.UserRepo
}

func NewUserSvc(r user.UserRepo) *UserSvc {
	return &UserSvc{repo: r}
}

func (s UserSvc) Register(i user.RegisterInput) (*user.User, *richerror.AppError) {

	exists, err := s.repo.EmailExists(i.Email)
	if err != nil {
		return nil, richerror.Internal("internal error", err)
	}

	if exists {
		return nil, richerror.InvalidInput("this email has already been used", nil)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, richerror.Internal("failed to hash password", err)
	}

	u := user.User{
		Name:         i.Name,
		Email:        i.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.repo.Create(u)
	if err != nil {
		return nil, richerror.Internal("failed to create user", err)
	}
	u.ID = createdUser.ID

	return &u, nil
}

func (s UserSvc) Login(i user.LoginInput) (*user.User, *richerror.AppError) {

	u, err := s.repo.GetByEmail(i.Email)
	if err != nil {
		// We do NOT reveal whether the email exists.
		return nil, richerror.InvalidInput("email or password incorrect", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(i.Password))
	if err != nil {
		return nil, richerror.InvalidInput("email or password incorrect", err)
	}

	return &u, nil
}
