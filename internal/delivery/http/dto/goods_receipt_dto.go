package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request: Create Goods Receipt
// ──────────────────────────────────────────────────────────────────────────────

type CreateGoodsReceiptRequest struct {
	PurchaseOrderID uuid.UUID                     `json:"purchase_order_id" validate:"required"`
	WarehouseID     uuid.UUID                     `json:"warehouse_id" validate:"required"`
	ReceiptDate     time.Time                     `json:"receipt_date" validate:"required"`
	DeliveryNoteNo  *string                       `json:"delivery_note_no" validate:"omitempty,max=50"`
	Notes           *string                       `json:"notes"`
	Items           []CreateGoodsReceiptItemInput `json:"items" validate:"required,min=1,dive"`
	OverridePIN     *string                       `json:"override_pin"`
}

// UpdateGoodsReceiptRequest digunakan untuk mengubah GR yang masih DRAFT.
type UpdateGoodsReceiptRequest struct {
	ReceiptDate    time.Time                     `json:"receipt_date" validate:"required"`
	DeliveryNoteNo *string                       `json:"delivery_note_no" validate:"omitempty,max=50"`
	Notes          *string                       `json:"notes"`
	Items          []CreateGoodsReceiptItemInput `json:"items" validate:"required,min=1,dive"`
	OverridePIN    *string                       `json:"override_pin"`
}

