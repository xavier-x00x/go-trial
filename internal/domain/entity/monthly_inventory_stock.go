package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// MonthlyInventoryStock (Kartu Stok / Mutasi Persediaan Bulanan)
// Menyimpan riwayat pergerakan fisik barang per gudang untuk setiap periode bulan.
type MonthlyInventoryStock struct {
	BaseModel
	PeriodMonth      string          `gorm:"type:varchar(7);not null;index" json:"period_month"` // Format: "YYYY-MM" (contoh: "2026-03")
	WarehouseID      uuid.UUID       `gorm:"type:char(36);not null;index" json:"warehouse_id"`   // Relasi ke Gudang
	Warehouse        Warehouse       `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	ProductID        uuid.UUID       `gorm:"type:char(36);not null;index" json:"product_id"` // Relasi ke Produk
	Product          Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	BeginningBalance decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"beginning_balance"` // Saldo awal (Bawaan dari EndingBalance bulan sebelumnya)
	TotalIn          decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"total_in"`          // Total barang masuk (Pembelian, Retur Pelanggan, Mutasi Masuk)
	TotalOut         decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"total_out"`         // Total barang keluar (Penjualan Kasir, Retur Supplier, Mutasi Keluar)
	EndingBalance    decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"ending_balance"`    // Beginning + In - Out (Sama persis dengan saldo di tabel InventoryStock)
	EndingValue      decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"ending_value"`      // Total Nilai Aset Uang (EndingBalance * AverageBuyPrice di akhir bulan)
}
