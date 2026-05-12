package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSupplierCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"omitempty"`
}

type UpdateSupplierCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"omitempty"`
}

type SupplierCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
