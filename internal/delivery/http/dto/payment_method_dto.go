package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreatePaymentMethodRequest struct {
	Code             string          `json:"code" validate:"required,max=20"`
	Name             string          `json:"name" validate:"required,max=50"`
	MdrPercentage    decimal.Decimal `json:"mdr_percentage"`
	DepositAccountID *uuid.UUID      `json:"deposit_account_id,omitempty"`
	ExpenseAccountID *uuid.UUID      `json:"expense_account_id,omitempty"`
}

func (r *CreatePaymentMethodRequest) GetMdrPercentage() decimal.Decimal {
	return r.MdrPercentage
}

type UpdatePaymentMethodRequest struct {
	Name             *string         `json:"name,omitempty" validate:"omitempty,max=50"`
	MdrPercentage   *decimal.Decimal `json:"mdr_percentage"`
	DepositAccountID *uuid.UUID     `json:"deposit_account_id,omitempty"`
	ExpenseAccountID *uuid.UUID    `json:"expense_account_id,omitempty"`
	IsActive        *bool          `json:"is_active"`
}

func (r *UpdatePaymentMethodRequest) GetMdrPercentage() decimal.Decimal {
	if r == nil || r.MdrPercentage == nil {
		return decimal.Zero
	}
	return *r.MdrPercentage
}

type PaymentMethodResponse struct {
	ID               uuid.UUID                `json:"id"`
	Code             string                  `json:"code"`
	Name             string                  `json:"name"`
	MdrPercentage    decimal.Decimal       `json:"mdr_percentage"`
	DepositAccountID *uuid.UUID            `json:"deposit_account_id,omitempty"`
	DepositAccount   *ChartOfAccountResponse `json:"deposit_account,omitempty"`
	ExpenseAccountID *uuid.UUID        `json:"expense_account_id,omitempty"`
	ExpenseAccount  *ChartOfAccountResponse `json:"expense_account,omitempty"`
	IsActive        bool                  `json:"is_active"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}