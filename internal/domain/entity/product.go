package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product merepresentasikan barang utama (Katalog Induk)
type Product struct {
	BaseModel
	SKU         string          `gorm:"type:varchar(50);not null" json:"sku"`            // Stock Keeping Unit
	Barcode     *string         `gorm:"type:varchar(50)" json:"barcode"`                 // Nomor Barcode
	Name        string          `gorm:"type:varchar(200);not null" json:"name"`          // Nama lengkap produk
	CategoryID  *uuid.UUID      `gorm:"type:char(36);index" json:"category_id"`          // Pemetaan ke kategori
	Category    ProductCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // Relasi ke ProductCategory
	BaseUOMID   uuid.UUID       `gorm:"type:char(36);not null;index" json:"base_uom_id"` // Satuan terkecil
	BaseUOM     UOM             `gorm:"foreignKey:BaseUOMID" json:"base_uom,omitempty"`  // Relasi ke UOM
	IsStockable bool            `gorm:"default:true" json:"is_stockable"`                // True jika barang fisik

	// Informasi Volumetrik & Planogram (Display Shelf Management)
	Length        decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"length"` // Panjang barang (cm)
	Width         decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"width"`  // Lebar barang (cm)
	Height        decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"height"` // Tinggi barang (cm)
	Weight        decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"weight"` // Berat barang (gram)
	IsStackable   bool            `gorm:"default:true" json:"is_stackable"`           // True = Barang bisa ditumpuk (karton/kotak). False = Rapuh/Bentuk tak beraturan.
	MaxStackLayer int             `gorm:"default:0" json:"max_stack_layer"`
}
