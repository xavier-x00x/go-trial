package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductRequest struct {
	SKU           string          `json:"sku" validate:"required,max=50"`
	Barcode       *string         `json:"barcode,omitempty" validate:"omitempty,max=50"`
	Name          string          `json:"name" validate:"required,max=200"`
	CategoryID    *uuid.UUID      `json:"category_id,omitempty"`
	BaseUOMID     uuid.UUID       `json:"base_uom_id" validate:"required"`
	IsStockable   bool            `json:"is_stockable"`
	Length        decimal.Decimal `json:"length"`
	Width         decimal.Decimal `json:"width"`
	Height        decimal.Decimal `json:"height"`
	Weight        decimal.Decimal `json:"weight"`
	IsStackable   bool            `json:"is_stackable"`
	MaxStackLayer int             `json:"max_stack_layer"`
}

type UpdateProductRequest struct {
	SKU           *string          `json:"sku,omitempty" validate:"omitempty,max=50"`
	Name          *string          `json:"name,omitempty" validate:"omitempty,max=200"`
	Barcode       *string          `json:"barcode,omitempty" validate:"omitempty,max=50"`
	CategoryID    *uuid.UUID       `json:"category_id,omitempty"`
	BaseUOMID     *uuid.UUID       `json:"base_uom_id,omitempty"`
	IsStockable   *bool            `json:"is_stockable,omitempty"`
	Length        *decimal.Decimal `json:"length,omitempty"`
	Width         *decimal.Decimal `json:"width,omitempty"`
	Height        *decimal.Decimal `json:"height,omitempty"`
	Weight        *decimal.Decimal `json:"weight,omitempty"`
	IsStackable   *bool            `json:"is_stackable,omitempty"`
	MaxStackLayer *int             `json:"max_stack_layer,omitempty"`
}

type ProductResponse struct {
	ID            uuid.UUID         `json:"id"`
	SKU           string            `json:"sku"`
	Barcode       *string           `json:"barcode,omitempty"`
	Name          string            `json:"name"`
	CategoryID    *uuid.UUID        `json:"category_id,omitempty"`
	Category      *CategoryResponse `json:"category,omitempty"`
	BaseUOMID     uuid.UUID         `json:"base_uom_id"`
	BaseUOM       *UOMResponse      `json:"base_uom,omitempty"`
	IsStockable   bool              `json:"is_stockable"`
	Length        decimal.Decimal   `json:"length"`
	Width         decimal.Decimal   `json:"width"`
	Height        decimal.Decimal   `json:"height"`
	Weight        decimal.Decimal   `json:"weight"`
	IsStackable   bool              `json:"is_stackable"`
	MaxStackLayer int               `json:"max_stack_layer"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}
