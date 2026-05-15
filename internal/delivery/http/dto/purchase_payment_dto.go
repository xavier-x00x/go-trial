package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request: Create Purchase Payment
// ──────────────────────────────────────────────────────────────────────────────

type CreatePurchasePaymentRequest struct {
	SupplierID       uuid.UUID                        `json:"supplier_id" validate:"required"`
	PaymentAccountID uuid.UUID                        `json:"payment_account_id" validate:"required"`
	APAccountID      uuid.UUID                        `json:"ap_account_id" validate:"required"`
	PaymentDate      time.Time                        `json:"payment_date" validate:"required"`
	PaymentMode      string                           `json:"payment_mode" validate:"required,oneof=CASH TRANSFER GIRO"`
	ReferenceNo      *string                          `json:"reference_no" validate:"omitempty,max=50"`
	GiroNumber       *string                          `json:"giro_number" validate:"omitempty,max=50"`
	GiroDueDate      *time.Time                       `json:"giro_due_date"`
	AdminFeeAmount    decimal.Decimal                  `json:"admin_fee_amount"`
	AdminFeeAccountID *uuid.UUID                       `json:"admin_fee_account_id"`
	DiscountAmount    decimal.Decimal                  `json:"discount_amount"`
	DiscountAccountID *uuid.UUID                       `json:"discount_account_id"`
	WHTAmount         decimal.Decimal                  `json:"wht_amount"`
	WHTAccountID      *uuid.UUID                       `json:"wht_account_id"`
	Notes            *string                          `json:"notes"`
	Items            []CreatePurchasePaymentItemInput `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchasePaymentItemInput struct {
	PurchaseInvoiceID *uuid.UUID      `json:"purchase_invoice_id"`
	PurchaseReturnID  *uuid.UUID      `json:"purchase_return_id"`
	PaidAmount        decimal.Decimal `json:"paid_amount" validate:"required"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: Post & Void
// ──────────────────────────────────────────────────────────────────────────────

type PostPurchasePaymentRequest struct {
	Notes *string `json:"notes"`
}

type VoidPurchasePaymentRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: List
// ──────────────────────────────────────────────────────────────────────────────

type PurchasePaymentListResponse struct {
	ID            uuid.UUID       `json:"id"`
	PaymentNumber string          `json:"payment_number"`
	SupplierID    uuid.UUID       `json:"supplier_id"`
	SupplierName  string          `json:"supplier_name"`
	PaymentDate   time.Time       `json:"payment_date"`
	PaymentMode   string          `json:"payment_mode"`
	TotalAmount   decimal.Decimal `json:"total_amount"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: Detail
// ──────────────────────────────────────────────────────────────────────────────

type PurchasePaymentDetailResponse struct {
	ID               uuid.UUID                     `json:"id"`
	PaymentNumber    string                        `json:"payment_number"`
	ReferenceNo      *string                       `json:"reference_no,omitempty"`
	SupplierID       uuid.UUID                     `json:"supplier_id"`
	SupplierName     string                        `json:"supplier_name"`
	PaymentAccountID uuid.UUID                     `json:"payment_account_id"`
	APAccountID      uuid.UUID                     `json:"ap_account_id"`
	PaymentDate      time.Time                     `json:"payment_date"`
	PaymentMode      string                        `json:"payment_mode"`
	GiroNumber       *string                       `json:"giro_number,omitempty"`
	GiroDueDate      *time.Time                    `json:"giro_due_date,omitempty"`
	TotalAmount      decimal.Decimal               `json:"total_amount"`
	AdminFeeAmount    decimal.Decimal               `json:"admin_fee_amount"`
	AdminFeeAccountID *uuid.UUID                    `json:"admin_fee_account_id,omitempty"`
	DiscountAmount    decimal.Decimal               `json:"discount_amount"`
	DiscountAccountID *uuid.UUID                    `json:"discount_account_id,omitempty"`
	WHTAmount         decimal.Decimal               `json:"wht_amount"`
	WHTAccountID      *uuid.UUID                    `json:"wht_account_id,omitempty"`
	Status           string                        `json:"status"`
	CreatedByID      uuid.UUID                     `json:"created_by_id"`
	PostedByID       *uuid.UUID                    `json:"posted_by_id,omitempty"`
	PostedAt         *time.Time                    `json:"posted_at,omitempty"`
	Notes            *string                       `json:"notes,omitempty"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
	Items            []PurchasePaymentItemResponse `json:"items"`
}

type PurchasePaymentItemResponse struct {
	ID                uuid.UUID       `json:"id"`
	SeqNo             int             `json:"seq_no"`
	PurchaseInvoiceID *uuid.UUID      `json:"purchase_invoice_id,omitempty"`
	InvoiceNumber     string          `json:"invoice_number,omitempty"`
	PurchaseReturnID  *uuid.UUID      `json:"purchase_return_id,omitempty"`
	ReturnNumber      string          `json:"return_number,omitempty"`
	DocumentAmount    decimal.Decimal `json:"document_amount"`
	PaidAmount        decimal.Decimal `json:"paid_amount"`
}
