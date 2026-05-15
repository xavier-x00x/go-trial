package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Status Purchase Payment
// ──────────────────────────────────────────────────────────────────────────────

const (
	PPStatusDraft  = "DRAFT"  // Pembayaran sedang disiapkan
	PPStatusPosted = "POSTED" // Dijurnal → Kas/Bank berkurang, Hutang berkurang
	PPStatusVoided = "VOIDED" // Dibatalkan setelah posting (jurnal balik)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchasePayment (Header)
// ──────────────────────────────────────────────────────────────────────────────

// PurchasePayment adalah dokumen pembayaran hutang ke supplier.
// Satu pembayaran bisa melunasi satu atau banyak faktur sekaligus (Batch Payment).
//
// Efek saat Status berubah ke POSTED:
//   - Jurnal Akuntansi: Debit Hutang Usaha (AP), Kredit Kas/Bank (PaymentAccountID)
//   - PurchaseInvoice.PaidAmount bertambah
//   - PurchaseInvoice.RemainingAmount berkurang
//   - PurchaseInvoice.Status → PARTIALLY_PAID atau PAID
//   - MonthlyAPBalance.TotalDebit bertambah
type PurchasePayment struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	PaymentNumber string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"payment_number"` // Nomor Bukti Bayar (cth: PAY/2026/04/0001)
	ReferenceNo   *string `gorm:"type:varchar(50)" json:"reference_no"`                        // Nomor referensi eksternal (cth: No. Transfer Bank, No. Giro)

	// ── Pihak Terkait ────────────────────────────────────────────────────
	SupplierID       uuid.UUID      `gorm:"type:char(36);not null;index" json:"supplier_id"`
	Supplier         Supplier       `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	PaymentAccountID uuid.UUID      `gorm:"type:char(36);not null;index" json:"payment_account_id"` // Akun Kas/Bank yang digunakan bayar (cth: Bank BCA, Kas Besar)
	PaymentAccount   ChartOfAccount `gorm:"foreignKey:PaymentAccountID" json:"payment_account,omitempty"`
	APAccountID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"ap_account_id"` // Akun Hutang yang dilunasi
	APAccount        ChartOfAccount `gorm:"foreignKey:APAccountID" json:"ap_account,omitempty"`

	// ── Tanggal & Metode ─────────────────────────────────────────────────
	PaymentDate time.Time  `gorm:"not null" json:"payment_date"`                  // Tanggal pembayaran dilakukan
	PaymentMode string     `gorm:"type:varchar(20);not null" json:"payment_mode"` // Metode: CASH, TRANSFER, GIRO
	GiroNumber  *string    `gorm:"type:varchar(50)" json:"giro_number"`           // Nomor Giro (wajib jika PaymentMode = GIRO)
	GiroDueDate *time.Time `json:"giro_due_date"`                                 // Tanggal jatuh tempo Giro

	// ── Nilai Transaksi ──────────────────────────────────────────────────
	TotalAmount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"total_amount"` // Total nilai pembayaran

	// ── Status & Tracking ────────────────────────────────────────────────
	Status      string     `gorm:"type:varchar(10);not null;default:'DRAFT';index" json:"status"`
	CreatedByID uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"`
	CreatedBy   User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	PostedByID  *uuid.UUID `gorm:"type:char(36);index" json:"posted_by_id"`
	PostedBy    *User      `gorm:"foreignKey:PostedByID" json:"posted_by,omitempty"`
	PostedAt    *time.Time `json:"posted_at"`

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes *string `gorm:"type:text" json:"notes"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []PurchasePaymentItem `gorm:"foreignKey:PurchasePaymentID" json:"items,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchasePaymentItem (Detail / Alokasi per Faktur)
// ──────────────────────────────────────────────────────────────────────────────

// PurchasePaymentItem merepresentasikan alokasi pembayaran ke masing-masing faktur.
// Satu pembayaran bisa dialokasikan ke banyak faktur (Batch Payment).
//
// Contoh: Bayar Rp 10.000.000 ke PT. Indofood, dialokasikan:
//   - Faktur PI/2026/03/0012: Rp 6.000.000 (lunas)
//   - Faktur PI/2026/03/0018: Rp 4.000.000 (cicilan, sisa Rp 2.000.000)
type PurchasePaymentItem struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	PurchasePaymentID uuid.UUID       `gorm:"type:char(36);not null;index" json:"purchase_payment_id"`
	PurchasePayment   PurchasePayment `gorm:"foreignKey:PurchasePaymentID" json:"purchase_payment,omitempty"`
	SeqNo             int             `gorm:"not null" json:"seq_no"`

	// ── Dokumen yang Dibayar / Di-offset ────────────────────────────────
	PurchaseInvoiceID *uuid.UUID       `gorm:"type:char(36);index" json:"purchase_invoice_id"`
	PurchaseInvoice   *PurchaseInvoice `gorm:"foreignKey:PurchaseInvoiceID" json:"purchase_invoice,omitempty"`
	PurchaseReturnID  *uuid.UUID       `gorm:"type:char(36);index" json:"purchase_return_id"`
	PurchaseReturn    *PurchaseReturn  `gorm:"foreignKey:PurchaseReturnID" json:"purchase_return,omitempty"`

	// ── Alokasi Nilai ────────────────────────────────────────────────────
	DocumentAmount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"document_amount"` // Snapshot total hutang/retur (untuk referensi)
	PaidAmount     decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"paid_amount"`     // Nilai yang dialokasikan (positif untuk Invoice, negatif untuk Return)
}
