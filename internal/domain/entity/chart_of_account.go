package entity

// ChartOfAccount (Bagan Akun) berfungsi sebagai motor penjurnalan ganda akuntansi (Double-Entry)
type ChartOfAccount struct {
	BaseModel
	AccountCode   string `gorm:"type:varchar(20);uniqueIndex;not null" json:"account_code"`                                // Nomor akun konvensi (cth: 1110.01)
	Name          string `gorm:"type:varchar(100);not null" json:"name"`                                                   // Nama entitas akun (cth: Kas di Tangan, HPP Grosir)
	AccountType   string `gorm:"type:enum('ASSET','LIABILITY','EQUITY','REVENUE','EXPENSE');not null" json:"account_type"` // Tipe Akun (Harta, Hutang, Modal, Pendapatan, Beban)
	NormalBalance string `gorm:"type:enum('DEBIT','CREDIT');not null" json:"normal_balance"`                               // Sifat penambahan saldo normal (Debit/Credit)
	IsActive      bool   `gorm:"default:true" json:"is_active"`                                                            // Status penggunaan akun untuk penjurnalan
}
