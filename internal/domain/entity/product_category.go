package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductCategory menggunakan pola Adjacency List untuk hierarki kategori tak terbatas
type ProductCategory struct {
	BaseModel
	ParentID         *uuid.UUID       `gorm:"type:char(36);index" json:"parent_id"`                  // Kategori induk. Null jika Kategori Utama (Root)
	Parent           *ProductCategory `gorm:"foreignKey:ParentID" json:"parent,omitempty"`           // Relasi ke kategori induk
	Name             string           `gorm:"type:varchar(100);not null" json:"name"`                // Nama kategori (cth: Sembako)
	Slug             string           `gorm:"type:varchar(120);uniqueIndex;not null" json:"slug"`    // SEO/URL Friendly identifier
	Code             string           `gorm:"type:char(3);uniqueIndex;not null" json:"code"`         // Kode kategori (3 huruf besar)
	Sequence         int              `gorm:"type:int;default:0" json:"sequence"`                    // Sequence number untuk SKU generator
	DefaultMarkupPct decimal.Decimal  `gorm:"type:decimal(7,2);default:0" json:"default_markup_pct"` // Markup default untuk auto-pricing (Fallback System)
}
