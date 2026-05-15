package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreatePurchaseReturnRequest struct {
	PurchaseInvoiceID uuid.UUID                   `json:"purchase_invoice_id" validate:"required"`
	ReturnDate        time.Time                   `json:"return_date" validate:"required"`
	Notes             *string                     `json:"notes"`
	Items             []CreatePurchaseReturnItem `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseReturnItem struct {
	PurchaseInvoiceItemID uuid.UUID       `json:"purchase_invoice_item_id" validate:"required"`
	ProductID           uuid.UUID       `json:"product_id" validate:"required"`
	UOMID               uuid.UUID       `json:"uom_id" validate:"required"`
	QtyReturn           decimal.Decimal `json:"qty_return" validate:"required,gt=0"`
	Notes               *string         `json:"notes"`
}

type PurchaseReturnDetailResponse struct {
	ID                uuid.UUID                  `json:"id"`
	ReturnNumber      string                     `json:"return_number"`
	ReturnDate        time.Time                  `json:"return_date"`
	PurchaseInvoiceID uuid.UUID                  `json:"purchase_invoice_id"`
	InvoiceNumber     string                     `json:"invoice_number"`
	SupplierName      string                     `json:"supplier_name"`
	StoreName         string                     `json:"store_name"`
	WarehouseName     string                     `json:"warehouse_name"`
	Subtotal          decimal.Decimal            `json:"subtotal"`
	DiscountAmount    decimal.Decimal            `json:"discount_amount"`
	TaxAmount         decimal.Decimal            `json:"tax_amount"`
	GrandTotal        decimal.Decimal            `json:"grand_total"`
	RemainingAmount   decimal.Decimal            `json:"remaining_amount"`
	Status            string                     `json:"status"`
	Notes             *string                    `json:"notes,omitempty"`
	Items             []PurchaseReturnItemResponse `json:"items"`
}

type PurchaseReturnItemResponse struct {
	ID              uuid.UUID       `json:"id"`
	SeqNo           int             `json:"seq_no"`
	ProductID       uuid.UUID       `json:"product_id"`
	ProductName     string          `json:"product_name"`
	UOMCode         string          `json:"uom_code"`
	QtyReturn       decimal.Decimal `json:"qty_return"`
	UnitPrice       decimal.Decimal `json:"unit_price"`
	DiscountAmount  decimal.Decimal `json:"discount_amount"`
	TaxAmount       decimal.Decimal `json:"tax_amount"`
	Subtotal        decimal.Decimal `json:"subtotal"`
}

type PurchaseReturnListResponse struct {
	ID            uuid.UUID `json:"id"`
	ReturnNumber  string    `json:"return_number"`
	ReturnDate    time.Time `json:"return_date"`
	InvoiceNumber string    `json:"invoice_number"`
	SupplierName  string    `json:"supplier_name"`
	GrandTotal    decimal.Decimal `json:"grand_total"`
	Status        string    `json:"status"`
}
