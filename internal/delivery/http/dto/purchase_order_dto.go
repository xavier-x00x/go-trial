package dto

import (
	"time"

	"github.com/google/uuid"
)

// ──────────────────────────────────────────────────────────────────────────────
// Request DTOs
// ──────────────────────────────────────────────────────────────────────────────

type CreatePurchaseOrderRequest struct {
	SupplierID       uuid.UUID                        `json:"supplier_id" validate:"required"`
	StoreID          uuid.UUID                        `json:"store_id" validate:"required"`
	WarehouseID      uuid.UUID                        `json:"warehouse_id" validate:"required"`
	ReferenceNo      *string                          `json:"reference_no"`
	OrderDate        time.Time                        `json:"order_date"`
	ExpectedDelivery *time.Time                       `json:"expected_delivery"`
	PaymentTermDays  int                              `json:"payment_term_days"`
	PaymentMode      string                           `json:"payment_mode"`
	Notes            *string                          `json:"notes"`
	SupplierNotes    *string                          `json:"supplier_notes"`
	Items            []CreatePurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// UpdatePurchaseOrderRequest digunakan untuk mengubah PO yang masih DRAFT.
type UpdatePurchaseOrderRequest struct {
	ReferenceNo      *string                          `json:"reference_no"`
	ExpectedDelivery *time.Time                       `json:"expected_delivery"`
	PaymentTermDays  int                              `json:"payment_term_days"`
	PaymentMode      string                           `json:"payment_mode"`
	Notes            *string                          `json:"notes"`
	SupplierNotes    *string                          `json:"supplier_notes"`
	Items            []CreatePurchaseOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseOrderItemRequest struct {
	ProductID         uuid.UUID  `json:"product_id" validate:"required"`
	UOMID             uuid.UUID  `json:"uom_id" validate:"required"`
	QtyOrdered        float64    `json:"qty_ordered" validate:"required,min=0.0001"`
	UnitPrice         float64    `json:"unit_price" validate:"required,min=0"`
	ProductSupplierID *uuid.UUID `json:"product_supplier_id"`
	PlanningID        *uuid.UUID `json:"planning_id"`
	Notes             *string    `json:"notes"`
}

type CreatePOFromPlanningRequest struct {
	StoreID      uuid.UUID `json:"store_id" validate:"required"`
	SupplierID  uuid.UUID `json:"supplier_id" validate:"required"`
	WarehouseID uuid.UUID `json:"warehouse_id" validate:"required"`
	PlanningIDs []string  `json:"planning_ids" validate:"required,min=1"`
	OrderDate   time.Time `json:"order_date"`
	Notes       *string  `json:"notes"`
}

type SubmitPORequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type ApprovePORequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type CancelPORequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	Reason string   `json:"reason"`
}

type BulkCreatePOFromPlanningRequest struct {
	StoreID      uuid.UUID `json:"store_id" validate:"required"`
	WarehouseID  uuid.UUID `json:"warehouse_id" validate:"required"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Response DTOs
// ──────────────────────────────────────────────────────────────────────────────

type PurchaseOrderListResponse struct {
	ID              uuid.UUID       `json:"id"`
	PONumber       string          `json:"po_number"`
	ReferenceNo    *string         `json:"reference_no,omitempty"`
	SupplierID     uuid.UUID       `json:"supplier_id"`
	SupplierName   string          `json:"supplier_name"`
	SupplierCode   string          `json:"supplier_code"`
	StoreID        uuid.UUID       `json:"store_id"`
	StoreName      string          `json:"store_name"`
	WarehouseID    uuid.UUID       `json:"warehouse_id"`
	WarehouseName  string          `json:"warehouse_name"`
	OrderDate      time.Time       `json:"order_date"`
	ExpectedDelivery *time.Time   `json:"expected_delivery,omitempty"`
	TotalAmount    float64         `json:"total_amount"`
	Status        string          `json:"status"`
	ApprovedByID   *uuid.UUID     `json:"approved_by_id,omitempty"`
	ApprovedAt    *time.Time      `json:"approved_at,omitempty"`
	CreatedByID   uuid.UUID      `json:"created_by_id"`
	CreatedByName string         `json:"created_by_name"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type PurchaseOrderDetailResponse struct {
	ID              uuid.UUID                     `json:"id"`
	PONumber       string                       `json:"po_number"`
	ReferenceNo   *string                     `json:"reference_no,omitempty"`
	Supplier      SupplierResponse             `json:"supplier"`
	Store         StoreResponse             `json:"store"`
	Warehouse     WarehouseResponse          `json:"warehouse"`
	OrderDate     time.Time                  `json:"order_date"`
	ExpectedDelivery *time.Time              `json:"expected_delivery,omitempty"`
	PaymentTermDays int                      `json:"payment_term_days"`
	PaymentMode   string                    `json:"payment_mode"`
	TotalAmount      float64                    `json:"total_amount"`
	Status       string                     `json:"status"`
	ApprovedByID *uuid.UUID                `json:"approved_by_id,omitempty"`
	ApprovedAt  *time.Time                 `json:"approved_at,omitempty"`
	ApprovedByName *string                 `json:"approved_by_name,omitempty"`
	CreatedByID  uuid.UUID                  `json:"created_by_id"`
	CreatedByName string                    `json:"created_by_name"`
	Notes            *string                    `json:"notes,omitempty"`
	SupplierNotes    *string                    `json:"supplier_notes,omitempty"`
	SupplierNameSnapshot      string          `json:"supplier_name"`
	SupplierCodeSnapshot      string          `json:"supplier_code"`
	SupplierAddressSnapshot *string           `json:"supplier_address,omitempty"`
	StoreCodeSnapshot        string          `json:"store_code"`
	StoreNameSnapshot         string          `json:"store_name"`
	StoreAddressSnapshot   *string           `json:"store_address,omitempty"`
	WarehouseNameSnapshot   string          `json:"warehouse_name"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedAt    time.Time                  `json:"updated_at"`
	Items        []PurchaseOrderItemResponse `json:"items"`
}

type PurchaseOrderItemResponse struct {
	ID                uuid.UUID   `json:"id"`
	SeqNo             int         `json:"seq_no"`
	ProductID         uuid.UUID   `json:"product_id"`
	ProductSKU        string      `json:"product_sku"`
	ProductName       string      `json:"product_name"`
	UOMID            uuid.UUID   `json:"uom_id"`
	UOMName          string      `json:"uom_name"`
	QtyOrdered       float64    `json:"qty_ordered"`
	QtyReceived      float64    `json:"qty_received"`
	UnitPrice           float64    `json:"unit_price"`
	Subtotal            float64    `json:"subtotal"`
	ProductSupplierID *uuid.UUID `json:"product_supplier_id,omitempty"`
	PlanningID       *uuid.UUID `json:"planning_id,omitempty"`
	Notes            *string   `json:"notes,omitempty"`
}

type PurchaseOrderResponse struct {
	ID              uuid.UUID       `json:"id"`
	PONumber       string          `json:"po_number"`
	SupplierID     uuid.UUID       `json:"supplier_id"`
	SupplierName   string          `json:"supplier_name"`
	StoreID        uuid.UUID       `json:"store_id"`
	StoreName      string          `json:"store_name"`
	GrandTotal     float64        `json:"grand_total"`
	Status        string          `json:"status"`
	OrderDate      time.Time       `json:"order_date"`
}