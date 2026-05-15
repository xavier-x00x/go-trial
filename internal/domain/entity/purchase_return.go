package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	PRStatusDraft  = "DRAFT"
	PRStatusPosted = "POSTED"
	PRStatusVoided = "VOIDED"
)

type PurchaseReturn struct {
	BaseModel

	ReturnNumber      string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"return_number"`
	ReturnDate        time.Time       `gorm:"not null" json:"return_date"`
	PurchaseInvoiceID uuid.UUID       `gorm:"type:char(36);not null;index" json:"purchase_invoice_id"`
	PurchaseInvoice   PurchaseInvoice `gorm:"foreignKey:PurchaseInvoiceID" json:"purchase_invoice,omitempty"`
	SupplierID        uuid.UUID       `gorm:"type:char(36);not null;index" json:"supplier_id"`
	Supplier          Supplier        `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	StoreID           uuid.UUID       `gorm:"type:char(36);not null;index" json:"store_id"`
	Store             Store           `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	WarehouseID       uuid.UUID       `gorm:"type:char(36);not null;index" json:"warehouse_id"`
	Warehouse         Warehouse       `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`

	Subtotal        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"subtotal"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"discount_amount"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"tax_amount"`
	GrandTotal      decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"grand_total"`
	RemainingAmount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"remaining_amount"`

	Status      string     `gorm:"type:varchar(10);not null;default:'DRAFT';index" json:"status"`
	Notes       *string    `gorm:"type:text" json:"notes"`
	CreatedByID uuid.UUID  `gorm:"type:char(36);not null;index" json:"created_by_id"`
	CreatedBy   User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	PostedByID  *uuid.UUID `gorm:"type:char(36);index" json:"posted_by_id"`
	PostedBy    *User      `gorm:"foreignKey:PostedByID" json:"posted_by,omitempty"`
	PostedAt    *time.Time `json:"posted_at"`

	Items []PurchaseReturnItem `gorm:"foreignKey:PurchaseReturnID" json:"items,omitempty"`
}

type PurchaseReturnItem struct {
	BaseModel

	PurchaseReturnID    uuid.UUID          `gorm:"type:char(36);not null;index" json:"purchase_return_id"`
	SeqNo               int                `gorm:"not null" json:"seq_no"`
	PurchaseInvoiceItemID uuid.UUID        `gorm:"type:char(36);not null;index" json:"purchase_invoice_item_id"`
	PurchaseInvoiceItem   PurchaseInvoiceItem `gorm:"foreignKey:PurchaseInvoiceItemID" json:"purchase_invoice_item,omitempty"`
	ProductID           uuid.UUID          `gorm:"type:char(36);not null;index" json:"product_id"`
	Product             Product            `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UOMID               uuid.UUID          `gorm:"type:char(36);not null;index" json:"uom_id"`
	UOM                 UOM                `gorm:"foreignKey:UOMID" json:"uom,omitempty"`

	QtyReturn      decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"qty_return"`
	UnitPrice      decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"unit_price"`
	Discount1Pct   decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"discount_1_pct"`
	Discount2Pct   decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"discount_2_pct"`
	Discount3Pct   decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"discount_3_pct"`
	DiscountAmount decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"discount_amount"`
	TaxPct         decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"tax_pct"`
	TaxAmount      decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"tax_amount"`
	Subtotal       decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"subtotal"`
}
