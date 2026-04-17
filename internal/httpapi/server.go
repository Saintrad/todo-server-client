package httpapi

import (
	"github.com/Saintrad/todo-server-client/internal/todo"
	"github.com/gin-gonic/gin"
)


type Server struct {
	svc *todo.Service
	router *gin.Engine
}

func NewServer(svc todo.Service) *Server {
    r := gin.New()
    
    r.Use(LoggerMiddleware())     // custom structured logger
    r.Use(gin.Recovery())         // panic recovery

    s := &Server{svc: &svc, router: r}
    s.registerRoutes()
    return s
}


func (s *Server) Start (addr string) error{
	return s.router.Run(addr)
}

func (s *Server) registerRoutes() {
	tasks := s.router.Group("v1/tasks")
	{
		tasks.GET("", s.listTaskHandler)
		tasks.POST("", s.createTaskHandler)
		tasks.GET("/:id", s.getByIdHandler)
		tasks.PATCH("/:id", s.updateTaskHandler)
		tasks.DELETE("/:id", s.deleteTaskHandler)
	}
}

