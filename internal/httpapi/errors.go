package httpapi

import (
	"errors"
	"net/http"

	"log/slog"

	"github.com/Saintrad/todo-server-client/internal/richerror"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func WriteError(c *gin.Context, err error) {
	var appErr *richerror.AppError

	if errors.As(err, &appErr) {
		// structured logging
		slog.Error("request failed",
			"code", appErr.Code,
			"message", appErr.Message,
			"error", appErr.Err,
			"path", c.FullPath(),
		)

		status := statusCodeFor(appErr.Code)

		c.JSON(status, ErrorResponse{
			Error: appErr.Message,
			Code:  string(appErr.Code),
		})
		return
	}

	// Fallback: unexpected error
	slog.Error("internal error",
		"error", err,
		"path", c.FullPath(),
	)

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
		Code:  "INTERNAL",
	})
}

func statusCodeFor(code richerror.ErrorCode) int {
	switch code {
	case richerror.ErrCodeNotFound:
		return http.StatusNotFound
	case richerror.ErrCodeInvalidInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
