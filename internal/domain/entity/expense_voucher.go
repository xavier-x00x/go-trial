package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	EVStatusDraft  = "DRAFT"
	EVStatusPosted = "POSTED"
	EVStatusVoided = "VOIDED"
)

type ExpenseVoucher struct {
	BaseModel

	VoucherNumber    string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"voucher_number"`
	VoucherDate      time.Time       `gorm:"not null" json:"voucher_date"`
	VendorName       string          `gorm:"type:varchar(255);not null" json:"vendor_name"`
	PaymentType      string          `gorm:"type:varchar(20);not null" json:"payment_type"` // CASH or CREDIT
	CreditAccountID  uuid.UUID       `gorm:"type:char(36);not null;index" json:"credit_account_id"`
	CreditAccount    ChartOfAccount  `gorm:"foreignKey:CreditAccountID" json:"credit_account,omitempty"`
	GrandTotal       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"grand_total"`
	Status           string          `gorm:"type:varchar(10);not null;default:'DRAFT';index" json:"status"`
	Notes            *string         `gorm:"type:text" json:"notes"`
	
	CreatedByID uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"`
	CreatedBy   User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	PostedByID  *uuid.UUID `gorm:"type:char(36);index" json:"posted_by_id"`
	PostedBy    *User      `gorm:"foreignKey:PostedByID" json:"posted_by,omitempty"`
	PostedAt    *time.Time `json:"posted_at"`

	Items []ExpenseVoucherItem `gorm:"foreignKey:ExpenseVoucherID" json:"items,omitempty"`
}

type ExpenseVoucherItem struct {
	BaseModel

	ExpenseVoucherID uuid.UUID      `gorm:"type:char(36);not null;index" json:"expense_voucher_id"`
	SeqNo            int            `gorm:"not null" json:"seq_no"`
	Description      string         `gorm:"type:varchar(255);not null" json:"description"`
	ExpenseAccountID uuid.UUID      `gorm:"type:char(36);not null;index" json:"expense_account_id"`
	ExpenseAccount   ChartOfAccount `gorm:"foreignKey:ExpenseAccountID" json:"expense_account,omitempty"`
	Amount           decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
}