type CreateGoodsReceiptItemInput struct {
	PurchaseOrderItemID uuid.UUID       `json:"purchase_order_item_id" validate:"required"`
	ProductID           uuid.UUID       `json:"product_id" validate:"required"`
	UOMID               uuid.UUID       `json:"uom_id" validate:"required"`
	QtyReceived         decimal.Decimal `json:"qty_received" validate:"required"`
	QtyRejected         decimal.Decimal `json:"qty_rejected"`
	UnitPrice           decimal.Decimal `json:"unit_price" validate:"required,min=0"`
	RejectReason        *string         `json:"reject_reason"`
	Notes               *string         `json:"notes"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Request: Confirm (Kepala Gudang mengkonfirmasi penerimaan → stok bertambah)
// ──────────────────────────────────────────────────────────────────────────────

type ConfirmGoodsReceiptRequest struct {
	Notes *string `json:"notes"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Request: Mode Langsung (Terima Barang + Catat Faktur sekaligus)
// ──────────────────────────────────────────────────────────────────────────────

// CreateGoodsReceiptWithInvoiceRequest digunakan saat barang dan faktur
// datang bersamaan (kasus umum di retail). Dalam satu kali submit:
//   - Sistem membuat GoodsReceipt + auto-confirm → Stok bertambah, HPP dihitung
//   - Sistem membuat PurchaseInvoice + auto-post → Hutang tercatat, jurnal dibuat
//
// Ini adalah shortcut UX. Di database tetap tersimpan sebagai 2 dokumen terpisah
// (GR dan PI) agar audit trail dan laporan keuangan tetap bersih.
type CreateGoodsReceiptWithInvoiceRequest struct {
	// ── Data Penerimaan (GR) ──────────────────────────────────────────────
	PurchaseOrderID uuid.UUID `json:"purchase_order_id" validate:"required"`
	WarehouseID     uuid.UUID `json:"warehouse_id" validate:"required"`
	ReceiptDate     time.Time `json:"receipt_date" validate:"required"`
	DeliveryNoteNo  *string   `json:"delivery_note_no" validate:"omitempty,max=50"`

	// ── Data Faktur (PI) ─────────────────────────────────────────────────
	SupplierInvoiceNumber string          `json:"supplier_invoice_number" validate:"required,max=50"`
	InvoiceDate           time.Time       `json:"invoice_date" validate:"required"`
	APAccountID           uuid.UUID       `json:"ap_account_id" validate:"required"`
	PaymentTermDays       int             `json:"payment_term_days" validate:"min=0"`
	DiscountAmount        decimal.Decimal `json:"discount_amount" validate:"min=0"`

	// ── Items (berlaku untuk GR dan PI sekaligus) ────────────────────────
	Notes *string                       `json:"notes"`
	Items []CreateGoodsReceiptItemInput `json:"items" validate:"required,min=1,dive"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: List (Header Only)
// ──────────────────────────────────────────────────────────────────────────────

type GoodsReceiptListResponse struct {
	ID              uuid.UUID `json:"id"`
	GRNumber        string    `json:"gr_number"`
	PurchaseOrderID uuid.UUID `json:"purchase_order_id"`
	PONumber        string    `json:"po_number"`
	WarehouseID     uuid.UUID `json:"warehouse_id"`
	WarehouseName   string    `json:"warehouse_name"`
	ReceiptDate     time.Time `json:"receipt_date"`
	DeliveryNoteNo  *string   `json:"delivery_note_no,omitempty"`
	Status          string    `json:"status"`
	ReceivedByID    uuid.UUID `json:"received_by_id"`
	SupplierName    string    `json:"supplier_name"`
	StoreName       string    `json:"store_name"`
	CreatedAt       time.Time `json:"created_at"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response: Detail (Header + Items)
// ──────────────────────────────────────────────────────────────────────────────

type GoodsReceiptDetailResponse struct {
	ID              uuid.UUID                  `json:"id"`
	GRNumber        string                     `json:"gr_number"`
	PurchaseOrderID uuid.UUID                  `json:"purchase_order_id"`
	PONumber        string                     `json:"po_number"`
	WarehouseID     uuid.UUID                  `json:"warehouse_id"`
	WarehouseName   string                     `json:"warehouse_name"`
	ReceiptDate     time.Time                  `json:"receipt_date"`
	DeliveryNoteNo  *string                    `json:"delivery_note_no,omitempty"`
	Status          string                     `json:"status"`
	ReceivedByID    uuid.UUID                  `json:"received_by_id"`
	ConfirmedByID   *uuid.UUID                 `json:"confirmed_by_id,omitempty"`
	ConfirmedAt     *time.Time                 `json:"confirmed_at,omitempty"`
	Notes                  *string                    `json:"notes,omitempty"`
	IsOverReceivedOverride bool                       `json:"is_over_received_override"`
	OverrideApprovedByID   *uuid.UUID                 `json:"override_approved_by_id,omitempty"`
	CreatedAt              time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	SupplierName    string                     `json:"supplier_name"`
	SupplierCode    string                     `json:"supplier_code"`
	SupplierAddress *string                    `json:"supplier_address,omitempty"`
	StoreName       string                     `json:"store_name"`
	Subtotal        decimal.Decimal            `json:"subtotal"`
	DiscountAmount  decimal.Decimal            `json:"discount_amount"`
	TaxAmount       decimal.Decimal            `json:"tax_amount"`
	FreightAmount   decimal.Decimal            `json:"freight_amount"`
	OtherCostAmount decimal.Decimal            `json:"other_cost_amount"`
	GrandTotal      decimal.Decimal            `json:"grand_total"`
	IsTaxInclusive  bool                       `json:"is_tax_inclusive"`
	Items           []GoodsReceiptItemResponse `json:"items"`
}

type GoodsReceiptItemResponse struct {
	ID                  uuid.UUID       `json:"id"`
	SeqNo               int             `json:"seq_no"`
	PurchaseOrderItemID uuid.UUID       `json:"purchase_order_item_id"`
	ProductID           uuid.UUID       `json:"product_id"`
	ProductName         string          `json:"product_name"`
	ProductSKU          string          `json:"product_sku"`
	UOMID               uuid.UUID       `json:"uom_id"`
	UOMCode             string          `json:"uom_code"`
	QtyOrdered          decimal.Decimal `json:"qty_ordered"`
	QtyReceived         decimal.Decimal `json:"qty_received"`
	QtyRejected         decimal.Decimal `json:"qty_rejected"`
	UnitPrice           decimal.Decimal `json:"unit_price"`
	Discount1Pct        decimal.Decimal `json:"discount_1_pct"`
	Discount2Pct        decimal.Decimal `json:"discount_2_pct"`
	Discount3Pct        decimal.Decimal `json:"discount_3_pct"`
	DiscountAmount      decimal.Decimal `json:"discount_amount"`
	TotalDiscountAmount decimal.Decimal `json:"total_discount_amount"`
	TaxPct              decimal.Decimal `json:"tax_pct"`
	TaxAmount           decimal.Decimal `json:"tax_amount"`
	LandedCostAmount    decimal.Decimal `json:"landed_cost_amount"`
	NetUnitPrice        decimal.Decimal `json:"net_unit_price"`
	RejectReason        *string         `json:"reject_reason,omitempty"`
	Notes               *string         `json:"notes,omitempty"`
}
