package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Status PO
// ──────────────────────────────────────────────────────────────────────────────

const (
	POStatusDraft             = "DRAFT"              // PO baru dibuat, belum dikirim ke supplier
	POStatusSubmitted         = "SUBMITTED"          // PO telah dikirim/dikomunikasikan ke supplier
	POStatusApproved          = "APPROVED"           // Disetujui oleh Manajer (otorisasi internal)
	POStatusPartiallyReceived = "PARTIALLY_RECEIVED" // Sebagian barang sudah diterima di gudang
	POStatusReceived          = "RECEIVED"           // Seluruh barang sudah diterima
	POStatusClosed            = "CLOSED"             // PO selesai (sudah di-invoice & matched)
	POStatusCancelled         = "CANCELLED"          // PO dibatalkan
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Notification Status
// ──────────────────────────────────────────────────────────────────────────────

const (
	NotificationStatusNone    = "NONE"    // Tidak perlu kirim notifikasi
	NotificationStatusPending = "PENDING" // Menunggu proses pengiriman
	NotificationStatusSent    = "SENT"    // Berhasil terkirim
	NotificationStatusFailed  = "FAILED"  // Gagal kirim
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchaseOrder (Header)
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseOrder adalah dokumen resmi pemesanan barang ke Supplier.
//
// Siklus Hidup:
//
//	DRAFT → SUBMITTED → APPROVED → PARTIALLY_RECEIVED → RECEIVED → CLOSED
//	                                    (atau langsung RECEIVED jika 1x kirim)
//	Bisa CANCELLED kapan saja sebelum RECEIVED.
type PurchaseOrder struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	PONumber    string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"po_number"` // Nomor PO unik (cth: PO/2026/04/0001)
	ReferenceNo *string `gorm:"type:varchar(50)" json:"reference_no"`                   // Nomor referensi eksternal (cth: No. Quotation dari supplier)

	// ── Pihak Terkait ────────────────────────────────────────────────────
	SupplierID  uuid.UUID `gorm:"type:char(36);not null;index" json:"supplier_id"`
	Supplier    Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	StoreID     uuid.UUID `gorm:"type:char(36);not null;index" json:"store_id"` // Cabang yang memesan
	Store       Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	WarehouseID uuid.UUID `gorm:"type:char(36);not null;index" json:"warehouse_id"` // Gudang tujuan penerimaan
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`

	// ── Tanggal & Jadwal ─────────────────────────────────────────────────
	OrderDate        time.Time  `gorm:"not null" json:"order_date"` // Tanggal PO dibuat/dikirim
	ExpectedDelivery *time.Time `json:"expected_delivery"`          // Estimasi tanggal barang tiba di gudang

	// ── Syarat Pembayaran (Override dari Supplier Master) ────────────────
	PaymentTermDays int    `gorm:"default:0" json:"payment_term_days"`                      // Termin hutang (bisa beda dari default supplier)
	PaymentMode     string `gorm:"type:varchar(20);default:'TRANSFER'" json:"payment_mode"` // Cara bayar: CASH, TRANSFER, GIRO

	// ── Nilai Transaksi ──────────────────────────────────────────────────
	TotalAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"total_amount"` // Total nilai pesanan (Qty * UnitPrice)

	// ── Status & Tracking ────────────────────────────────────────────────
	Status       string     `gorm:"type:varchar(20);not null;default:'DRAFT';index" json:"status"`
	ApprovedByID *uuid.UUID `gorm:"type:char(36);index" json:"approved_by_id"`
	ApprovedBy   *User      `gorm:"foreignKey:ApprovedByID" json:"approved_by,omitempty"`
	ApprovedAt   *time.Time `json:"approved_at"`
	CreatedByID  uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"` // Staf yang membuat PO
	CreatedBy    User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`

	// ── Notifikasi (Email / WA) ──────────────────────────────────────────
	NotificationMethod string     `gorm:"type:varchar(20);not null;default:'EMAIL'" json:"notification_method"`      // WHATSAPP, EMAIL, NONE
	NotificationStatus string     `gorm:"type:varchar(20);not null;default:'NONE';index" json:"notification_status"` // Status pengiriman: NONE, PENDING, SENT, FAILED
	SentAt             *time.Time `json:"sent_at"`                                                                   // Kapan PO berhasil terkirim ke supplier

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes         *string `gorm:"type:text" json:"notes"`          // Catatan internal
	SupplierNotes *string `gorm:"type:text" json:"supplier_notes"` // Catatan/instruksi khusus untuk supplier

	// ── Snapshot Data (Master Data at Time of Transaction) ───────────────
	SupplierName    string  `gorm:"type:varchar(150)" json:"supplier_name"`
	SupplierCode    string  `gorm:"type:varchar(20)" json:"supplier_code"`
	SupplierAddress *string `gorm:"type:text" json:"supplier_address"`
	StoreCode       string  `gorm:"type:varchar(20)" json:"store_code"`
	StoreName       string  `gorm:"type:varchar(150)" json:"store_name"`
	StoreAddress    *string `gorm:"type:text" json:"store_address"`
	WarehouseName   string  `gorm:"type:varchar(100)" json:"warehouse_name"`
	CreatedByName   string  `gorm:"type:varchar(100)" json:"created_by_name"`
	ApprovedByName  *string `gorm:"type:varchar(100)" json:"approved_by_name"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchaseOrderItem (Detail / Baris Pesanan)
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseOrderItem adalah baris detail barang dalam satu PO.
// Setiap baris merepresentasikan satu produk yang dipesan beserta qty, harga, dan UOM-nya.
type PurchaseOrderItem struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	PurchaseOrderID uuid.UUID     `gorm:"type:char(36);not null;index" json:"purchase_order_id"`
	PurchaseOrder   PurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`
	SeqNo           int           `gorm:"not null" json:"seq_no"` // Nomor urut baris (1, 2, 3, ...)

	// ── Produk & Satuan ──────────────────────────────────────────────────
	ProductID uuid.UUID `gorm:"type:char(36);not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UOMID     uuid.UUID `gorm:"type:char(36);not null" json:"uom_id"` // Satuan pesanan (bisa PCS, DUS, KG)
	UOM       UOM       `gorm:"foreignKey:UOMID" json:"uom,omitempty"`

	// ── Kuantitas ────────────────────────────────────────────────────────
	QtyOrdered  decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"qty_ordered"`   // Jumlah yang dipesan
	QtyReceived decimal.Decimal `gorm:"type:decimal(15,4);default:0" json:"qty_received"` // Jumlah yang sudah diterima (denormalisasi, di-update saat Goods Receipt)

	// ── Harga & Nilai ────────────────────────────────────────────────────
	UnitPrice decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"unit_price"` // Harga beli per satuan (gross)
	Subtotal  decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"subtotal"`   // QtyOrdered × UnitPrice

	// ── Traceability ─────────────────────────────────────────────────────
	ProductSupplierID *uuid.UUID       `gorm:"type:char(36);index" json:"product_supplier_id"` // Referensi ke kontrak supplier (opsional)
	ProductSupplier   *ProductSupplier `gorm:"foreignKey:ProductSupplierID" json:"product_supplier,omitempty"`
	PlanningID        *uuid.UUID       `gorm:"type:char(36);index" json:"planning_id"` // Referensi ke PurchaseOrderPlanning (jika PO dibuat dari rekomendasi sistem)

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes *string `gorm:"type:text" json:"notes"` // Catatan per baris item

	// ── Snapshot Data (Master Data at Time of Transaction) ───────────────
	ProductName string `gorm:"type:varchar(200)" json:"product_name"`
	ProductSKU  string `gorm:"type:varchar(50)" json:"product_sku"`
	UOMName     string `gorm:"type:varchar(50)" json:"uom_name"`
}

/*
Alur kerjanya nanti di Service Layer (Backend):

Manajer mengklik "Approve". API /approve dipanggil.
Database di-update: Status = "APPROVED" dan NotificationStatus = "PENDING".
Backend memasukkan tugas pengiriman WA/Email ke antrean (Message Queue, misalnya Redis/RabbitMQ/Kafka) di background.
Respons instan dikembalikan ke frontend: "PO Berhasil di-Approve!". (Tidak perlu menunggu WA terkirim).
Di background, Worker akan memproses antrean. Ia mengambil data PO, merakit PDF, dan memanggil API WhatsApp/Email.
Jika API berhasil membalas "200 OK", Worker meng-update database: NotificationStatus = "SENT" dan SentAt = WaktuSekarang.
Jika API WhatsApp down / gagal, status bisa di-update jadi "FAILED", dan di frontend bisa muncul tombol "Kirim Ulang (Resend)".
*/
