package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateCustomerRequest struct {
	Code         *string         `json:"code,omitempty" validate:"omitempty,max=20"`
	Name         *string         `json:"name,omitempty" validate:"omitempty,max=150"`
	PhoneNumber  *string         `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Email        *string         `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Address      *string         `json:"address,omitempty"`
	IsActive     *bool           `json:"is_active,omitempty"`
	PointBalance *decimal.Decimal `json:"point_balance,omitempty"`
	CreditLimit  *decimal.Decimal `json:"credit_limit,omitempty"`
	ARAccountID  *uuid.UUID     `json:"ar_account_id,omitempty"`
}

type CreateCustomerRequest struct {
	Code          string          `json:"code" validate:"required,max=20"`
	Name          string          `json:"name" validate:"required,max=150"`
	PhoneNumber   *string         `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Email         *string         `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Address       *string         `json:"address,omitempty"`
	PointBalance  decimal.Decimal  `json:"point_balance,omitempty"`
	CreditLimit   decimal.Decimal  `json:"credit_limit,omitempty"`
	ARAccountID  *uuid.UUID     `json:"ar_account_id,omitempty"`
}

type CustomerResponse struct {
	ID            uuid.UUID         `json:"id"`
	Code          string           `json:"code"`
	Name          string          `json:"name"`
	PhoneNumber   *string         `json:"phone_number,omitempty"`
	Email         *string        `json:"email,omitempty"`
	Address       *string        `json:"address,omitempty"`
	IsActive      bool            `json:"is_active"`
	PointBalance decimal.Decimal `json:"point_balance"`
	CreditLimit  decimal.Decimal `json:"credit_limit"`
	ARAccountID  *uuid.UUID     `json:"ar_account_id,omitempty"`
	ARAccount    *ChartOfAccountResponse `json:"ar_account,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}