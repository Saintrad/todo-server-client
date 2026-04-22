package httpapi

import "github.com/Saintrad/todo-server-client/internal/user"

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}



type UserLoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
}


type UserLoginResponse struct {
    Token string          `json:"token"`
    User  *UserPublicView `json:"user"`
}

// UserPublicView is reused in several endpoints.
type UserPublicView struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}

func (r UserRegisterRequest) ToDomain() user.RegisterInput {
	return user.RegisterInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r UserLoginRequest) ToDomain() user.LoginInput {
	return user.LoginInput{
		Email:    r.Email,
		Password: r.Password,
	}
}