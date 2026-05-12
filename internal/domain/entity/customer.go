package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Customer merepresentasikan pelanggan, member loyalitas, atau agen
type Customer struct {
	BaseModel
	Code         string          `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`  // Kode referensi pelanggan/member
	Name         string          `gorm:"type:varchar(150);not null" json:"name"`             // Nama individu/institusi
	PhoneNumber  *string         `gorm:"type:varchar(20);uniqueIndex" json:"phone_number"`   // No WhatsApp untuk integrasi loyalitas
	Email        *string         `gorm:"type:varchar(100)" json:"email"`                     // Email pelanggan untuk e-receipt
	Address      *string         `gorm:"type:text" json:"address"`                           // Alamat domisili/pengiriman
	IsActive     bool            `gorm:"default:true" json:"is_active"`                      // Status member aktif/diblokir
	PointBalance decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"point_balance"`  // Saldo poin akumulatif reward
	CreditLimit  decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"credit_limit"`   // Plafon utang (kasbon khusus pelanggan B2B/Grosir)
	ARAccountID  *uuid.UUID      `gorm:"type:char(36);index" json:"ar_account_id"`           // Akun Piutang Usaha (Account Receivable) spesifik
	ARAccount    *ChartOfAccount `gorm:"foreignKey:ARAccountID" json:"ar_account,omitempty"` // Relasi ke tabel ChartOfAccount
}
