package httpapi

import (
	"github.com/Saintrad/todo-server-client/internal/auth"
	"github.com/Saintrad/todo-server-client/internal/task"
	"github.com/Saintrad/todo-server-client/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	userSvc *user.UserSvc
	taskSvc *task.TaskSvc
	router  *gin.Engine
	validator *validator.Validate
    jwt       *auth.JWTManager
}

func NewServer(u *user.UserSvc, t *task.TaskSvc, v *validator.Validate, jwt *auth.JWTManager) *Server {
	r := gin.New()

	r.Use(LoggerMiddleware())
	r.Use(gin.Recovery())

	s := &Server{userSvc: u, taskSvc: t, validator: v, jwt: jwt, router: r}
	s.registerRoutes()
	return s
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) registerRoutes() {
	authMW := NewAuthMiddleware(s.jwt)

    // =====================
    // Protected routes
    // =====================
    tasks := s.router.Group("/v1/tasks")
    tasks.Use(authMW.Handle())
    {
        tasks.GET("", s.listTaskHandler)
        tasks.POST("", s.createTaskHandler)
        tasks.GET("/:id", s.getByIdHandler)
        tasks.PATCH("/:id", s.updateTaskHandler)
        tasks.DELETE("/:id", s.deleteTaskHandler)
    }

    // =====================
    // Public routes
    // =====================
    users := s.router.Group("/v1/users")
    {
        users.POST("", s.userRegisterHandler)
        users.POST("/login", s.userLoginHandler) // if you have a login handler
    }
}
