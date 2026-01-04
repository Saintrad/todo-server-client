package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Saintrad/todo-server-client/internal/todo"
)

type Server struct {
	svc todo.Service
}

func NewServer(svc todo.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/tasks", s.tasksHandler)     // exact path
	mux.HandleFunc("/v1/tasks/", s.taskByIDHandler) // prefix match

	return mux
}

func (s *Server) taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	tail := strings.TrimPrefix(path, "/v1/tasks/")
	if tail == path || tail == "" {
		http.NotFound(w, r)
		return
	}
	if strings.Contains(tail, "/") {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(tail)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, err := s.svc.GetByID(id)
		if err != nil {
			s.writeDomainError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, ToTaskResponse(task))

	case http.MethodPatch:
		var req UpdateTaskRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		// Optional: reject empty PATCH (no fields provided)
		if req.Title == nil && req.Category == nil && req.DueDate == nil && req.IsDone == nil {
			writeError(w, http.StatusBadRequest, "no fields provided for update")
			return
		}

		task, err := s.svc.UpdateTask(id, req.ToDomain())
		if err != nil {
			s.writeDomainError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, ToTaskResponse(task))

	case http.MethodDelete:
		_, err := s.svc.Delete(id)
		if err != nil {
			s.writeDomainError(w, err)
			return
		}
		// Common: 204 No Content on successful delete
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


func (s *Server) tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := s.svc.ListTask()
		if err != nil {
			s.writeDomainError(w, err)
			return
		}

		// Convert domain tasks to response DTOs
		out := make([]TaskResponse, 0, len(tasks))
		for _, t := range tasks {
			out = append(out, ToTaskResponse(t))
		}

		writeJSON(w, http.StatusOK, out)

	case http.MethodPost:
		var req CreateTaskRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		// Minimal validation at transport boundary (optional).
		// Domain should still validate as the source of truth.
		if req.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}

		task, err := s.svc.CreateTask(req.ToDomain())
		if err != nil {
			s.writeDomainError(w, err)
			return
		}

		// 201 for resource creation
		writeJSON(w, http.StatusCreated, ToTaskResponse(task))

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// writeDomainError maps domain sentinel errors to HTTP status codes and returns JSON error body.
func (s *Server) writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, todo.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, "task not found")
	case errors.Is(err, todo.ErrEmptyTitle):
		writeError(w, http.StatusBadRequest, "title is required")
	default:
		// Avoid leaking internal details to clients
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

// writeJSON writes a JSON response with status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a consistent JSON error shape.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
