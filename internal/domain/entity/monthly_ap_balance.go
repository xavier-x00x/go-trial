package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MonthlyAPBalance (Mutasi Saldo Hutang Bulanan) menggunakan pola Roll-Forward
// Menyimpan riwayat pergerakan hutang per supplier untuk setiap periode bulan
type MonthlyAPBalance struct {
	BaseModel
	PeriodMonth      string          `gorm:"type:varchar(7);not null;index" json:"period_month"` // Format: "YYYY-MM" (contoh: "2026-03")
	SupplierID       uuid.UUID       `gorm:"type:char(36);not null;index" json:"supplier_id"`
	Supplier         Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	BeginningBalance decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"beginning_balance"` // Saldo awal (Bawaan dari EndingBalance bulan sebelumnya)
	TotalDebit       decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"total_debit"`       // Total pembayaran cicilan/lunas di bulan ini (-)
	TotalCredit      decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"total_credit"`      // Total faktur barang masuk di bulan ini (+)
	EndingBalance    decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"ending_balance"`    // Beginning + Credit - Debit (Saldo riil saat ini)
}
