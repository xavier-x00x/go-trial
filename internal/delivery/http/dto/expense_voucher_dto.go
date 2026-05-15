package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateExpenseVoucherRequest struct {
	VoucherDate     time.Time                  `json:"voucher_date" validate:"required"`
	VendorName      string                     `json:"vendor_name" validate:"required"`
	PaymentType     string                     `json:"payment_type" validate:"required,oneof=CASH CREDIT"`
	CreditAccountID uuid.UUID                  `json:"credit_account_id" validate:"required"`
	Notes           *string                    `json:"notes"`
	Items           []CreateExpenseVoucherItem `json:"items" validate:"required,min=1,dive"`
}

type UpdateExpenseVoucherRequest struct {
	VoucherDate     time.Time                  `json:"voucher_date" validate:"required"`
	VendorName      string                     `json:"vendor_name" validate:"required"`
	PaymentType     string                     `json:"payment_type" validate:"required,oneof=CASH CREDIT"`
	CreditAccountID uuid.UUID                  `json:"credit_account_id" validate:"required"`
	Notes           *string                    `json:"notes"`
	Items           []CreateExpenseVoucherItem `json:"items" validate:"required,min=1,dive"`
}

type CreateExpenseVoucherItem struct {
	Description      string          `json:"description" validate:"required"`
	ExpenseAccountID uuid.UUID       `json:"expense_account_id" validate:"required"`
	Amount           decimal.Decimal `json:"amount" validate:"required,gt=0"`
}

type ExpenseVoucherDetailResponse struct {
	ID               uuid.UUID                    `json:"id"`
	VoucherNumber    string                       `json:"voucher_number"`
	VoucherDate      time.Time                    `json:"voucher_date"`
	VendorName       string                       `json:"vendor_name"`
	PaymentType      string                       `json:"payment_type"`
	CreditAccountID  uuid.UUID                    `json:"credit_account_id"`
	CreditAccountName string                      `json:"credit_account_name"`
	GrandTotal       decimal.Decimal              `json:"grand_total"`
	Status           string                       `json:"status"`
	Notes            *string                      `json:"notes,omitempty"`
	Items            []ExpenseVoucherItemResponse `json:"items"`
}

type ExpenseVoucherItemResponse struct {
	ID                 uuid.UUID       `json:"id"`
	SeqNo              int             `json:"seq_no"`
	Description        string          `json:"description"`
	ExpenseAccountID   uuid.UUID       `json:"expense_account_id"`
	ExpenseAccountName string          `json:"expense_account_name"`
	Amount             decimal.Decimal `json:"amount"`
}

type ExpenseVoucherListResponse struct {
	ID            uuid.UUID       `json:"id"`
	VoucherNumber string          `json:"voucher_number"`
	VoucherDate   time.Time       `json:"voucher_date"`
	VendorName    string          `json:"vendor_name"`
	PaymentType   string          `json:"payment_type"`
	GrandTotal    decimal.Decimal `json:"grand_total"`
	Status        string          `json:"status"`
}
