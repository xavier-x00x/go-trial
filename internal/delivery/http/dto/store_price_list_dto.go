package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateStorePriceListRequest struct {
	StoreID     uuid.UUID `json:"store_id" validate:"required"`
	PriceListID uuid.UUID `json:"price_list_id" validate:"required"`
	Priority    int       `json:"priority"`
	IsActive    bool      `json:"is_active"`
}

type UpdateStorePriceListRequest struct {
	Priority int  `json:"priority"`
	IsActive bool `json:"is_active"`
}

type StorePriceListResponse struct {
	ID          uuid.UUID `json:"id"`
	StoreID     uuid.UUID `json:"store_id"`
	PriceListID uuid.UUID `json:"price_list_id"`
	Priority    int       `json:"priority"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
