package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Status
// ──────────────────────────────────────────────────────────────────────────────

const (
	PlanningStatusPending  = "PENDING"  // Menunggu keputusan staf Purchasing
	PlanningStatusApproved = "APPROVED" // Disetujui → akan dibuatkan PO
	PlanningStatusIgnored  = "IGNORED"  // Diabaikan (stok masih cukup / keputusan manual)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseOrderPlanning adalah rekomendasi pemesanan yang di-generate secara otomatis
// oleh Background Worker/Cron Job berdasarkan algoritma Reorder Point & Dynamic Safety Stock.
//
// Alur Kerja:
//  1. Cron Job berjalan setiap malam (atau interval tertentu).
//  2. Sistem menghitung stok vs ROP per produk per cabang.
//  3. Jika stok <= ROP → sistem membuat record Planning dengan Status = PENDING.
//  4. Staf Purchasing meninjau daftar rekomendasi.
//  5. Jika Approve → staf memilih item, lalu sistem membuatkan PurchaseOrder.
//  6. Jika Ignore → rekomendasi ditandai IGNORED (tidak dibuatkan PO).
type PurchaseOrderPlanning struct {
	BaseModel

	// ── Identifikasi Target ──────────────────────────────────────────────
	StoreID           uuid.UUID       `gorm:"type:char(36);not null;index" json:"store_id"`                   // Cabang yang membutuhkan barang
	Store             Store           `gorm:"foreignKey:StoreID" json:"store,omitempty"`                      // Relasi ke Store
	ProductID         uuid.UUID       `gorm:"type:char(36);not null;index" json:"product_id"`                 // Produk yang perlu dipesan
	Product           Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`                  // Relasi ke Product
	ProductSupplierID uuid.UUID       `gorm:"type:char(36);not null;index" json:"product_supplier_id"`        // Supplier yang direkomendasikan (Primary Supplier)
	ProductSupplier   ProductSupplier `gorm:"foreignKey:ProductSupplierID" json:"product_supplier,omitempty"` // Relasi ke ProductSupplier

	// ── Snapshot Metrik Saat Perhitungan ──────────────────────────────────
	// Nilai-nilai ini di-snapshot pada saat CalculatedDate agar bisa di-audit.
	// Jangan dibandingkan dengan data real-time karena stok bisa berubah setiap detik.
	CurrentStock       decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"current_stock"`        // Total stok GABUNGAN semua gudang dalam 1 toko (SUM InventoryStock WHERE store's warehouses)
	SafetyStock        decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"safety_stock"`         // Safety Stock statis (dari StoreProductAssortment)
	DynamicSafetyStock decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"dynamic_safety_stock"` // Safety Stock dinamis (hasil kalkulasi algoritma)
	MaxStockQty        decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"max_stock_qty"`        // Kapasitas maksimum stok cabang
	ReorderPoint       decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"reorder_point"`        // Titik pemesanan kembali (ROP)
	AverageDailySales  decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"average_daily_sales"`  // Rata-rata penjualan harian
	LeadTimeDays       int             `gorm:"not null" json:"lead_time_days"`                          // Estimasi waktu pengiriman (hari)
	LeadTimeDemand     decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"lead_time_demand"`     // Permintaan selama lead time (ADS × LeadTimeDays)

	RecommendedOrderQty decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"recommended_order_qty"`        // Qty pesanan yang disarankan (MaxStock - CurrentStock)
	OrderQty            decimal.Decimal `gorm:"type:decimal(15,4);not null;default:0" json:"order_qty"`          // Qty pesanan real yang bisa diedit user
	IsManualSupplier    bool            `gorm:"default:false" json:"is_manual_supplier"`                         // Flag penanda supplier direvisi manual
	IsSelected          bool            `gorm:"default:false" json:"is_selected"`                                // Flag penanda baris dipilih oleh user
	CalculatedDate      time.Time       `gorm:"not null" json:"calculated_date"`                                 // Waktu kalkulasi dilakukan oleh Cron Job
	Status              string          `gorm:"type:varchar(10);not null;default:'PENDING';index" json:"status"` // PENDING, APPROVED, IGNORED
	ProcessedByID       *uuid.UUID      `gorm:"type:char(36);index" json:"processed_by_id"`                      // User yang mengambil keputusan
	ProcessedBy         *User           `gorm:"foreignKey:ProcessedByID" json:"processed_by,omitempty"`          // Relasi ke User yang memproses keputusan
	ProcessedDate       *time.Time      `json:"processed_date"`                                                  // Waktu keputusan diambil
	PurchaseOrderID     *uuid.UUID      `gorm:"type:char(36);index" json:"purchase_order_id"`                    // FK ke PO yang dihasilkan (NULL jika IGNORED)
}
