package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request: Create Purchase Invoice
// ──────────────────────────────────────────────────────────────────────────────

type CreatePurchaseInvoiceRequest struct {
	PurchaseOrderID       uuid.UUID                        `json:"purchase_order_id" validate:"required"`
	SupplierID            uuid.UUID                        `json:"supplier_id" validate:"required"`
	SupplierInvoiceNumber string                           `json:"supplier_invoice_number" validate:"required,max=50"`
	APAccountID           uuid.UUID                        `json:"ap_account_id" validate:"required"`
	InvoiceDate           time.Time                        `json:"invoice_date" validate:"required"`
	ReceivedDate          time.Time                        `json:"received_date" validate:"required"`
	PaymentTermDays       int                              `json:"payment_term_days" validate:"min=0"`
	PaymentMode           string                           `json:"payment_mode" validate:"omitempty,oneof=CASH TRANSFER GIRO"`
	DiscountAmount        decimal.Decimal                  `json:"discount_amount" validate:"min=0"`
	Notes                 *string                          `json:"notes"`
	Items                 []CreatePurchaseInvoiceItemInput `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseInvoiceItemInput struct {
	PurchaseOrderItemID *uuid.UUID      `json:"purchase_order_item_id"`
	GoodsReceiptItemID  *uuid.UUID      `json:"goods_receipt_item_id"`
	ProductID           uuid.UUID       `json:"product_id" validate:"required"`
	UOMID               uuid.UUID       `json:"uom_id" validate:"required"`
	QtyInvoiced         decimal.Decimal `json:"qty_invoiced" validate:"required"`
	UnitPrice           decimal.Decimal `json:"unit_price" validate:"required,min=0"`
	Discount1Pct        decimal.Decimal `json:"discount_1_pct" validate:"min=0"`
	Discount2Pct        decimal.Decimal `json:"discount_2_pct" validate:"min=0"`
	Discount3Pct        decimal.Decimal `json:"discount_3_pct" validate:"min=0"`
	DiscountAmount      decimal.Decimal `json:"discount_amount" validate:"min=0"`
	TaxPct              decimal.Decimal `json:"tax_pct" validate:"min=0"`
	Notes               *string         `json:"notes"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Request: Verify & Post
// ──────────────────────────────────────────────────────────────────────────────

// VerifyPurchaseInvoiceRequest digunakan oleh staf Finance setelah mencocokkan
// faktur dengan PO dan GR (3-Way Match).
type VerifyPurchaseInvoiceRequest struct {
	Notes *string `json:"notes"`
}

// PostPurchaseInvoiceRequest digunakan oleh Manajer Akuntansi untuk
// memposting jurnal hutang ke buku besar.
type PostPurchaseInvoiceRequest struct {
	Notes *string `json:"notes"`
}

type CancelPurchaseInvoiceRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: List (Header Only)
// ──────────────────────────────────────────────────────────────────────────────

type PurchaseInvoiceListResponse struct {
	ID                    uuid.UUID       `json:"id"`
	InvoiceNumber         string          `json:"invoice_number"`
	SupplierInvoiceNumber string          `json:"supplier_invoice_number"`
	PurchaseOrderID       uuid.UUID       `json:"purchase_order_id"`
	PONumber              string          `json:"po_number"`
	SupplierID            uuid.UUID       `json:"supplier_id"`
	SupplierName          string          `json:"supplier_name"`
	InvoiceDate           time.Time       `json:"invoice_date"`
	DueDate               time.Time       `json:"due_date"`
	GrandTotal            decimal.Decimal `json:"grand_total"`
	PaidAmount            decimal.Decimal `json:"paid_amount"`
	RemainingAmount       decimal.Decimal `json:"remaining_amount"`
	Status                string          `json:"status"`
	CreatedAt             time.Time       `json:"created_at"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: Detail (Header + Items)
// ──────────────────────────────────────────────────────────────────────────────

type PurchaseInvoiceDetailResponse struct {
	ID                    uuid.UUID                     `json:"id"`
	InvoiceNumber         string                        `json:"invoice_number"`
	SupplierInvoiceNumber string                        `json:"supplier_invoice_number"`
	PurchaseOrderID       uuid.UUID                     `json:"purchase_order_id"`
	PONumber              string                        `json:"po_number"`
	SupplierID            uuid.UUID                     `json:"supplier_id"`
	SupplierName          string                        `json:"supplier_name"`
	APAccountID           uuid.UUID                     `json:"ap_account_id"`
	InvoiceDate           time.Time                     `json:"invoice_date"`
	ReceivedDate          time.Time                     `json:"received_date"`
	DueDate               time.Time                     `json:"due_date"`
	PaymentTermDays       int                           `json:"payment_term_days"`
	PaymentMode           string                        `json:"payment_mode"`
	Subtotal              decimal.Decimal               `json:"subtotal"`
	DiscountAmount        decimal.Decimal               `json:"discount_amount"`
	TaxAmount             decimal.Decimal               `json:"tax_amount"`
	FreightAmount         decimal.Decimal               `json:"freight_amount"`
	OtherCostAmount       decimal.Decimal               `json:"other_cost_amount"`
	GrandTotal            decimal.Decimal               `json:"grand_total"`
	IsTaxInclusive        bool                          `json:"is_tax_inclusive"`
	PaidAmount            decimal.Decimal               `json:"paid_amount"`
	RemainingAmount       decimal.Decimal               `json:"remaining_amount"`
	Status                string                        `json:"status"`
	VerifiedByID          *uuid.UUID                    `json:"verified_by_id,omitempty"`
	VerifiedAt            *time.Time                    `json:"verified_at,omitempty"`
	PostedByID            *uuid.UUID                    `json:"posted_by_id,omitempty"`
	PostedAt              *time.Time                    `json:"posted_at,omitempty"`
	CreatedByID           uuid.UUID                     `json:"created_by_id"`
	Notes                 *string                       `json:"notes,omitempty"`
	CreatedAt             time.Time                     `json:"created_at"`
	UpdatedAt             time.Time                     `json:"updated_at"`
	Items                 []PurchaseInvoiceItemResponse `json:"items"`
}

type PurchaseInvoiceItemResponse struct {
	ID                  uuid.UUID       `json:"id"`
	SeqNo               int             `json:"seq_no"`
	PurchaseOrderItemID *uuid.UUID      `json:"purchase_order_item_id,omitempty"`
	GoodsReceiptItemID  *uuid.UUID      `json:"goods_receipt_item_id,omitempty"`
	ProductID           uuid.UUID       `json:"product_id"`
	ProductName         string          `json:"product_name"`
	ProductSKU          string          `json:"product_sku"`
	UOMID               uuid.UUID       `json:"uom_id"`
	UOMCode             string          `json:"uom_code"`
	QtyInvoiced         decimal.Decimal `json:"qty_invoiced"`
	UnitPrice           decimal.Decimal `json:"unit_price"`
	Discount1Pct        decimal.Decimal `json:"discount_1_pct"`
	Discount2Pct        decimal.Decimal `json:"discount_2_pct"`
	Discount3Pct        decimal.Decimal `json:"discount_3_pct"`
	DiscountAmount      decimal.Decimal `json:"discount_amount"`
	TotalDiscountAmount decimal.Decimal `json:"total_discount_amount"`
	TaxPct              decimal.Decimal `json:"tax_pct"`
	TaxAmount           decimal.Decimal `json:"tax_amount"`
	LandedCostAmount    decimal.Decimal `json:"landed_cost_amount"`
	Subtotal            decimal.Decimal `json:"subtotal"`
	NetUnitPrice        decimal.Decimal `json:"net_unit_price"`
	Notes               *string         `json:"notes,omitempty"`
}
