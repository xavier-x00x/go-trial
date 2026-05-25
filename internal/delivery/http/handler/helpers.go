package handler

import (
	"errors"

	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func handleFieldErrors(c *fiber.Ctx, err error) bool {
	var fieldErrors *usecase.FieldErrors
	if err == nil {
		return false
	}
	if errors.As(err, &fieldErrors) {
		errs := make([]validator.ValidationError, len(fieldErrors.Errors))
		for i, fe := range fieldErrors.Errors {
			errs[i] = validator.ValidationError{Field: fe.Field, Message: fe.Message}
		}
		response.ValidationError(c, "Validation failed", errs)
		return true
	}
	return false
}
