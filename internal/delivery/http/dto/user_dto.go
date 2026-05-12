package dto

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	StoreID   *uuid.UUID `json:"store_id,omitempty"`
	Name      string    `json:"name" validate:"required,max=255"`
	Username  string    `json:"username" validate:"required,max=100"`
	Email     string    `json:"email" validate:"required,email,max=255"`
	Phone     *string   `json:"phone,omitempty" validate:"omitempty,max=20"`
	Password  string    `json:"password" validate:"required,min=8"`
	Role      string    `json:"role" validate:"omitempty,max=50"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
}

type UpdateUserRequest struct {
	StoreID   *uuid.UUID `json:"store_id,omitempty"`
	Name      *string   `json:"name,omitempty" validate:"omitempty,max=255"`
	Username  *string   `json:"username,omitempty" validate:"omitempty,max=100"`
	Email     *string   `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Phone     *string   `json:"phone,omitempty" validate:"omitempty,max=20"`
	Password  *string   `json:"password,omitempty" validate:"omitempty,min=8"`
	Role      *string   `json:"role,omitempty" validate:"omitempty,max=50"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
}

type UpdateUserByAdminRequest struct {
	StoreID   *string `json:"store_id,omitempty"`
	Name      *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Username  *string `json:"username,omitempty" validate:"omitempty,max=100"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Password  *string `json:"password,omitempty" validate:"omitempty,min=8"`
	Role      *string `json:"role,omitempty" validate:"omitempty,max=50"`
	IsActive  *bool   `json:"is_active,omitempty"`
}