package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StoreProductAssortment mengatur status ketersediaan barang secara eksplisit di cabang tertentu
type StoreProductAssortment struct {
	BaseModel
	StoreID   uuid.UUID `gorm:"type:char(36);not null;index" json:"store_id"`                             // Relasi ke cabang/toko
	Store     Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`                                // Relasi ke tabel Store
	ProductID uuid.UUID `gorm:"type:char(36);not null;index" json:"product_id"`                           // Relasi ke produk
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`                            // Relasi ke tabel Product
	Status    string    `gorm:"type:enum('ACTIVE','DISCONTINUED','NOT_ASSORTED');not null" json:"status"` // Status katalog di cabang tsb

	DisplayFacingQty     int             `gorm:"default:1" json:"display_facing_qty"`                        // Berapa muka/baris barang ini dipajang di rak (Planogram per Cabang)
	DisplayShelfCapacity decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"display_shelf_capacity"` // Total kapasitas maksimum rak fisik khusus cabang ini (Facing x Kedalaman x Stack)

	// Min-Max Inventory Control (Spesifik per Cabang)
	VelocityClass        string          `gorm:"type:varchar(20);default:'MEDIUM'" json:"velocity_class"` // Kelas pergerakan cabang ini: FAST, MEDIUM, SLOW
	VelocityLookbackDays int             `gorm:"default:30" json:"velocity_lookback_days"`                // Jumlah hari histori ke belakang untuk hitung rata-rata penjualan cabang ini ( 7, 30, 90, 180 hari)
	SafetyStockQty       decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"safety_stock_qty"`    // Stok bantalan/penyangga cabang ini
	ReorderPointQty      decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"reorder_point_qty"`   // Titik trigger Pemesanan Kembali cabang ini
	MaxStockQty          decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"max_stock_qty"`       // Kapasitas maksimum finansial/fisik untuk cabang ini

	// Caching Metrik (Di-update secara berkala oleh Background Worker/Cron Job)
	AverageDailySales decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"average_daily_sales"` // Nilai cache rata-rata penjualan harian cabang ini
}

// PlanningData is joined data for planning calculation
type PlanningData struct {
	StoreID             uuid.UUID
	ProductID           uuid.UUID
	ProductSupplierID   uuid.UUID
	DefaultLeadTimeDays int
	AverageDailySales   decimal.Decimal
	SafetyStockQty      decimal.Decimal
	MaxStockQty        decimal.Decimal
	CurrentStock       decimal.Decimal // Total stock from all warehouses in this store
}

/*
Catatan:
- AverageDailySales = Total Penjualan Cabang ini dalam periode VelocityLookbackDays / VelocityLookbackDays
- ReorderPointQty = (AverageDailySales x LeadTimeDays) + SafetyStockQty
*/
