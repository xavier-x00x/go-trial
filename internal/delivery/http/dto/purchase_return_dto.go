package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request: Create Purchase Return
// ──────────────────────────────────────────────────────────────────────────────

type CreatePurchaseReturnRequest struct {
	PurchaseOrderID   uuid.UUID                       `json:"purchase_order_id" validate:"required"`
	PurchaseInvoiceID *uuid.UUID                      `json:"purchase_invoice_id"`
	SupplierID        uuid.UUID                       `json:"supplier_id" validate:"required"`
	WarehouseID       uuid.UUID                       `json:"warehouse_id" validate:"required"`
	ReturnDate        time.Time                       `json:"return_date" validate:"required"`
	Reason            string                          `json:"reason" validate:"required"`
	Notes             *string                         `json:"notes"`
	Items             []CreatePurchaseReturnItemInput `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseReturnItemInput struct {
	GoodsReceiptItemID *uuid.UUID      `json:"goods_receipt_item_id"`
	ProductID          uuid.UUID       `json:"product_id" validate:"required"`
	UOMID              uuid.UUID       `json:"uom_id" validate:"required"`
	QtyReturned        decimal.Decimal `json:"qty_returned" validate:"required"`
	UnitPrice          decimal.Decimal `json:"unit_price" validate:"required,min=0"`
	ReturnReason       string          `json:"return_reason" validate:"required"`
	Notes              *string         `json:"notes"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Request: Confirm & Cancel
// ──────────────────────────────────────────────────────────────────────────────

type ConfirmPurchaseReturnRequest struct {
	Notes *string `json:"notes"`
}

type CancelPurchaseReturnRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: List
// ──────────────────────────────────────────────────────────────────────────────

type PurchaseReturnListResponse struct {
	ID              uuid.UUID       `json:"id"`
	ReturnNumber    string          `json:"return_number"`
	PurchaseOrderID uuid.UUID       `json:"purchase_order_id"`
	PONumber        string          `json:"po_number"`
	SupplierID      uuid.UUID       `json:"supplier_id"`
	SupplierName    string          `json:"supplier_name"`
	ReturnDate      time.Time       `json:"return_date"`
	TotalAmount     decimal.Decimal `json:"total_amount"`
	Status          string          `json:"status"`
	Reason          string          `json:"reason"`
	CreatedAt       time.Time       `json:"created_at"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: Detail
// ──────────────────────────────────────────────────────────────────────────────

type PurchaseReturnDetailResponse struct {
	ID                uuid.UUID                    `json:"id"`
	ReturnNumber      string                       `json:"return_number"`
	PurchaseOrderID   uuid.UUID                    `json:"purchase_order_id"`
	PONumber          string                       `json:"po_number"`
	PurchaseInvoiceID *uuid.UUID                   `json:"purchase_invoice_id,omitempty"`
	InvoiceNumber     *string                      `json:"invoice_number,omitempty"`
	SupplierID        uuid.UUID                    `json:"supplier_id"`
	SupplierName      string                       `json:"supplier_name"`
	WarehouseID       uuid.UUID                    `json:"warehouse_id"`
	WarehouseName     string                       `json:"warehouse_name"`
	ReturnDate        time.Time                    `json:"return_date"`
	TotalAmount       decimal.Decimal              `json:"total_amount"`
	Status            string                       `json:"status"`
	CreatedByID       uuid.UUID                    `json:"created_by_id"`
	ConfirmedByID     *uuid.UUID                   `json:"confirmed_by_id,omitempty"`
	ConfirmedAt       *time.Time                   `json:"confirmed_at,omitempty"`
	Reason            string                       `json:"reason"`
	Notes             *string                      `json:"notes,omitempty"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
	Items             []PurchaseReturnItemResponse `json:"items"`
}

type PurchaseReturnItemResponse struct {
	ID                 uuid.UUID       `json:"id"`
	SeqNo              int             `json:"seq_no"`
	GoodsReceiptItemID *uuid.UUID      `json:"goods_receipt_item_id,omitempty"`
	ProductID          uuid.UUID       `json:"product_id"`
	ProductName        string          `json:"product_name"`
	ProductSKU         string          `json:"product_sku"`
	UOMID              uuid.UUID       `json:"uom_id"`
	UOMCode            string          `json:"uom_code"`
	QtyReturned        decimal.Decimal `json:"qty_returned"`
	UnitPrice          decimal.Decimal `json:"unit_price"`
	Subtotal           decimal.Decimal `json:"subtotal"`
	ReturnReason       string          `json:"return_reason"`
	Notes              *string         `json:"notes,omitempty"`
}
