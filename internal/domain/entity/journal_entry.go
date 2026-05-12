package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta
// ──────────────────────────────────────────────────────────────────────────────

// SourceDocumentType mendefinisikan jenis dokumen sumber yang menghasilkan jurnal.
const (
	JournalSourcePurchaseInvoice = "PURCHASE_INVOICE"
	JournalSourcePurchasePayment = "PURCHASE_PAYMENT"
	JournalSourcePurchaseReturn  = "PURCHASE_RETURN"
	JournalSourceExpenseVoucher  = "EXPENSE_VOUCHER"
	JournalSourceSalesInvoice    = "SALES_INVOICE" // Untuk modul penjualan nanti
	JournalSourceSalesReturn     = "SALES_RETURN"  // Untuk modul penjualan nanti
	JournalSourceManualEntry     = "MANUAL_ENTRY"  // Jurnal Umum / Penyesuaian manual
)

// Status Jurnal
const (
	JournalStatusPosted   = "POSTED"   // Jurnal aktif dan berpengaruh ke buku besar
	JournalStatusReversed = "REVERSED" // Jurnal telah dibalik (void/cancel dari dokumen sumber)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: JournalEntry (Header)
// ──────────────────────────────────────────────────────────────────────────────

// JournalEntry adalah catatan jurnal akuntansi (General Ledger Entry).
// Setiap transaksi keuangan yang di-posting akan menghasilkan satu JournalEntry
// dengan dua atau lebih JournalEntryLine yang memenuhi prinsip double-entry:
//
//	Total Debit = Total Kredit
//
// Jurnal bisa dibuat otomatis oleh sistem (dari dokumen sumber) atau manual
// oleh akuntan (jurnal penyesuaian akhir bulan/tahun).
type JournalEntry struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	EntryNumber string `gorm:"type:varchar(50);uniqueIndex;not null" json:"entry_number"` // Nomor jurnal unik (cth: JE/2026/04/0001)

	// ── Dokumen Sumber ───────────────────────────────────────────────────
	SourceDocumentType string     `gorm:"type:varchar(30);not null;index" json:"source_document_type"` // Jenis dokumen sumber (PURCHASE_INVOICE, dll)
	SourceDocumentID   *uuid.UUID `gorm:"type:char(36);index" json:"source_document_id"`               // ID dokumen sumber. NULL jika MANUAL_ENTRY.
	SourceDocumentNo   *string    `gorm:"type:varchar(50)" json:"source_document_no"`                  // Nomor dokumen sumber untuk display (cth: PI/2026/04/0001)

	// ── Tanggal & Periode ────────────────────────────────────────────────
	EntryDate time.Time `gorm:"not null;index" json:"entry_date"`             // Tanggal jurnal (biasanya = tanggal dokumen sumber)
	Period    string    `gorm:"type:varchar(7);not null;index" json:"period"` // Periode akuntansi YYYY-MM (cth: 2026-04). Untuk filter laporan bulanan.

	// ── Nilai ────────────────────────────────────────────────────────────
	TotalDebit  decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"total_debit"`  // Total sisi Debit (harus = TotalCredit)
	TotalCredit decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"total_credit"` // Total sisi Kredit

	// ── Deskripsi ────────────────────────────────────────────────────────
	Description string `gorm:"type:varchar(200);not null" json:"description"` // Keterangan jurnal (cth: "Pembelian barang dari PT. Indofood")

	// ── Status & Tracking ────────────────────────────────────────────────
	Status       string        `gorm:"type:varchar(10);not null;default:'POSTED';index" json:"status"`
	ReversalOfID *uuid.UUID    `gorm:"type:char(36);index" json:"reversal_of_id"` // Jika ini jurnal balik, menunjuk ke jurnal asli yang dibalik
	ReversalOf   *JournalEntry `gorm:"foreignKey:ReversalOfID" json:"reversal_of,omitempty"`
	PostedByID   uuid.UUID     `gorm:"type:char(36);not null;index" json:"posted_by_id"`
	PostedBy     User          `gorm:"foreignKey:PostedByID" json:"posted_by,omitempty"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Lines []JournalEntryLine `gorm:"foreignKey:JournalEntryID" json:"lines,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: JournalEntryLine (Detail / Baris Jurnal)
// ──────────────────────────────────────────────────────────────────────────────

// JournalEntryLine adalah satu baris dalam jurnal akuntansi.
// Setiap baris mencatat pergerakan di satu akun: bisa Debit ATAU Kredit (tidak boleh keduanya).
type JournalEntryLine struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	JournalEntryID uuid.UUID    `gorm:"type:char(36);not null;index" json:"journal_entry_id"`
	JournalEntry   JournalEntry `gorm:"foreignKey:JournalEntryID" json:"journal_entry,omitempty"`
	SeqNo          int          `gorm:"not null" json:"seq_no"`

	// ── Akun ─────────────────────────────────────────────────────────────
	AccountID uuid.UUID      `gorm:"type:char(36);not null;index" json:"account_id"` // Akun COA yang terpengaruh
	Account   ChartOfAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`

	// ── Nilai ────────────────────────────────────────────────────────────
	DebitAmount  decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"debit_amount"`  // Nilai Debit (0 jika baris ini Kredit)
	CreditAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"credit_amount"` // Nilai Kredit (0 jika baris ini Debit)

	// ── Deskripsi ────────────────────────────────────────────────────────
	Description *string `gorm:"type:varchar(200)" json:"description"` // Keterangan per baris (opsional, override dari header)
}
