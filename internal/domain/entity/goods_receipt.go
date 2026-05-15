package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Konstanta Status Goods Receipt
// ──────────────────────────────────────────────────────────────────────────────

const (
	GRStatusDraft     = "DRAFT"     // GR baru dibuat, staf gudang sedang menghitung fisik barang
	GRStatusConfirmed = "CONFIRMED" // Dikonfirmasi → stok otomatis bertambah di InventoryStock
	GRStatusCancelled = "CANCELLED" // Dibatalkan (hanya bisa dari DRAFT)
)

// ──────────────────────────────────────────────────────────────────────────────
// Entity: GoodsReceipt (Header)
// ──────────────────────────────────────────────────────────────────────────────

// GoodsReceipt adalah dokumen Bukti Terima Barang (BTB) yang dibuat saat
// barang fisik tiba di gudang dari supplier.
//
// Efek saat Status berubah ke CONFIRMED:
//   - InventoryStock.Quantity bertambah sebesar QtyReceived (dikonversi ke Base UOM)
//   - InventoryStock.AverageBuyPrice dihitung ulang (Weighted Average)
//   - PurchaseOrderItem.QtyReceived di-update
//   - PurchaseOrder.Status berubah ke PARTIALLY_RECEIVED atau RECEIVED
//
// Satu PO bisa memiliki banyak GoodsReceipt (pengiriman bertahap/parsial).
type GoodsReceipt struct {
	BaseModel

	// ── Identifikasi Dokumen ─────────────────────────────────────────────
	GRNumber        string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"gr_number"` // Nomor BTB unik (cth: GR/2026/04/0001)
	PurchaseOrderID uuid.UUID     `gorm:"type:char(36);not null;index" json:"purchase_order_id"`  // Referensi ke PO induk
	PurchaseOrder   PurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`

	// ── Lokasi Penerimaan ────────────────────────────────────────────────
	WarehouseID uuid.UUID `gorm:"type:char(36);not null;index" json:"warehouse_id"` // Gudang yang menerima barang
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`

	// ── Tanggal ──────────────────────────────────────────────────────────
	ReceiptDate    time.Time `gorm:"not null" json:"receipt_date"`             // Tanggal fisik barang diterima
	DeliveryNoteNo *string   `gorm:"type:varchar(50)" json:"delivery_note_no"` // Nomor Surat Jalan dari supir/supplier

	// ── Status & Tracking ────────────────────────────────────────────────
	Status        string     `gorm:"type:varchar(10);not null;default:'DRAFT';index" json:"status"`
	ReceivedByID  uuid.UUID  `gorm:"type:char(36);not null;index" json:"received_by_id"` // Staf gudang yang menerima dan menghitung
	ReceivedBy    User       `gorm:"foreignKey:ReceivedByID" json:"received_by,omitempty"`
	ConfirmedByID *uuid.UUID `gorm:"type:char(36);index" json:"confirmed_by_id"` // Kepala Gudang / Manajer yang mengkonfirmasi
	ConfirmedBy   *User      `gorm:"foreignKey:ConfirmedByID" json:"confirmed_by,omitempty"`
	ConfirmedAt   *time.Time `json:"confirmed_at"`

	// ── Catatan ──────────────────────────────────────────────────────────
	Notes *string `gorm:"type:text" json:"notes"`

	// ── Nilai Transaksi (Snapshot from PO) ───────────────────────────────
	Subtotal        decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"subtotal"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"discount_amount"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"tax_amount"`
	FreightAmount   decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"freight_amount"`
	OtherCostAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"other_cost_amount"`
	GrandTotal      decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"grand_total"`
	IsTaxInclusive  bool            `gorm:"default:false" json:"is_tax_inclusive"`

	// ── Override Kuantitas ───────────────────────────────────────────────
	IsOverReceivedOverride bool       `gorm:"default:false" json:"is_over_received_override"`
	OverrideApprovedByID   *uuid.UUID `gorm:"type:char(36);index" json:"override_approved_by_id"`
	OverrideApprovedBy     *User      `gorm:"foreignKey:OverrideApprovedByID" json:"override_approved_by,omitempty"`

	// ── Snapshot Data (Master Data at Time of Transaction) ───────────────
	SupplierCode    string  `gorm:"type:varchar(20)" json:"supplier_code"`
	SupplierName    string  `gorm:"type:varchar(150)" json:"supplier_name"`
	SupplierAddress *string `gorm:"type:text" json:"supplier_address"`
	StoreCode       string  `gorm:"type:varchar(20)" json:"store_code"`
	StoreName       string  `gorm:"type:varchar(150)" json:"store_name"`
	StoreAddress    *string `gorm:"type:text" json:"store_address"`
	WarehouseName   string  `gorm:"type:varchar(100)" json:"warehouse_name"`
	CreatedByName   string  `gorm:"type:varchar(100)" json:"created_by_name"`
	ApprovedByName  *string `gorm:"type:varchar(100)" json:"approved_by_name"`
	ReceivedByName  string  `gorm:"type:varchar(100)" json:"received_by_name"`
	ConfirmedByName *string `gorm:"type:varchar(100)" json:"confirmed_by_name"`

	// ── Relasi Detail ────────────────────────────────────────────────────
	Items []GoodsReceiptItem `gorm:"foreignKey:GoodsReceiptID" json:"items,omitempty"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Entity: GoodsReceiptItem (Detail / Baris Penerimaan)
// ──────────────────────────────────────────────────────────────────────────────

// GoodsReceiptItem adalah baris detail barang yang diterima per item.
// Staf gudang mengisi QtyReceived (barang bagus) dan QtyRejected (barang rusak/salah).
type GoodsReceiptItem struct {
	BaseModel

	// ── Relasi Header ────────────────────────────────────────────────────
	GoodsReceiptID uuid.UUID    `gorm:"type:char(36);not null;index" json:"goods_receipt_id"`
	GoodsReceipt   GoodsReceipt `gorm:"foreignKey:GoodsReceiptID" json:"goods_receipt,omitempty"`
	SeqNo          int          `gorm:"not null" json:"seq_no"`

	// ── Referensi ke PO Item ─────────────────────────────────────────────
	PurchaseOrderItemID uuid.UUID         `gorm:"type:char(36);not null;index" json:"purchase_order_item_id"` // Link ke baris PO yang sedang diterima
	PurchaseOrderItem   PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderItemID" json:"purchase_order_item,omitempty"`

	// ── Produk & Satuan ──────────────────────────────────────────────────
	ProductID uuid.UUID `gorm:"type:char(36);not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UOMID     uuid.UUID `gorm:"type:char(36);not null" json:"uom_id"` // Satuan penerimaan (harus sama dengan UOM di PO Item)
	UOM       UOM       `gorm:"foreignKey:UOMID" json:"uom,omitempty"`

	// ── Kuantitas ────────────────────────────────────────────────────────
	QtyOrdered  decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"qty_ordered"`   // Snapshot qty dari PO (untuk referensi cepat saat cek fisik)
	QtyReceived decimal.Decimal `gorm:"type:decimal(15,4);not null" json:"qty_received"`  // Jumlah barang yang diterima dalam kondisi baik
	QtyRejected decimal.Decimal `gorm:"type:decimal(15,4);default:0" json:"qty_rejected"` // Jumlah barang ditolak (rusak, expired, salah kirim)

	// ── Harga & Nilai (Snapshot from PO for HPP calculation) ─────────────
	UnitPrice           decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"unit_price"`
	Discount1Pct        decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_1_pct"`
	Discount2Pct        decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_2_pct"`
	Discount3Pct        decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_3_pct"`
	DiscountAmount      decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"discount_amount"`
	TotalDiscountAmount decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"total_discount_amount"`
	TaxPct              decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"tax_pct"`
	TaxAmount           decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"tax_amount"`
	LandedCostAmount    decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"landed_cost_amount"`
	NetUnitPrice        decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"net_unit_price"`

	// ── Catatan ──────────────────────────────────────────────────────────
	RejectReason *string `gorm:"type:text" json:"reject_reason"` // Alasan penolakan (wajib diisi jika QtyRejected > 0)
	Notes        *string `gorm:"type:text" json:"notes"`

	// ── Snapshot Data (Master Data at Time of Transaction) ───────────────
	ProductName string `gorm:"type:varchar(200)" json:"product_name"`
	ProductSKU  string `gorm:"type:varchar(50)" json:"product_sku"`
	UOMName     string `gorm:"type:varchar(50)" json:"uom_name"`
}
