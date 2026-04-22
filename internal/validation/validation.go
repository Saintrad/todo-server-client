package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("password", passwordValidator)
	return v
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var re = regexp.MustCompile(`^[A-Za-z0-9]{8,}$`)
	return re.MatchString(password)
}
