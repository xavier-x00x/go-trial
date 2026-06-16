package entity

import "github.com/google/uuid"

// PriceList mengakomodir harga berbeda (retail vs grosir) atau masa promo
// tabel header dari product_price
type PriceList struct {
	BaseModel
	Code         string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`  // Kode buku harga (cth: RETAIL-01)
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`             // Nama daftar harga
	CurrencyCode string     `gorm:"type:varchar(3);default:'IDR'" json:"currency_code"` // Kode mata uang (ISO 4217)
	StoreID      *uuid.UUID `gorm:"type:char(36);index" json:"store_id,omitempty"`
	Store        *Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"` // Status ketersediaan price list
}
