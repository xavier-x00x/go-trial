package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateWarehouseRequest struct {
	StoreID  uuid.UUID `json:"store_id" validate:"required"`
	Code     string    `json:"code" validate:"required,max=20"`
	Name    string    `json:"name" validate:"required,max=100"`
	IsActive bool      `json:"is_active"`
}

type UpdateWarehouseRequest struct {
	StoreID  *uuid.UUID `json:"store_id,omitempty"`
	Code     *string   `json:"code,omitempty" validate:"omitempty,max=20"`
	Name     *string   `json:"name,omitempty" validate:"omitempty,max=100"`
	IsActive *bool     `json:"is_active,omitempty"`
}

type WarehouseResponse struct {
	ID        uuid.UUID      `json:"id"`
	StoreID   uuid.UUID      `json:"store_id"`
	Store     *StoreResponse `json:"store,omitempty"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}