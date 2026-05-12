package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// InventoryStock melacak saldo fisik per lokasi diukur dalam satuan dasar (base_uom)
type InventoryStock struct {
	BaseModel
	WarehouseID     uuid.UUID       `gorm:"type:char(36);not null;index" json:"warehouse_id"`      // Relasi ke gudang penyimpan
	Warehouse       Warehouse       `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`     // Relasi ke tabel Warehouse
	ProductID       uuid.UUID       `gorm:"type:char(36);not null;index" json:"product_id"`        // Relasi ke produk
	Product         Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`         // Relasi ke tabel Product
	Quantity        decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"quantity"`          // Saldo fisik ketersediaan stok saat ini
	ReservedQty     decimal.Decimal `gorm:"type:decimal(15,3);default:0" json:"reserved_qty"`      // Stok di-booking pesanan (belum diambil/dikirim)
	AverageBuyPrice decimal.Decimal `gorm:"type:decimal(19,2);default:0" json:"average_buy_price"` // HPP / Harga Beli Rata-rata Bergerak per Gudang
	LastAuditAt     *time.Time      `json:"last_audit_at"`                                         // Waktu terakhir dilakukan Stock Opname
}

/*
	AverageBuyPrice :
	Bulan Lalu: Anda beli 10 sak dengan harga Rp 50.000/sak. Modal awal di sistem tercatat BuyPrice = 50.000.
	Hari Ini (Transaksi): Harga sembako naik. Anda beli lagi 10 sak dari supplier dengan harga Rp 60.000/sak.
		Saat Manajer menekan tombol [ACC Faktur Pembelian], Backend Golang Anda akan langsung bekerja.
	Kalkulasi Rata-rata: ((10 sak lama * 50.000) + (10 sak baru * 60.000)) / Total 20 sak = Rp 55.000
*/
