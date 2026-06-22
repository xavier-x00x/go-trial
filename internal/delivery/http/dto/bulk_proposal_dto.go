package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ──────────────────────────────────────────────────────────────────────────────
// DTO Khusus Use-Case: Bulk Link Product-Supplier
// ──────────────────────────────────────────────────────────────────────────────

// BulkCreateProductSupplierProposalRequest adalah DTO convenience untuk
// use-case spesifik: menghubungkan banyak barang ke satu supplier sekaligus.
//
// Di Service Layer, DTO ini akan di-transform menjadi:
//   - 1 MasterDataProposal Header (EntityType: PRODUCT_SUPPLIER, ActionType: CREATE)
//   - N MasterDataProposalItem (masing-masing berisi PayloadJSON dari setiap item)
type BulkCreateProductSupplierProposalRequest struct {
	SupplierID uuid.UUID                      `json:"supplier_id" validate:"required"`
	Reason     string                         `json:"reason"`
	Items      []ProductSupplierLinkDetailDTO `json:"items" validate:"required,min=1,dive"`
}

// ProductSupplierLinkDetailDTO berisi detail spesifik kontrak barang dengan supplier.
type ProductSupplierLinkDetailDTO struct {
	ProductID           uuid.UUID       `json:"product_id" validate:"required"`
	StoreID             *uuid.UUID      `json:"store_id"`
	SupplierSKU         *string         `json:"supplier_sku"`
	IsPrimary           bool            `json:"is_primary"`
	IsConsignment       bool            `json:"is_consignment"`
	IsReturnable        bool            `json:"is_returnable"`
	DefaultLeadTimeDays int             `json:"default_lead_time_days"`
	PurchaseUOMID       *uuid.UUID      `json:"purchase_uom_id"`
	OfferedPrice        decimal.Decimal `json:"offered_price" validate:"required"`
	MinOrderQty         decimal.Decimal `json:"min_order_qty" validate:"required"`
}

// BulkProposalResponse adalah response untuk operasi bulk proposal.
type BulkProposalResponse struct {
	ReferenceNumbers []string                              `json:"reference_numbers"`
	TotalCount    int                                 `json:"total_count"`
	SuccessCount  int                                 `json:"success_count"`
	FailedCount  int                                 `json:"failed_count"`
	Proposals    []MasterDataProposalDetailResponse `json:"proposals"`
}
