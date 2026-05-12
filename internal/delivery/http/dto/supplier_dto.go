package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateSupplierRequest struct {
	Code                          string           `json:"code" validate:"required,max=20"`
	Name                         string           `json:"name" validate:"required,max=150"`
	ContactPerson                *string          `json:"contact_person,omitempty" validate:"omitempty,max=100"`
	ContactPhone                 *string          `json:"contact_phone,omitempty" validate:"omitempty,max=20"`
	PhoneNumber                  *string          `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Email                        *string          `json:"email,omitempty" validate:"omitempty,email,max=100"`
	PreferredNotificationMethod string           `json:"preferred_notification_method,omitempty" validate:"omitempty,oneof=WHATSAPP EMAIL NONE"`
	Address                      *string          `json:"address,omitempty"`
	TaxRegNumber                 *string          `json:"tax_reg_number,omitempty" validate:"omitempty,max=50"`
	SupplierCategoryID            *uuid.UUID       `json:"supplier_category_id,omitempty"`
	IsPKP                        bool             `json:"is_pkp,omitempty"`
	PaymentTermDays              int              `json:"payment_term_days,omitempty" validate:"min=0"`
	PaymentMode                   string           `json:"payment_mode,omitempty" validate:"omitempty,oneof=CASH TRANSFER GIRO"`
	MinOrderAmount               json.Number      `json:"min_order_amount,omitempty"`
	BankName                     *string          `json:"bank_name,omitempty" validate:"omitempty,max=50"`
	BankAccount                  *string          `json:"bank_account,omitempty" validate:"omitempty,max=50"`
	BankAccountName              *string          `json:"bank_account_name,omitempty" validate:"omitempty,max=150"`
	APAccountID                  *uuid.UUID       `json:"ap_account_id,omitempty"`
}

func (r *CreateSupplierRequest) GetMinOrderAmount() decimal.Decimal {
	if r.MinOrderAmount == "" {
		return decimal.Zero
	}
	d, _ := decimal.NewFromString(string(r.MinOrderAmount))
	return d
}

type UpdateSupplierRequest struct {
	Code                          *string          `json:"code,omitempty" validate:"omitempty,max=20"`
	Name                         *string          `json:"name,omitempty" validate:"omitempty,max=150"`
	ContactPerson                *string          `json:"contact_person,omitempty" validate:"omitempty,max=100"`
	ContactPhone                 *string          `json:"contact_phone,omitempty" validate:"omitempty,max=20"`
	PhoneNumber                 *string          `json:"phone_number,omitempty" validate:"omitempty,max=20"`
	Email                       *string          `json:"email,omitempty" validate:"omitempty,email,max=100"`
	PreferredNotificationMethod *string          `json:"preferred_notification_method,omitempty" validate:"omitempty,oneof=WHATSAPP EMAIL NONE"`
	Address                     *string          `json:"address,omitempty"`
	TaxRegNumber                *string          `json:"tax_reg_number,omitempty" validate:"omitempty,max=50"`
	SupplierCategoryID          *uuid.UUID      `json:"supplier_category_id,omitempty"`
	IsPKP                       *bool           `json:"is_pkp,omitempty"`
	PaymentTermDays             *int            `json:"payment_term_days,omitempty" validate:"min=0"`
	PaymentMode                 *string          `json:"payment_mode,omitempty" validate:"omitempty,oneof=CASH TRANSFER GIRO"`
	MinOrderAmount              *json.Number     `json:"min_order_amount,omitempty"`
	BankName                    *string          `json:"bank_name,omitempty" validate:"omitempty,max=50"`
	BankAccount                 *string          `json:"bank_account,omitempty" validate:"omitempty,max=50"`
	BankAccountName             *string          `json:"bank_account_name,omitempty" validate:"omitempty,max=150"`
	IsActive                    *bool            `json:"is_active,omitempty"`
	APAccountID                 *uuid.UUID      `json:"ap_account_id,omitempty"`
}

func (r *UpdateSupplierRequest) GetMinOrderAmount() decimal.Decimal {
	if r.MinOrderAmount == nil || *r.MinOrderAmount == "" {
		return decimal.Zero
	}
	d, _ := decimal.NewFromString(string(*r.MinOrderAmount))
	return d
}

func (r *UpdateSupplierRequest) Get() decimal.Decimal {
	if r.MinOrderAmount == nil || *r.MinOrderAmount == "" {
		return decimal.Zero
	}
	d, _ := decimal.NewFromString(string(*r.MinOrderAmount))
	return d
}

type SupplierResponse struct {
	ID                          uuid.UUID                `json:"id"`
	Code                       string                  `json:"code"`
	Name                       string                  `json:"name"`
	ContactPerson              *string                  `json:"contact_person,omitempty"`
	ContactPhone               *string                  `json:"contact_phone,omitempty"`
	PhoneNumber               *string                  `json:"phone_number,omitempty"`
	Email                     *string                  `json:"email,omitempty"`
	PreferredNotificationMethod string                 `json:"preferred_notification_method"`
	Address                   *string                  `json:"address,omitempty"`
	TaxRegNumber               *string                  `json:"tax_reg_number,omitempty"`
	SupplierCategoryID       *uuid.UUID               `json:"supplier_category_id,omitempty"`
	SupplierCategory         *SupplierCategoryResponse `json:"supplier_category,omitempty"`
	IsPKP                     bool                     `json:"is_pkp"`
	PaymentTermDays           int                      `json:"payment_term_days"`
	PaymentMode               string                   `json:"payment_mode"`
	MinOrderAmount             decimal.Decimal          `json:"min_order_amount"`
	BankName                  *string                  `json:"bank_name,omitempty"`
	BankAccount              *string                  `json:"bank_account,omitempty"`
	BankAccountName           *string                  `json:"bank_account_name,omitempty"`
	IsActive                  bool                     `json:"is_active"`
	APAccountID              *uuid.UUID                `json:"ap_account_id,omitempty"`
	APAccount                *ChartOfAccountResponse  `json:"ap_account,omitempty"`
	CreatedAt                 time.Time                `json:"created_at"`
	UpdatedAt                 time.Time                `json:"updated_at"`
}