package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta
// ──────────────────────────────────────────────────────────────────────────────

// Status ExpenseVoucher
const (
	EVStatusDraft    = "DRAFT"    // Voucher baru dibuat
	EVStatusApproved = "APPROVED" // Disetujui oleh Manajer
	EVStatusPosted   = "POSTED"   // Dijurnal → Beban tercatat, Kas/Bank berkurang
	EVStatusVoided   = "VOIDED"   // Dibatalkan setelah posting (jurnal balik)
)

// Tipe Pembayaran
const (
	EVPaymentTypeCash   = "CASH"   // Langsung dibayar tunai/transfer saat itu juga
	EVPaymentTypeCredit = "CREDIT" // Hutang ke vendor, dibayar nanti (tempo)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: ExpenseVoucher (Header)
// ──────────────────────────────────────────────────────────────────────────────

// ExpenseVoucher adalah dokumen Bukti Pengeluaran untuk pembelian non-dagangan.
// Contoh: ATK, jasa service AC, biaya listrik, biaya kebersihan, dll.
//
// Perbedaan dengan PurchaseInvoice:
//   - Tidak terkait dengan Product Master (barang tidak masuk stok)
//   - Tidak memerlukan PO dan Goods Receipt
//   - Jurnal langsung ke akun Beban (Expense), bukan Persediaan (Inventory)
//
// Efek saat Status berubah ke POSTED:
//   - Jika PaymentType = CASH:
//     Jurnal: Debit Beban (ExpenseAccountID), Kredit Kas/Bank (PaymentAccountID)
//   - Jika PaymentType = CREDIT:
//     Jurnal: Debit Beban (ExpenseAccountID), Kredit Hutang Lain (PayableAccountID)
type ExpenseVoucher struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	VoucherNumber string `gorm:"type:varchar(50);uniqueIndex;not null" json:"voucher_number"` // Nomor Bukti Pengeluaran (cth: EXP/2026/04/0001)

	// ── Pihak Terkait ────────────────────────────────────────────────────
	StoreID    uuid.UUID  `gorm:"type:char(36);not null;index" json:"store_id"` // Cabang yang mengeluarkan biaya
	Store      Store      `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	SupplierID *uuid.UUID `gorm:"type:char(36);index" json:"supplier_id"` // Vendor (opsional, bisa tidak ada di master supplier)
	Supplier   *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	VendorName *string    `gorm:"type:varchar(150)" json:"vendor_name"` // Nama vendor ad-hoc (jika tidak ada di master supplier)

	// ── Tanggal ──────────────────────────────────────────────────────────
	VoucherDate time.Time `gorm:"not null" json:"voucher_date"` // Tanggal transaksi pengeluaran

	// ── Pembayaran ───────────────────────────────────────────────────────
	PaymentType      string          `gorm:"type:varchar(10);not null" json:"payment_type"` // CASH atau CREDIT
	PaymentAccountID *uuid.UUID      `gorm:"type:char(36);index" json:"payment_account_id"` // Akun Kas/Bank (wajib jika CASH)
	PaymentAccount   *ChartOfAccount `gorm:"foreignKey:PaymentAccountID" json:"payment_account,omitempty"`
	PayableAccountID *uuid.UUID      `gorm:"type:char(36);index" json:"payable_account_id"` // Akun Hutang Lain-lain (wajib jika CREDIT)
	PayableAccount   *ChartOfAccount `gorm:"foreignKey:PayableAccountID" json:"payable_account,omitempty"`

	// ── Nilai Transaksi ──────────────────────────────────────────────────
	Subtotal   decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"subtotal"`
	TaxAmount  decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"tax_amount"`
	GrandTotal decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"grand_total"`

	// ── Status & Tracking ────────────────────────────────────────────────
	Status       string     `gorm:"type:varchar(10);not null;default:'DRAFT';index" json:"status"`
	CreatedByID  uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"`
	CreatedBy    User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	ApprovedByID *uuid.UUID `gorm:"type:char(36);index" json:"approved_by_id"`
	ApprovedBy   *User      `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at"`
	PostedByID   *uuid.UUID `gorm:"type:char(36);index" json:"posted_by_id"`
	PostedBy     *User      `gorm:"foreignKey:PostedByID" json:"posted_by,omitempty"`
	PostedAt     *time.Time `json:"posted_at"`

	// ── Catatan ──────────────────────────────────────────────────────────
	Description string  `gorm:"type:varchar(200);not null" json:"description"` // Deskripsi singkat keperluan (cth: "Pembelian ATK Bulan April")
	Notes       *string `gorm:"type:text" json:"notes"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []ExpenseVoucherItem `gorm:"foreignKey:ExpenseVoucherID" json:"items,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: ExpenseVoucherItem (Detail / Baris Pengeluaran)
// ──────────────────────────────────────────────────────────────────────────────

// ExpenseVoucherItem adalah baris detail pengeluaran.
// Setiap baris bisa menunjuk ke akun beban yang berbeda.
//
// Contoh 1 voucher dengan 3 baris:
//   - Baris 1: Kertas HVS 5 rim × Rp 50.000   → Akun: Beban ATK (6110)
//   - Baris 2: Tinta Printer 2 pcs × Rp 150.000 → Akun: Beban ATK (6110)
//   - Baris 3: Jasa Service AC 1 unit × Rp 500.000 → Akun: Beban Pemeliharaan (6210)
type ExpenseVoucherItem struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	ExpenseVoucherID uuid.UUID      `gorm:"type:char(36);not null;index" json:"expense_voucher_id"`
	ExpenseVoucher   ExpenseVoucher `gorm:"foreignKey:ExpenseVoucherID" json:"expense_voucher,omitempty"`
	SeqNo            int            `gorm:"not null" json:"seq_no"`

	// ── Akun Beban ───────────────────────────────────────────────────────
	ExpenseAccountID uuid.UUID      `gorm:"type:char(36);not null;index" json:"expense_account_id"` // Akun beban tujuan (cth: Beban ATK, Beban Listrik)
	ExpenseAccount   ChartOfAccount `gorm:"foreignKey:ExpenseAccountID" json:"expense_account,omitempty"`

	// ── Detail Item ──────────────────────────────────────────────────────
	Description string          `gorm:"type:varchar(200);not null" json:"description"` // Nama barang/jasa (free text, tidak dari master product)
	Qty         decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"qty"`
	UnitPrice   decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"unit_price"`
	TaxPct      decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"tax_pct"`
	TaxAmount   decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"tax_amount"`
	Subtotal    decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"subtotal"` // Qty × UnitPrice

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes *string `gorm:"type:text" json:"notes"`
}
