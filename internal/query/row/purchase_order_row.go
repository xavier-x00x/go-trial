package row

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PurchaseOrderListRow struct {
	ID               uuid.UUID  `json:"id" gorm:"column:id"`
	PONumber         string     `json:"po_number" gorm:"column:po_number"`
	ReferenceNo      *string    `json:"reference_no,omitempty" gorm:"column:reference_no"`
	SupplierID       uuid.UUID  `json:"supplier_id" gorm:"column:supplier_id"`
	SupplierName     string     `json:"supplier_name" gorm:"column:supplier_name"`
	SupplierCode     string     `json:"supplier_code" gorm:"column:supplier_code"`
	StoreID          uuid.UUID  `json:"store_id" gorm:"column:store_id"`
	StoreName        string     `json:"store_name" gorm:"column:store_name"`
	WarehouseID      uuid.UUID  `json:"warehouse_id" gorm:"column:warehouse_id"`
	WarehouseName    string     `json:"warehouse_name" gorm:"column:warehouse_name"`
	OrderDate        time.Time  `json:"order_date" gorm:"column:order_date"`
	ExpectedDelivery *time.Time `json:"expected_delivery,omitempty" gorm:"column:expected_delivery"`
	TotalAmount      decimal.Decimal `json:"total_amount" gorm:"column:total_amount"`
	Status           string     `json:"status" gorm:"column:status"`
	ApprovedByID     *uuid.UUID `json:"approved_by_id,omitempty" gorm:"column:approved_by_id"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty" gorm:"column:approved_at"`
	CreatedByID      uuid.UUID  `json:"created_by_id" gorm:"column:created_by_id"`
	CreatedByName    string     `json:"created_by_name" gorm:"column:created_by_name"`
	CreatedAt        time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"column:updated_at"`
}

type PurchaseOrderDetailRow struct {
	ID                     uuid.UUID  `json:"id" gorm:"column:id"`
	PONumber               string     `json:"po_number" gorm:"column:po_number"`
	ReferenceNo            *string    `json:"reference_no,omitempty" gorm:"column:reference_no"`
	SupplierID             uuid.UUID  `json:"supplier_id" gorm:"column:supplier_id"`
	SupplierNameSnapshot   string     `json:"supplier_name" gorm:"column:supplier_name"`
	SupplierCodeSnapshot   string     `json:"supplier_code" gorm:"column:supplier_code"`
	SupplierAddressSnapshot *string   `json:"supplier_address,omitempty" gorm:"column:supplier_address"`
	StoreID                uuid.UUID  `json:"store_id" gorm:"column:store_id"`
	StoreCodeSnapshot      string     `json:"store_code" gorm:"column:store_code"`
	StoreNameSnapshot      string     `json:"store_name" gorm:"column:store_name"`
	StoreAddressSnapshot   *string    `json:"store_address,omitempty" gorm:"column:store_address"`
	WarehouseID            uuid.UUID  `json:"warehouse_id" gorm:"column:warehouse_id"`
	WarehouseNameSnapshot  string     `json:"warehouse_name" gorm:"column:warehouse_name"`
	OrderDate              time.Time  `json:"order_date" gorm:"column:order_date"`
	ExpectedDelivery       *time.Time `json:"expected_delivery,omitempty" gorm:"column:expected_delivery"`
	PaymentTermDays        int        `json:"payment_term_days" gorm:"column:payment_term_days"`
	PaymentMode            string     `json:"payment_mode" gorm:"column:payment_mode"`
	TotalAmount            decimal.Decimal `json:"total_amount" gorm:"column:total_amount"`
	Status                 string     `json:"status" gorm:"column:status"`
	ApprovedByID           *uuid.UUID `json:"approved_by_id,omitempty" gorm:"column:approved_by_id"`
	ApprovedAt             *time.Time `json:"approved_at,omitempty" gorm:"column:approved_at"`
	ApprovedByName         *string    `json:"approved_by_name,omitempty" gorm:"column:approved_by_name"`
	CreatedByID            uuid.UUID  `json:"created_by_id" gorm:"column:created_by_id"`
	CreatedByName          string     `json:"created_by_name" gorm:"column:created_by_name"`
	Notes                  *string    `json:"notes,omitempty" gorm:"column:notes"`
	SupplierNotes          *string    `json:"supplier_notes,omitempty" gorm:"column:supplier_notes"`
	CreatedAt              time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt              time.Time  `json:"updated_at" gorm:"column:updated_at"`

	Items []PurchaseOrderItemRow `json:"items" gorm:"-"`
}

type PurchaseOrderItemRow struct {
	ID                uuid.UUID       `json:"id" gorm:"column:id"`
	PurchaseOrderID   uuid.UUID       `json:"purchase_order_id" gorm:"column:purchase_order_id"`
	SeqNo             int             `json:"seq_no" gorm:"column:seq_no"`
	ProductID         uuid.UUID       `json:"product_id" gorm:"column:product_id"`
	ProductSKU        string          `json:"product_sku" gorm:"column:product_sku"`
	ProductName       string          `json:"product_name" gorm:"column:product_name"`
	UOMID             uuid.UUID       `json:"uom_id" gorm:"column:uom_id"`
	UOMName           string          `json:"uom_name" gorm:"column:uom_name"`
	QtyOrdered        decimal.Decimal `json:"qty_ordered" gorm:"column:qty_ordered"`
	QtyReceived       decimal.Decimal `json:"qty_received" gorm:"column:qty_received"`
	UnitPrice         decimal.Decimal `json:"unit_price" gorm:"column:unit_price"`
	Subtotal          decimal.Decimal `json:"subtotal" gorm:"column:subtotal"`
	ProductSupplierID *uuid.UUID      `json:"product_supplier_id,omitempty" gorm:"column:product_supplier_id"`
	PlanningID        *uuid.UUID      `json:"planning_id,omitempty" gorm:"column:planning_id"`
	Notes             *string         `json:"notes,omitempty" gorm:"column:notes"`
}
