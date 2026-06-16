package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePriceListRequest struct {
	Code         string `json:"code" validate:"required,max=20"`
	Name         string `json:"name" validate:"required,max=100"`
	CurrencyCode string `json:"currency_code" validate:"omitempty,max=3"`
}

type UpdatePriceListRequest struct {
	Code         *string `json:"code,omitempty" validate:"omitempty,max=20"`
	Name         *string `json:"name,omitempty" validate:"omitempty,max=100"`
	CurrencyCode *string `json:"currency_code,omitempty" validate:"omitempty,max=3"`
	IsActive     *bool   `json:"is_active"`
}

type PriceListResponse struct {
	ID           uuid.UUID      `json:"id"`
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	CurrencyCode string         `json:"currency_code"`
	StoreID      *uuid.UUID     `json:"store_id,omitempty"`
	Store        *StoreResponse `json:"store,omitempty"`
	IsActive     bool           `json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}