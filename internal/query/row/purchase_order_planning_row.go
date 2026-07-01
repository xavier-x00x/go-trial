package row

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PurchaseOrderPlanningRow struct {
	ID                    uuid.UUID       `json:"id" gorm:"column:id"`
	StoreID               uuid.UUID       `json:"store_id" gorm:"column:store_id"`
	ProductID             uuid.UUID       `json:"product_id" gorm:"column:product_id"`
	ProductSupplierID     uuid.UUID       `json:"product_supplier_id" gorm:"column:product_supplier_id"`
	ProductSKU            string          `json:"product_sku" gorm:"column:product_sku"`
	ProductName           string          `json:"product_name" gorm:"column:product_name"`
	SupplierCode          string          `json:"supplier_code" gorm:"column:supplier_code"`
	SupplierName          string          `json:"supplier_name" gorm:"column:supplier_name"`
	CurrentStock          decimal.Decimal `json:"current_stock" gorm:"column:current_stock"`
	SafetyStock           decimal.Decimal `json:"safety_stock" gorm:"column:safety_stock"`
	DynamicSafetyStock    decimal.Decimal `json:"dynamic_safety_stock" gorm:"column:dynamic_safety_stock"`
	ReorderPoint          decimal.Decimal `json:"reorder_point" gorm:"column:reorder_point"`
	AverageDailySales     decimal.Decimal `json:"average_daily_sales" gorm:"column:average_daily_sales"`
	LeadTimeDays          int             `json:"lead_time_days" gorm:"column:lead_time_days"`
	LeadTimeDemand        decimal.Decimal `json:"lead_time_demand" gorm:"column:lead_time_demand"`
	Status                string          `json:"status" gorm:"column:status"`
	RecommendedOrderQty   decimal.Decimal `json:"recommended_order_qty" gorm:"column:recommended_order_qty"`
	OrderQty              decimal.Decimal `json:"order_qty" gorm:"column:order_qty"`
	IsManualSupplier      bool            `json:"is_manual_supplier" gorm:"column:is_manual_supplier"`
	IsSelected            bool            `json:"is_selected" gorm:"column:is_selected"`
	CalculatedDate        time.Time       `json:"calculated_date" gorm:"column:calculated_date"`
	ProcessedDate         *time.Time      `json:"processed_date,omitempty" gorm:"column:processed_date"`
	ProcessedByID         *uuid.UUID      `json:"processed_by_id,omitempty" gorm:"column:processed_by_id"`
}
