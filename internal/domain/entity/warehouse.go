package entity

import "github.com/google/uuid"

// Warehouse merepresentasikan Gudang atau Lokasi Penyimpanan
type Warehouse struct {
	BaseModel
	StoreID  uuid.UUID `gorm:"type:char(36);not null;index" json:"store_id"`      // Relasi ke toko kepemilikan
	Store    Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`         // Relasi ke tabel Store
	Code     string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"` // Kode referensi gudang (cth: WH-DEPAN)
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`            // Nama deskriptif gudang penyimpanan
	IsActive bool      `gorm:"default:true" json:"is_active"`                     // Status gudang aktif menerima barang
}

func (w *Warehouse) TableName() string {
	return "warehouse"
}
