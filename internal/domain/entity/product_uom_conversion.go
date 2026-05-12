package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductUOMConversion menyimpan konversi multi-satuan untuk keperluan Grosir
type ProductUOMConversion struct {
	BaseModel
	ProductID      uuid.UUID       `gorm:"type:char(36);not null;index" json:"product_id"`     // Relasi ke produk spesifik
	Product        Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`      // Relasi ke tabel Product
	UOMID          uuid.UUID       `gorm:"type:char(36);not null;index" json:"uom_id"`         // Satuan grosir/alternatif (cth: DUS)
	UOM            UOM             `gorm:"foreignKey:UOMID" json:"uom,omitempty"`              // Relasi ke tabel UOM
	ConversionRate decimal.Decimal `gorm:"type:decimal(15,3);not null" json:"conversion_rate"` // Multiplier (contoh: 1 DUS = 24 PCS, isi 24.0000)
	Barcode        *string         `gorm:"type:varchar(50);uniqueIndex" json:"barcode"`        // Barcode spesifik untuk kemasan grosir ini

	// Informasi Volumetrik & Planogram Khusus Kemasan Grosir (Warehouse Management)
	Length        decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"length"` // Panjang kemasan grosir (cm)
	Width         decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"width"`  // Lebar kemasan grosir (cm)
	Height        decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"height"` // Tinggi kemasan grosir (cm)
	Weight        decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"weight"` // Berat kemasan grosir (gram/kg)
	IsStackable   bool            `gorm:"default:true" json:"is_stackable"`           // True = Kemasan grosir ini bisa ditumpuk di atas Pallet
	MaxStackLayer int             `gorm:"default:0" json:"max_stack_layer"`           // Batas maksimal tumpukan vertikal karton/krat agar tidak hancur
}
