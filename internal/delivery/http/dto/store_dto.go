package dto

import "time"

// CreateStoreRequest is the payload for creating a new store.
type CreateStoreRequest struct {
	Code       string  `json:"code" validate:"required,min=1,max=20"`
	Name       string  `json:"name" validate:"required,min=1,max=150"`
	NPWP       *string `json:"npwp,omitempty" validate:"omitempty,max=50"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty" validate:"omitempty,max=100"`
	Province   *string `json:"province,omitempty" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code,omitempty" validate:"omitempty,max=10"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Email      *string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	IsMain     bool    `json:"is_main"`
}

// UpdateStoreRequest is the payload for updating a store (partial update).
type UpdateStoreRequest struct {
	Code       *string `json:"code,omitempty" validate:"omitempty,min=1,max=20"`
	Name       *string `json:"name,omitempty" validate:"omitempty,min=1,max=150"`
	NPWP       *string `json:"npwp,omitempty" validate:"omitempty,max=50"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty" validate:"omitempty,max=100"`
	Province   *string `json:"province,omitempty" validate:"omitempty,max=100"`
	PostalCode *string `json:"postal_code,omitempty" validate:"omitempty,max=10"`
	Phone      *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Email      *string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	IsMain     *bool   `json:"is_main,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// StoreResponse is the public representation of a store.
type StoreResponse struct {
	ID         string    `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	NPWP       *string   `json:"npwp,omitempty"`
	Address    *string   `json:"address,omitempty"`
	City       *string   `json:"city,omitempty"`
	Province   *string   `json:"province,omitempty"`
	PostalCode *string   `json:"postal_code,omitempty"`
	Phone      *string   `json:"phone,omitempty"`
	Email      *string   `json:"email,omitempty"`
	IsMain     bool      `json:"is_main"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
