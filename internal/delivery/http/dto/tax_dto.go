package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateTaxRequest struct {
	Name           string          `json:"name" validate:"required,max=50"`
	RatePercentage decimal.Decimal `json:"rate_percentage" validate:"required"`
	TaxAccountID   *uuid.UUID     `json:"tax_account_id,omitempty"`
}

func (r *CreateTaxRequest) GetRatePercentage() decimal.Decimal {
	return r.RatePercentage
}

type UpdateTaxRequest struct {
	Name           *string         `json:"name,omitempty" validate:"omitempty,max=50"`
	RatePercentage *decimal.Decimal `json:"rate_percentage"`
	TaxAccountID   *uuid.UUID     `json:"tax_account_id,omitempty"`
}

func (r *UpdateTaxRequest) GetRatePercentage() decimal.Decimal {
	if r.RatePercentage == nil {
		return decimal.Zero
	}
	return *r.RatePercentage
}

type TaxResponse struct {
	ID             uuid.UUID                `json:"id"`
	Name           string                  `json:"name"`
	RatePercentage decimal.Decimal       `json:"rate_percentage"`
	TaxAccountID   *uuid.UUID           `json:"tax_account_id,omitempty"`
	TaxAccount     *ChartOfAccountResponse `json:"tax_account,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}