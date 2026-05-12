package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductPrice memisahkan harga jual dari tabel produk untuk fleksibilitas (Dynamic Pricing)
type ProductPrice struct {
	BaseModel
	PriceListID uuid.UUID       `gorm:"type:char(36);not null;index" json:"price_list_id"`  // Relasi ke daftar harga
	PriceList   PriceList       `gorm:"foreignKey:PriceListID" json:"price_list,omitempty"` // Relasi ke tabel PriceList
	ProductID   uuid.UUID       `gorm:"type:char(36);not null;index" json:"product_id"`     // Relasi ke entitas produk
	Product     Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`      // Relasi ke tabel Product
	UOMID       uuid.UUID       `gorm:"type:char(36);not null;index" json:"uom_id"`         // Satuan spesifik harga (1 PCS beda margin dgn 1 DUS)
	UOM         UOM             `gorm:"foreignKey:UOMID" json:"uom,omitempty"`              // Relasi ke tabel UOM
	MarkupPct   decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"markup_pct"`      // Persentase Markup. Jika diisi > 0, SellPrice otomatis menjadi: HPP + (HPP * Markup%).
	SellPrice   decimal.Decimal `gorm:"type:decimal(19,2);not null" json:"sell_price"`      // Harga Jual Final ke konsumen
	DiscountPct decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_pct"`    // Persentase diskon langsung (Markdown)
}
