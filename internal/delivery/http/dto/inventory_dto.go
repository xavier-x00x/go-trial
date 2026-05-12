package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InventoryStockResponse struct {
	ID              uuid.UUID        `json:"id"`
	WarehouseID     uuid.UUID        `json:"warehouse_id"`
	Warehouse       *WarehouseResponse `json:"warehouse,omitempty"`
	ProductID       uuid.UUID        `json:"product_id"`
	Product         *ProductResponse   `json:"product,omitempty"`
	Quantity       decimal.Decimal `json:"quantity"`
	ReservedQty    decimal.Decimal `json:"reserved_qty"`
	AverageBuyPrice decimal.Decimal `json:"average_buy_price"`
	LastAuditAt     *time.Time     `json:"last_audit_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type MonthlyInventoryStockResponse struct {
	ID               uuid.UUID        `json:"id"`
	PeriodMonth     string          `json:"period_month"`
	WarehouseID     uuid.UUID        `json:"warehouse_id"`
	Warehouse       *WarehouseResponse `json:"warehouse,omitempty"`
	ProductID       uuid.UUID        `json:"product_id"`
	Product         *ProductResponse   `json:"product,omitempty"`
	BeginningBalance decimal.Decimal `json:"beginning_balance"`
	TotalIn         decimal.Decimal `json:"total_in"`
	TotalOut        decimal.Decimal `json:"total_out"`
	EndingBalance   decimal.Decimal `json:"ending_balance"`
	EndingValue     decimal.Decimal `json:"ending_value"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type MonthlyAPBalanceResponse struct {
	ID                uuid.UUID        `json:"id"`
	PeriodMonth      string          `json:"period_month"`
	SupplierID       uuid.UUID        `json:"supplier_id"`
	Supplier         *SupplierResponse   `json:"supplier,omitempty"`
	BeginningBalance decimal.Decimal `json:"beginning_balance"`
	TotalDebit       decimal.Decimal `json:"total_debit"`
	TotalCredit      decimal.Decimal `json:"total_credit"`
	EndingBalance    decimal.Decimal `json:"ending_balance"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type StoreProductAssortmentResponse struct {
	ID                   uuid.UUID        `json:"id"`
	StoreID             uuid.UUID        `json:"store_id"`
	Store               *StoreResponse   `json:"store,omitempty"`
	ProductID           uuid.UUID        `json:"product_id"`
	Product             *ProductResponse `json:"product,omitempty"`
	Status              string          `json:"status"`
	DisplayFacingQty    int             `json:"display_facing_qty"`
	DisplayShelfCapacity decimal.Decimal `json:"display_shelf_capacity"`
	VelocityClass      string          `json:"velocity_class"`
	VelocityLookbackDays int             `json:"velocity_lookback_days"`
	SafetyStockQty     decimal.Decimal `json:"safety_stock_qty"`
	ReorderPointQty    decimal.Decimal `json:"reorder_point_qty"`
	MaxStockQty        decimal.Decimal `json:"max_stock_qty"`
	AverageDailySales decimal.Decimal `json:"average_daily_sales"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}