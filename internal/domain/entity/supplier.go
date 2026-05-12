package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Supplier merepresentasikan vendor atau pemasok barang masuk
type Supplier struct {
	BaseModel
	Code                        string            `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`                     // Kode register vendor ERP
	Name                        string            `gorm:"type:varchar(150);not null" json:"name"`                                // Nama bisnis distributor/pabrik
	ContactPerson               *string           `gorm:"type:varchar(100)" json:"contact_person"`                               // PIC dari supplier (Nama Sales)
	ContactPhone                *string           `gorm:"type:varchar(20)" json:"contact_phone"`                                 // No HP/WhatsApp spesifik milik PIC
	PhoneNumber                 *string           `gorm:"type:varchar(20)" json:"phone_number"`                                  // Kontak seluler/kantor supplier
	Email                       *string           `gorm:"type:varchar(100)" json:"email"`                                        // Email untuk pengiriman PO
	PreferredNotificationMethod string            `gorm:"type:varchar(20);default:'EMAIL'" json:"preferred_notification_method"` // WHATSAPP, EMAIL, NONE
	Address                     *string           `gorm:"type:text" json:"address"`                                              // Alamat fisik perusahaan/gudang
	TaxRegNumber                *string           `gorm:"type:varchar(50)" json:"tax_reg_number"`                                // NPWP Supplier untuk faktur
	SupplierCategoryID          *uuid.UUID        `gorm:"type:char(36);index" json:"supplier_category_id"`                       // Relasi ke kategori komoditas supplier
	SupplierCategory            *SupplierCategory `gorm:"foreignKey:SupplierCategoryID" json:"supplier_category,omitempty"`      // Data relasi Kategori
	IsPKP                       bool              `gorm:"default:false" json:"is_pkp"`                                           // Status Pengusaha Kena Pajak (Menerbitkan Faktur Pajak PPN)
	PaymentTermDays             int               `gorm:"default:0" json:"payment_term_days"`                                    // Termin hutang (0 = Cash/COD, 30 = Tempo 30 Hari)
	PaymentMode                 string            `gorm:"type:varchar(20);default:'TRANSFER'" json:"payment_mode"`               // Cara pembayaran (CASH, TRANSFER, GIRO)
	MinOrderAmount              decimal.Decimal   `gorm:"type:decimal(19,2);default:0" json:"min_order_amount"`                  // Minimum total nilai belanja (Rupiah) agar diproses/gratis ongkir
	BankName                    *string           `gorm:"type:varchar(50)" json:"bank_name"`                                     // Nama Bank Supplier (cth: BCA)
	BankAccount                 *string           `gorm:"type:varchar(50)" json:"bank_account"`                                  // Nomor Rekening Bank Supplier
	BankAccountName             *string           `gorm:"type:varchar(150)" json:"bank_account_name"`                            // Nama Pemilik Rekening (Atas Nama)
	IsActive                    bool              `gorm:"default:true" json:"is_active"`                                         // Status kerjasama aktif/diblokir
	APAccountID                 *uuid.UUID        `gorm:"type:char(36);index" json:"ap_account_id"`                              // Akun Utang Usaha (Account Payable) spesifik
	APAccount                   *ChartOfAccount   `gorm:"foreignKey:APAccountID" json:"ap_account,omitempty"`                    // Relasi ke tabel ChartOfAccount
}
