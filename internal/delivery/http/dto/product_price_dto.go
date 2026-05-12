package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductPriceRequest struct {
	PriceListID uuid.UUID       `json:"price_list_id" validate:"required"`
	ProductID   uuid.UUID       `json:"product_id" validate:"required"`
	UOMID       uuid.UUID       `json:"uom_id" validate:"required"`
	MarkupPct   decimal.Decimal `json:"markup_pct"`
	SellPrice   decimal.Decimal `json:"sell_price" validate:"required"`
	DiscountPct decimal.Decimal `json:"discount_pct"`
}

type UpdateProductPriceRequest struct {
	PriceListID *uuid.UUID       `json:"price_list_id"`
	ProductID   *uuid.UUID       `json:"product_id"`
	UOMID       *uuid.UUID       `json:"uom_id"`
	MarkupPct   *decimal.Decimal `json:"markup_pct"`
	SellPrice   *decimal.Decimal `json:"sell_price"`
	DiscountPct *decimal.Decimal `json:"discount_pct"`
}

type ProductPriceResponse struct {
	ID          uuid.UUID          `json:"id"`
	PriceListID uuid.UUID          `json:"price_list_id"`
	PriceList   *PriceListResponse `json:"price_list,omitempty"`
	ProductID   uuid.UUID          `json:"product_id"`
	Product     *ProductResponse   `json:"product,omitempty"`
	UOMID       uuid.UUID          `json:"uom_id"`
	UOM         *UOMResponse       `json:"uom,omitempty"`
	MarkupPct   decimal.Decimal    `json:"markup_pct"`
	SellPrice   decimal.Decimal    `json:"sell_price"`
	DiscountPct decimal.Decimal    `json:"discount_pct"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
