package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaymentMethod merepresentasikan metode pembayaran elektronik/tunai yang tersedia di Kasir
type PaymentMethod struct {
	BaseModel
	Code             string          `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`            // Identifier internal metode (CASH, QRIS_BCA)
	Name             string          `gorm:"type:varchar(50);not null" json:"name"`                        // Nama metode bayar yang tampil di layar POS
	MdrPercentage    decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"mdr_percentage"`            // Potongan admin bank / Merchant Discount Rate
	DepositAccountID *uuid.UUID      `gorm:"type:char(36);index" json:"deposit_account_id,omitempty"`       // Akun penerima dana masuk (Bank BCA, Kas Laci) - nullable
	DepositAccount   *ChartOfAccount `gorm:"foreignKey:DepositAccountID" json:"deposit_account,omitempty"` // Relasi ke tabel ChartOfAccount
	ExpenseAccountID *uuid.UUID      `gorm:"type:char(36);index" json:"expense_account_id,omitempty"`                // Akun beban untuk pencatatan otomatis potongan MDR
	ExpenseAccount   *ChartOfAccount `gorm:"foreignKey:ExpenseAccountID" json:"expense_account,omitempty"` // Relasi ke tabel ChartOfAccount
	IsActive         bool            `gorm:"default:true" json:"is_active"`                                // Ketersediaan metode di terminal kasir
}