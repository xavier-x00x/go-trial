package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request: Create Expense Voucher
// ──────────────────────────────────────────────────────────────────────────────

type CreateExpenseVoucherRequest struct {
	StoreID          uuid.UUID                       `json:"store_id" validate:"required"`
	SupplierID       *uuid.UUID                      `json:"supplier_id"`
	VendorName       *string                         `json:"vendor_name" validate:"omitempty,max=150"`
	VoucherDate      time.Time                       `json:"voucher_date" validate:"required"`
	PaymentType      string                          `json:"payment_type" validate:"required,oneof=CASH CREDIT"`
	PaymentAccountID *uuid.UUID                      `json:"payment_account_id" validate:"required_if=PaymentType CASH"`
	PayableAccountID *uuid.UUID                      `json:"payable_account_id" validate:"required_if=PaymentType CREDIT"`
	Description      string                          `json:"description" validate:"required,max=200"`
	Notes            *string                         `json:"notes"`
	Items            []CreateExpenseVoucherItemInput `json:"items" validate:"required,min=1,dive"`
}

type CreateExpenseVoucherItemInput struct {
	ExpenseAccountID uuid.UUID       `json:"expense_account_id" validate:"required"`
	Description      string          `json:"description" validate:"required,max=200"`
	Qty              decimal.Decimal `json:"qty" validate:"required"`
	UnitPrice        decimal.Decimal `json:"unit_price" validate:"required,min=0"`
	TaxPct           decimal.Decimal `json:"tax_pct" validate:"min=0"`
	Notes            *string         `json:"notes"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Request: Approve, Post, Void
// ──────────────────────────────────────────────────────────────────────────────

type ApproveExpenseVoucherRequest struct {
	Notes *string `json:"notes"`
}

type PostExpenseVoucherRequest struct {
	Notes *string `json:"notes"`
}

type VoidExpenseVoucherRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: List
// ──────────────────────────────────────────────────────────────────────────────

type ExpenseVoucherListResponse struct {
	ID            uuid.UUID       `json:"id"`
	VoucherNumber string          `json:"voucher_number"`
	StoreID       uuid.UUID       `json:"store_id"`
	StoreName     string          `json:"store_name"`
	VendorDisplay string          `json:"vendor_display"` // SupplierName atau VendorName
	VoucherDate   time.Time       `json:"voucher_date"`
	PaymentType   string          `json:"payment_type"`
	GrandTotal    decimal.Decimal `json:"grand_total"`
	Description   string          `json:"description"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: Detail
// ──────────────────────────────────────────────────────────────────────────────

type ExpenseVoucherDetailResponse struct {
	ID               uuid.UUID                    `json:"id"`
	VoucherNumber    string                       `json:"voucher_number"`
	StoreID          uuid.UUID                    `json:"store_id"`
	StoreName        string                       `json:"store_name"`
	SupplierID       *uuid.UUID                   `json:"supplier_id,omitempty"`
	SupplierName     *string                      `json:"supplier_name,omitempty"`
	VendorName       *string                      `json:"vendor_name,omitempty"`
	VoucherDate      time.Time                    `json:"voucher_date"`
	PaymentType      string                       `json:"payment_type"`
	PaymentAccountID *uuid.UUID                   `json:"payment_account_id,omitempty"`
	PayableAccountID *uuid.UUID                   `json:"payable_account_id,omitempty"`
	Subtotal         decimal.Decimal              `json:"subtotal"`
	TaxAmount        decimal.Decimal              `json:"tax_amount"`
	GrandTotal       decimal.Decimal              `json:"grand_total"`
	Description      string                       `json:"description"`
	Status           string                       `json:"status"`
	CreatedByID      uuid.UUID                    `json:"created_by_id"`
	ApprovedByID     *uuid.UUID                   `json:"approved_by_id,omitempty"`
	ApprovedAt       *time.Time                   `json:"approved_at,omitempty"`
	PostedByID       *uuid.UUID                   `json:"posted_by_id,omitempty"`
	PostedAt         *time.Time                   `json:"posted_at,omitempty"`
	Notes            *string                      `json:"notes,omitempty"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
	Items            []ExpenseVoucherItemResponse `json:"items"`
}

type ExpenseVoucherItemResponse struct {
	ID               uuid.UUID       `json:"id"`
	SeqNo            int             `json:"seq_no"`
	ExpenseAccountID uuid.UUID       `json:"expense_account_id"`
	AccountName      string          `json:"account_name"`
	Description      string          `json:"description"`
	Qty              decimal.Decimal `json:"qty"`
	UnitPrice        decimal.Decimal `json:"unit_price"`
	TaxPct           decimal.Decimal `json:"tax_pct"`
	TaxAmount        decimal.Decimal `json:"tax_amount"`
	Subtotal         decimal.Decimal `json:"subtotal"`
	Notes            *string         `json:"notes,omitempty"`
}
