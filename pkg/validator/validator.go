package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator for struct validation.
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidationError represents a single field validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates a struct and returns user-friendly error messages.
func (v *Validator) Validate(s interface{}) []ValidationError {
	var errs []ValidationError

	err := v.validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: msgForTag(e),
			})
		}
	}

	return errs
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", toSnakeCase(fe.Field()))
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", toSnakeCase(fe.Field()), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", toSnakeCase(fe.Field()), fe.Param())
	case "unique":
		return fmt.Sprintf("%s must be unique", toSnakeCase(fe.Field()))
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", toSnakeCase(fe.Field()), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", toSnakeCase(fe.Field()))
	}
}

// toSnakeCase converts PascalCase/camelCase to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		isUpper := r >= 'A' && r <= 'Z'
		prevIsLower := i > 0 && s[i-1] >= 'a' && s[i-1] <= 'z'
		nextIsLower := i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z'

		if i > 0 && isUpper && (prevIsLower || nextIsLower) {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
