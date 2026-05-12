package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductSupplierRequest struct {
	ProductID           uuid.UUID       `json:"product_id" validate:"required"`
	SupplierID          uuid.UUID       `json:"supplier_id" validate:"required"`
	StoreID             *uuid.UUID      `json:"store_id"`
	SupplierSKU         *string         `json:"supplier_sku" validate:"omitempty,max=50"`
	IsPrimary           bool            `json:"is_primary"`
	IsConsignment       bool            `json:"is_consignment"`
	IsReturnable        bool            `json:"is_returnable"`
	DefaultLeadTimeDays int             `json:"default_lead_time_days" validate:"min=0"`
	OfferedPrice        decimal.Decimal `json:"offered_price" validate:"min=0"`
	MinOrderQty         decimal.Decimal `json:"min_order_qty" validate:"min=0"`
}

type CreateProductSupplierItemRequest struct {
	ProductID           uuid.UUID       `json:"product_id" validate:"required"`
	SupplierID          uuid.UUID       `json:"supplier_id" validate:"required"`
	StoreID             *uuid.UUID      `json:"store_id"`
	SupplierSKU         *string         `json:"supplier_sku" validate:"omitempty,max=50"`
	IsPrimary           bool            `json:"is_primary"`
	IsConsignment       bool            `json:"is_consignment"`
	IsReturnable        bool            `json:"is_returnable"`
	DefaultLeadTimeDays int             `json:"default_lead_time_days" validate:"min=0"`
	OfferedPrice        decimal.Decimal `json:"offered_price" validate:"min=0"`
	MinOrderQty         decimal.Decimal `json:"min_order_qty" validate:"min=0"`
}

type UpdateProductSupplierRequest struct {
	StoreID             *uuid.UUID      `json:"store_id"`
	SupplierSKU         *string         `json:"supplier_sku" validate:"omitempty,max=50"`
	IsPrimary           bool            `json:"is_primary"`
	IsConsignment       bool            `json:"is_consignment"`
	IsReturnable        bool            `json:"is_returnable"`
	DefaultLeadTimeDays int             `json:"default_lead_time_days" validate:"min=0"`
	OfferedPrice        decimal.Decimal `json:"offered_price" validate:"min=0"`
	MinOrderQty         decimal.Decimal `json:"min_order_qty" validate:"min=0"`
}

type ProductSupplierResponse struct {
	ID                  uuid.UUID       `json:"id"`
	ProductID           uuid.UUID       `json:"product_id"`
	SupplierID          uuid.UUID       `json:"supplier_id"`
	StoreID             *uuid.UUID      `json:"store_id,omitempty"`
	SupplierSKU         *string         `json:"supplier_sku,omitempty"`
	IsPrimary           bool            `json:"is_primary"`
	IsConsignment       bool            `json:"is_consignment"`
	IsReturnable        bool            `json:"is_returnable"`
	DefaultLeadTimeDays int             `json:"default_lead_time_days"`
	OfferedPrice        decimal.Decimal `json:"offered_price"`
	MinOrderQty         decimal.Decimal `json:"min_order_qty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}
