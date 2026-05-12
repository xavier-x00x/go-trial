package entity

// SupplierCategory untuk mengelompokkan vendor berdasarkan jenis komoditas yang mereka suplai
type SupplierCategory struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"` // Nama Kategori (cth: Sayuran Segar, FMCG, Kemasan)
	Description string `gorm:"type:text" json:"description"`                       // Keterangan opsional mengenai kategori ini
}
