package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Status Purchase Return
// ──────────────────────────────────────────────────────────────────────────────

const (
	PRStatusDraft     = "DRAFT"     // Dokumen retur sedang disiapkan
	PRStatusConfirmed = "CONFIRMED" // Dikonfirmasi → stok berkurang, debit note ke supplier
	PRStatusClosed    = "CLOSED"    // Supplier sudah mengakui (nilai hutang sudah dipotong)
	PRStatusCancelled = "CANCELLED" // Dibatalkan (hanya dari DRAFT)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchaseReturn (Header)
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseReturn adalah dokumen retur barang ke supplier.
// Retur bisa terjadi karena barang rusak, expired, salah kirim, atau kelebihan qty.
//
// Efek saat Status berubah ke CONFIRMED:
//   - InventoryStock.Quantity berkurang sebesar QtyReturned (barang keluar dari gudang)
//   - InventoryStock.AverageBuyPrice dihitung ulang
//   - Debit Note / Nota Retur diterbitkan ke supplier
//   - PurchaseInvoice.RemainingAmount berkurang (hutang dipotong senilai retur)
//   - MonthlyAPBalance.TotalDebit bertambah
type PurchaseReturn struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	ReturnNumber      string           `gorm:"type:varchar(50);uniqueIndex;not null" json:"return_number"` // Nomor Retur unik (cth: PR/2026/04/0001)
	PurchaseInvoiceID *uuid.UUID       `gorm:"type:char(36);index" json:"purchase_invoice_id"`             // Referensi ke Faktur (opsional, retur bisa sebelum invoice)
	PurchaseInvoice   *PurchaseInvoice `gorm:"foreignKey:PurchaseInvoiceID" json:"purchase_invoice,omitempty"`
	PurchaseOrderID   uuid.UUID        `gorm:"type:char(36);not null;index" json:"purchase_order_id"` // Referensi ke PO asal barang
	PurchaseOrder     PurchaseOrder    `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`

	// ── Pihak Terkait ────────────────────────────────────────────────────
	SupplierID  uuid.UUID `gorm:"type:char(36);not null;index" json:"supplier_id"`
	Supplier    Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	WarehouseID uuid.UUID `gorm:"type:char(36);not null;index" json:"warehouse_id"` // Gudang asal barang yang diretur
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`

	// ── Tanggal ──────────────────────────────────────────────────────────
	ReturnDate time.Time `gorm:"not null" json:"return_date"` // Tanggal retur dilakukan

	// ── Nilai Transaksi ──────────────────────────────────────────────────
	TotalAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"total_amount"` // Total nilai retur (pengurang hutang)

	// ── Status & Tracking ────────────────────────────────────────────────
	Status        string     `gorm:"type:varchar(10);not null;default:'DRAFT';index" json:"status"`
	CreatedByID   uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"`
	CreatedBy     User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	ConfirmedByID *uuid.UUID `gorm:"type:char(36);index" json:"confirmed_by_id"`
	ConfirmedBy   *User      `gorm:"foreignKey:ConfirmedByID" json:"confirmed_by,omitempty"`
	ConfirmedAt   *time.Time `json:"confirmed_at"`

	// ── Catatan ──────────────────────────────────────────────────────────
	Reason string  `gorm:"type:text;not null" json:"reason"` // Alasan retur (wajib diisi)
	Notes  *string `gorm:"type:text" json:"notes"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []PurchaseReturnItem `gorm:"foreignKey:PurchaseReturnID" json:"items,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchaseReturnItem (Detail)
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseReturnItem adalah baris detail barang yang diretur ke supplier.
type PurchaseReturnItem struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	PurchaseReturnID uuid.UUID      `gorm:"type:char(36);not null;index" json:"purchase_return_id"`
	PurchaseReturn   PurchaseReturn `gorm:"foreignKey:PurchaseReturnID" json:"purchase_return,omitempty"`
	SeqNo            int            `gorm:"not null" json:"seq_no"`

	// ── Referensi ────────────────────────────────────────────────────────
	GoodsReceiptItemID *uuid.UUID        `gorm:"type:char(36);index" json:"goods_receipt_item_id"` // Baris GR mana yang diretur
	GoodsReceiptItem   *GoodsReceiptItem `gorm:"foreignKey:GoodsReceiptItemID" json:"goods_receipt_item,omitempty"`

	// ── Produk & Satuan ──────────────────────────────────────────────────
	ProductID uuid.UUID `gorm:"type:char(36);not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UOMID     uuid.UUID `gorm:"type:char(36);not null" json:"uom_id"`
	UOM       UOM       `gorm:"foreignKey:UOMID" json:"uom,omitempty"`

	// ── Kuantitas & Harga ────────────────────────────────────────────────
	QtyReturned decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"qty_returned"` // Jumlah yang dikembalikan ke supplier
	UnitPrice   decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"unit_price"`   // Harga beli per satuan (dari GR/Invoice)
	Subtotal    decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"subtotal"`     // QtyReturned × UnitPrice

	// ── Catatan ──────────────────────────────────────────────────────────
	ReturnReason string  `gorm:"type:text;not null" json:"return_reason"` // Alasan retur per item (Rusak, Expired, Salah Kirim, Kelebihan)
	Notes        *string `gorm:"type:text" json:"notes"`
}
