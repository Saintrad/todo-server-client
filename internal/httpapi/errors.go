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

		status := StatusCodeFor(appErr.Code)

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

func WriteValidationErrors(c *gin.Context, errs []richerror.AppError) {
    var fields = make(map[string]string)

    for _, e := range errs {
        fields[e.Message] = e.Err.Error()
    }

    slog.Warn("validation failed",
        "path", c.FullPath(),
        "errors", errs,
    )

    c.JSON(http.StatusBadRequest, gin.H{
        "error":  "validation_error",
        "message": "invalid input",
        "fields": fields,
    })
}

func StatusCodeFor(code richerror.ErrorCode) int {
	switch code {
	case richerror.ErrCodeNotFound:
		return http.StatusNotFound
	case richerror.ErrCodeInvalidInput:
		return http.StatusBadRequest
	case richerror.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
