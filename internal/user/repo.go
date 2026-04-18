package user

import "github.com/Saintrad/todo-server-client/internal/richerror"

type repo interface{
	Register(User) (User, richerror.AppError)
	Login(email string, pw string) (User, richerror.AppError)
}