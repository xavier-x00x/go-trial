package entity

import "time"

// PriceList mengakomodir harga berbeda (retail vs grosir) atau masa promo
// tabel header dari product_price
type PriceList struct {
	BaseModel
	Code         string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`  // Kode buku harga (cth: RETAIL-01)
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`             // Nama daftar harga
	CurrencyCode string     `gorm:"type:varchar(3);default:'IDR'" json:"currency_code"` // Kode mata uang (ISO 4217)
	StartDate    *time.Time `json:"start_date"`                                         // Tanggal mulai berlaku (untuk promo campaign)
	EndDate      *time.Time `json:"end_date"`                                           // Tanggal berakhir (opsional)
	IsActive     bool       `gorm:"default:true" json:"is_active"`                      // Status ketersediaan price list
}
