package entity

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel berisi kolom standar untuk semua tabel
type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"` // Identifier unik
	CreatedAt time.Time  `json:"created_at"`                         // Waktu pencatatan data
	UpdatedAt time.Time  `json:"updated_at"`                         // Waktu pembaruan data terakhir
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`  // Waktu penghapusan data (Soft Delete)
}

// GenerateID adalah fungsi murni (Pure Go) tanpa dependensi ke database/GORM.
// Fungsi ini dipanggil secara manual oleh lapisan Usecase sebelum menyimpan data.
func (b *BaseModel) GenerateID() error {
	if b.ID == uuid.Nil {
		u, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = u
	}
	return nil
}

// Meta Response
type Meta struct {
	Page          int `json:"page"`
	Limit         int `json:"limit"`
	Total         int `json:"total"`
	TotalFiltered int `json:"total_filtered"`
	LastPage      int `json:"last_page"`
	Draw          int `json:"draw"`
}

type QueryFilter struct {
	Page         int
	Limit        int
	Search       string
	OrderColumn  string
	OrderDir     string
	SearchColumn []string
	Conditions   map[string]interface{} // fleksibel
}
