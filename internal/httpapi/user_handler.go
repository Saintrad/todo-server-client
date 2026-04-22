package httpapi

import (
	"net/http"
	"time"

	"github.com/Saintrad/todo-server-client/internal/validation"
	"github.com/gin-gonic/gin"
)

func (s *Server) userRegisterHandler(c *gin.Context) {
    var req UserRegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        WriteError(c, err)
		return
    }

    if err := s.validator.Struct(req); err != nil {
        appErrors := validation.TranslateError(err)
        WriteValidationErrors(c ,appErrors)
        return
    }


    newUser, rErr := s.userSvc.Register(req.ToDomain())
    if rErr != nil {
        WriteError(c, rErr)
    }

    c.JSON(http.StatusCreated, newUser)
}

func (s *Server) userLoginHandler(c *gin.Context) {
    var req UserLoginRequest

    // --- Parse JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        WriteError(c, err)
        return
    }

    // --- Validate DTO
    if err := s.validator.Struct(req); err != nil {
        appErrors := validation.TranslateError(err)
        WriteValidationErrors(c ,appErrors)
        return
    }

    // --- Domain login
    u, err := s.userSvc.Login(req.ToDomain())
    if err != nil {
        WriteError(c, err)
        return
    }

    // --- Generate JWT
    token, err := s.jwt.GenerateToken(u.ID)
    if err != nil {
        WriteError(c, err)
        return
    }

    // --- Response payload
    resp := UserLoginResponse{
        Token: token,
        User: &UserPublicView{
            ID:        u.ID,
            Name:      u.Name,
            Email:     u.Email,
            CreatedAt: u.CreatedAt.Format(time.RFC3339),
        },
    }

    c.JSON(http.StatusOK, resp)
}