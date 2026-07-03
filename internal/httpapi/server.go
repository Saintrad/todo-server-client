package httpapi

import (
	"context"
	"net/http"

	"github.com/Saintrad/todo-server-client/internal/auth"
	taskservice "github.com/Saintrad/todo-server-client/internal/service/task"
	userservice "github.com/Saintrad/todo-server-client/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	userSvc *userservice.UserSvc
	taskSvc *taskservice.TaskSvc
	router  *gin.Engine
	validator *validator.Validate
    server *http.Server
    jwt       *auth.JWTManager
}

func NewServer(u *userservice.UserSvc, t *taskservice.TaskSvc, v *validator.Validate, jwt *auth.JWTManager) *Server {
	r := gin.New()

	r.Use(LoggerMiddleware())
	r.Use(gin.Recovery())

	s := &Server{userSvc: u, taskSvc: t, validator: v, jwt: jwt, router: r}
	s.registerRoutes()
	return s
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) Start(addr string) error {
    s.server = &http.Server{
        Addr:    addr,
        Handler: s.router,
    }

    err := s.server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        return err
    }

    return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}

func (s *Server) registerRoutes() {
	authMW := NewAuthMiddleware(s.jwt)

    s.router.GET("/health", s.healthHandler)

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
        users.POST("/login", s.userLoginHandler)
    }
}
