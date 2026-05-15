package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Status Purchase Invoice
// ──────────────────────────────────────────────────────────────────────────────

const (
	PurchaseInvoiceStatusDraft          = "DRAFT"          // Faktur baru diinput, belum diverifikasi
	PurchaseInvoiceStatusSubmitted     = "SUBMITTED"      // Faktur telah disubmit ke sistem
	PurchaseInvoiceStatusVerified      = "VERIFIED"       // Sudah dicocokkan dengan PO & GR (3-Way Match)
	PurchaseInvoiceStatusPosted        = "POSTED"         // Dijurnal → Hutang resmi tercatat di Akuntansi
	PurchaseInvoiceStatusPartiallyPaid = "PARTIALLY_PAID" // Sebagian hutang sudah dibayar
	PurchaseInvoiceStatusPaid          = "PAID"           // Lunas
	PurchaseInvoiceStatusCancelled     = "CANCELLED"      // Dibatalkan (hanya dari DRAFT/VERIFIED)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchaseInvoice (Header)
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseInvoice adalah dokumen faktur pembelian dari supplier.
// Dokumen ini merupakan dasar pencatatan hutang (Account Payable) ke supplier.
//
// Siklus Hidup:
//
//	DRAFT → SUBMITTED → VERIFIED → POSTED → PARTIALLY_PAID → PAID
//	                  (atau langsung POSTED jika 3-Way Match otomatis)
//	Bisa CANCELLED kapan saja sebelum POSTED.
//
// Efek saat Status berubah ke POSTED:
//   - Jurnal Akuntansi dicatat: Debit Persediaan, Kredit Hutang Usaha (APAccountID)
//   - MonthlyAPBalance.TotalCredit bertambah sebesar GrandTotal
//   - PurchaseOrder.Status → CLOSED (jika seluruh qty sudah di-invoice)
//   - Jatuh tempo pembayaran mulai aktif dihitung
//
// 3-Way Matching:
//
//	PO (apa yang dipesan) vs GR (apa yang diterima) vs Invoice (apa yang ditagih)
//	Staf Finance harus memastikan ketiganya cocok sebelum POSTED.
type PurchaseInvoice struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	InvoiceNumber         string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_number"` // Nomor internal faktur (cth: PI/2026/04/0001)
	SupplierInvoiceNumber string        `gorm:"type:varchar(50);not null" json:"supplier_invoice_number"`    // Nomor faktur ASLI dari supplier (tercetak di kertas faktur)
	ReferenceNo           *string       `gorm:"type:varchar(50)" json:"reference_no"`                        // Nomor referensi eksternal
	PurchaseOrderID       uuid.UUID     `gorm:"type:char(36);not null;index" json:"purchase_order_id"`       // Referensi ke PO induk
	PurchaseOrder         PurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`

	// ── Pihak Terkait ────────────────────────────────────────────────────
	SupplierID  uuid.UUID      `gorm:"type:char(36);not null;index" json:"supplier_id"`
	Supplier    Supplier       `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	StoreID     uuid.UUID      `gorm:"type:char(36);not null;index" json:"store_id"` //Cabang yang terkait
	Store       Store          `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	WarehouseID uuid.UUID      `gorm:"type:char(36);not null;index" json:"warehouse_id"` // Gudang tujuan
	Warehouse   Warehouse      `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	APAccountID        uuid.UUID      `gorm:"type:char(36);not null;index" json:"ap_account_id"` // Akun Hutang Usaha untuk penjurnalan
	APAccount          ChartOfAccount `gorm:"foreignKey:APAccountID" json:"ap_account,omitempty"`
	InventoryAccountID uuid.UUID      `gorm:"type:char(36);not null;index" json:"inventory_account_id"` // Akun Persediaan untuk penjurnalan
	InventoryAccount   ChartOfAccount `gorm:"foreignKey:InventoryAccountID" json:"inventory_account,omitempty"`

	// ── Tanggal & Jadwal ─────────────────────────────────────────────────
	InvoiceDate      time.Time  `gorm:"not null" json:"invoice_date"`   // Tanggal faktur diterbitkan oleh supplier
	ReceivedDate     time.Time  `gorm:"not null" json:"received_date"`  // Tanggal faktur diterima oleh staf Finance kita
	DueDate          time.Time  `gorm:"not null;index" json:"due_date"` // Tanggal jatuh tempo pembayaran (InvoiceDate + PaymentTermDays)
	ExpectedDelivery *time.Time `json:"expected_delivery"`              // Estimasi tanggal pembayaran

	// ── Syarat Pembayaran ────────────────────────────────────────────────
	PaymentTermDays int    `gorm:"default:0" json:"payment_term_days"`
	PaymentMode     string `gorm:"type:varchar(20);default:'TRANSFER'" json:"payment_mode"`

	// ── Nilai Transaksi ──────────────────────────────────────────────────
	Subtotal        decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"subtotal"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"discount_amount"` // Potongan tambahan di faktur
	TaxAmount       decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"tax_amount"`      // Total PPN
	FreightAmount   decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"freight_amount"`
	OtherCostAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"other_cost_amount"`
	GrandTotal      decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"grand_total"`      // Nilai hutang yang harus dibayar
	PaidAmount      decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"paid_amount"`      // Total yang sudah dibayar (di-update saat pembayaran)
	RemainingAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"remaining_amount"` // Sisa hutang (GrandTotal - PaidAmount)
	IsTaxInclusive  bool            `gorm:"default:false" json:"is_tax_inclusive"`

	// ── Status & Tracking ────────────────────────────────────────────────
	Status       string     `gorm:"type:varchar(20);not null;default:'DRAFT';index" json:"status"`
	VerifiedByID *uuid.UUID `gorm:"type:char(36);index" json:"verified_by_id"` // Staf Finance yang memverifikasi 3-Way Match
	VerifiedBy   *User      `gorm:"foreignKey:VerifiedByID" json:"verified_by,omitempty"`
	VerifiedAt   *time.Time `json:"verified_at"`
	PostedByID   *uuid.UUID `gorm:"type:char(36);index" json:"posted_by_id"` // Manajer Akuntansi yang mem-posting jurnal
	PostedBy     *User      `gorm:"foreignKey:PostedByID" json:"posted_by,omitempty"`
	PostedAt     *time.Time `json:"posted_at"`
	CreatedByID  uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"`
	CreatedBy    User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`

	

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes *string `gorm:"type:text" json:"notes"`

	// ── Snapshot Data (Master Data at Time of Transaction) ───────────────
	SupplierCode    string  `gorm:"type:varchar(20)" json:"supplier_code"`
	SupplierName    string  `gorm:"type:varchar(150)" json:"supplier_name"`
	SupplierAddress *string `gorm:"type:text" json:"supplier_address"`
	StoreCode       string  `gorm:"type:varchar(20)" json:"store_code"`
	StoreName       string  `gorm:"type:varchar(150)" json:"store_name"`
	StoreAddress    *string `gorm:"type:text" json:"store_address"`
	WarehouseName   string  `gorm:"type:varchar(100)" json:"warehouse_name"`
	VerifiedByName  *string `gorm:"type:varchar(100)" json:"verified_by_name"`
	PostedByName    *string `gorm:"type:varchar(100)" json:"posted_by_name"`
	CreatedByName   string  `gorm:"type:varchar(100)" json:"created_by_name"`
	ApprovedByName  *string `gorm:"type:varchar(100)" json:"approved_by_name"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []PurchaseInvoiceItem `gorm:"foreignKey:PurchaseInvoiceID" json:"items,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: PurchaseInvoiceItem (Detail / Baris Faktur)
// ──────────────────────────────────────────────────────────────────────────────

// PurchaseInvoiceItem adalah baris detail barang dalam faktur pembelian.
// Staf Finance mencocokkan setiap baris ini dengan PO Item dan GR Item.
type PurchaseInvoiceItem struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	PurchaseInvoiceID uuid.UUID       `gorm:"type:char(36);not null;index" json:"purchase_invoice_id"`
	PurchaseInvoice   PurchaseInvoice `gorm:"foreignKey:PurchaseInvoiceID" json:"purchase_invoice,omitempty"`
	SeqNo             int             `gorm:"not null" json:"seq_no"`

	// ── Referensi 3-Way Matching ─────────────────────────────────────────
	PurchaseOrderItemID *uuid.UUID         `gorm:"type:char(36);index" json:"purchase_order_item_id"` // Link ke baris PO
	PurchaseOrderItem   *PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderItemID" json:"purchase_order_item,omitempty"`
	GoodsReceiptItemID  *uuid.UUID         `gorm:"type:char(36);index" json:"goods_receipt_item_id"` // Link ke baris GR
	GoodsReceiptItem    *GoodsReceiptItem  `gorm:"foreignKey:GoodsReceiptItemID" json:"goods_receipt_item,omitempty"`

	// ── Produk & Satuan ──────────────────────────────────────────────────
	ProductID uuid.UUID `gorm:"type:char(36);not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UOMID     uuid.UUID `gorm:"type:char(36);not null" json:"uom_id"`
	UOM       UOM       `gorm:"foreignKey:UOMID" json:"uom,omitempty"`

	// ── Kuantitas & Harga ────────────────────────────────────────────────
	QtyInvoiced         decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"qty_invoiced"` // Jumlah yang ditagihkan supplier
	UnitPrice           decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"unit_price"`   // Harga per satuan di faktur (gross)
	Discount1Pct        decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_1_pct"`
	Discount2Pct        decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_2_pct"`
	Discount3Pct        decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_3_pct"`
	DiscountAmount      decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"discount_amount"`
	TotalDiscountAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"total_discount_amount"`
	TaxPct              decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"tax_pct"`
	TaxAmount           decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"tax_amount"`
	LandedCostAmount    decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"landed_cost_amount"`
	Subtotal            decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"subtotal"`        // (QtyInvoiced × UnitPrice) - TotalDiscountAmount
	NetUnitPrice        decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"net_unit_price"` // (Subtotal + TaxAmount + LandedCostAmount) / QtyInvoiced

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes *string `gorm:"type:text" json:"notes"`

	// ── Snapshot Data (Master Data at Time of Transaction) ───────────────
	ProductName string `gorm:"type:varchar(200)" json:"product_name"`
	ProductSKU  string `gorm:"type:varchar(50)" json:"product_sku"`
	UOMName     string `gorm:"type:varchar(50)" json:"uom_name"`
}
