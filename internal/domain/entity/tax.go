package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Tax menyimpan konfigurasi pajak transaksional untuk faktur/struk
type Tax struct {
	BaseModel
	Name           string          `gorm:"type:varchar(50);not null" json:"name"`                // Label pajak (cth: PPN 11%, PB1)
	RatePercentage decimal.Decimal `gorm:"type:decimal(6,2);not null" json:"rate_percentage"`    // Besaran persentase tarif pajak (cth: 11.00)
	TaxAccountID   *uuid.UUID      `gorm:"type:char(36);index" json:"tax_account_id,omitempty"`  // Mapping ke Akun Liabilitas Pajak (Hutang Pajak) - nullable
	TaxAccount     *ChartOfAccount `gorm:"foreignKey:TaxAccountID" json:"tax_account,omitempty"` // Relasi ke tabel ChartOfAccount
}
