package dto

import (
	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name        string    `json:"name" validate:"required"`
	Permissions []uuid.UUID `json:"permissions,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string   `json:"name,omitempty" validate:"omitempty"`
	Permissions []uuid.UUID `json:"permissions,omitempty"`
}

type CreatePermissionRequest struct {
	Path string `json:"path" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdatePermissionRequest struct {
	Path *string `json:"path,omitempty" validate:"omitempty"`
	Name *string `json:"name,omitempty" validate:"omitempty"`
}