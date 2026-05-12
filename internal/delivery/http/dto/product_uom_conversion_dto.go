package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductUOMConversionRequest struct {
	ProductID      uuid.UUID        `json:"product_id" validate:"required"`
	UOMID          uuid.UUID        `json:"uom_id" validate:"required"`
	ConversionRate decimal.Decimal `json:"conversion_rate" validate:"required"`
	Barcode        *string         `json:"barcode" validate:"omitempty,max=50"`
	Length        decimal.Decimal `json:"length"`
	Width         decimal.Decimal `json:"width"`
	Height        decimal.Decimal `json:"height"`
	Weight        decimal.Decimal `json:"weight"`
	IsStackable   bool            `json:"is_stackable"`
	MaxStackLayer int             `json:"max_stack_layer"`
}

type UpdateProductUOMConversionRequest struct {
	UOMID          *uuid.UUID      `json:"uom_id"`
	ConversionRate *decimal.Decimal `json:"conversion_rate"`
	Barcode        *string         `json:"barcode" validate:"omitempty,max=50"`
	Length        *decimal.Decimal `json:"length"`
	Width         *decimal.Decimal `json:"width"`
	Height        *decimal.Decimal `json:"height"`
	Weight        *decimal.Decimal `json:"weight"`
	IsStackable    *bool           `json:"is_stackable"`
	MaxStackLayer *int            `json:"max_stack_layer"`
}

type ProductUOMConversionResponse struct {
	ID             uuid.UUID        `json:"id"`
	ProductID     uuid.UUID        `json:"product_id"`
	Product       *ProductResponse  `json:"product,omitempty"`
	UOMID         uuid.UUID        `json:"uom_id"`
	UOM           *UOMResponse    `json:"uom,omitempty"`
	ConversionRate decimal.Decimal `json:"conversion_rate"`
	Barcode        *string        `json:"barcode,omitempty"`
	Length        decimal.Decimal `json:"length"`
	Width         decimal.Decimal `json:"width"`
	Height        decimal.Decimal `json:"height"`
	Weight        decimal.Decimal `json:"weight"`
	IsStackable   bool            `json:"is_stackable"`
	MaxStackLayer int             `json:"max_stack_layer"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}