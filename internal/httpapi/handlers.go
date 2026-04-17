package httpapi

import (
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Saintrad/todo-server-client/internal/todo"
)

func (s *Server) createTaskHandler(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.BindJSON(&req); err != nil {
		WriteError(c, todo.InvalidInput("invalid create task payload", err))
		return
	}

	task, cErr := s.svc.CreateTask(req.ToDomain())
	if cErr != nil {
		WriteError(c, cErr)
		return
	}

	c.JSON(http.StatusCreated, ToTaskResponse(task))
}

func (s *Server) getByIdHandler(c *gin.Context) {
	// Extract the :id parameter from path
    idStr := c.Param("id")

    // Convert to int
    id, err := strconv.Atoi(idStr)
    if err != nil {
        WriteError(c, todo.InvalidInput("invalid task id", err))
        return
    }

    task, err := s.svc.GetByID(id)
    if err != nil {
        WriteError(c, err)
		return
    }

    c.JSON(http.StatusOK, ToTaskResponse(task))
}

func (s *Server) updateTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
        WriteError(c, todo.InvalidInput("invalid task id", err))
        return
    }

	var req UpdateTaskRequest
	if err := c.BindJSON(&req); err != nil {
		WriteError(c, todo.InvalidInput("invalid update payload", err))
	}

	task, uErr := s.svc.UpdateTask(id, req.ToDomain())
	if uErr != nil {
		WriteError(c, uErr)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (s *Server) listTaskHandler(c *gin.Context) {
	tasks, err := s.svc.ListTask()
	if err != nil {
		WriteError(c, err)
		return
	}

	responses := make([]TaskResponse, len(tasks))
	for i := range tasks {
		responses[i] = ToTaskResponse(tasks[i])
	}
	c.JSON(http.StatusOK, responses)
}

func (s *Server) deleteTaskHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
        WriteError(c, todo.InvalidInput("invalid task id", err))
        return
    }

	task, dErr := s.svc.Delete(id)
	if dErr != nil {
		WriteError(c, dErr)
		return
	}

	c.JSON(http.StatusNoContent, task)
}