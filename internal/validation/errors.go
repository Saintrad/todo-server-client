package validation

import (
	"fmt"

	"github.com/Saintrad/todo-server-client/internal/richerror"
	"github.com/go-playground/validator/v10"
)


func TranslateError(err error) []richerror.AppError {
	var errs []richerror.AppError

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []richerror.AppError{
			{Message: "invalid input",
			Code: richerror.ErrCodeInvalidInput,
			Err: err,
			},
		}
	}

	for _, e := range validationErrors {
		errs = append(errs, richerror.AppError{
			Code: richerror.ErrCodeInvalidInput,
			Message: e.Field() + " " + messageForTag(e),
			Err: e,
		})
	}

	return errs
}

func messageForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "password":
		return "must be at least 8 characters, contain letters and numbers"
	default:
		return "is invalid"
	}
}
