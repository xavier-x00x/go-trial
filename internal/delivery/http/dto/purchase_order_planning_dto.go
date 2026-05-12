package dto

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrderPlanningResponse struct {
	ID                    uuid.UUID `json:"id"`
	StoreID              uuid.UUID `json:"store_id"`
	ProductID            uuid.UUID `json:"product_id"`
	ProductSupplierID    uuid.UUID `json:"product_supplier_id"`
	ProductSKU           string    `json:"product_sku"`
	ProductName          string    `json:"product_name"`
	SupplierCode         string    `json:"supplier_code"`
	SupplierName         string    `json:"supplier_name"`
	CurrentStock         float64   `json:"current_stock"`
	SafetyStock          float64   `json:"safety_stock"`
	DynamicSafetyStock   float64   `json:"dynamic_safety_stock"`
	ReorderPoint        float64   `json:"reorder_point"`
	AverageDailySales   float64   `json:"average_daily_sales"`
	LeadTimeDays        int       `json:"lead_time_days"`
	LeadTimeDemand      float64   `json:"lead_time_demand"`
	Status              string    `json:"status"`
	RecommendedOrderQty float64   `json:"recommended_order_qty"`
	CalculatedDate      time.Time `json:"calculated_date"`
	ProcessedDate       *time.Time `json:"processed_date,omitempty"`
	ProcessedByID       *uuid.UUID `json:"processed_by_id,omitempty"`
}

type PlanningSummaryResponse struct {
	ProcessedCount   int       `json:"processed_count"`
	NeedsOrderCount int       `json:"needs_order_count"`
	CalculatedDate time.Time `json:"calculated_date"`
}

type ApprovePlanningRequest struct {
	ProductIDs      []uuid.UUID `json:"product_ids" validate:"required"`
	OrderQuantities []float64   `json:"order_quantities" validate:"required"`
	ProcessedByID  uuid.UUID  `json:"processed_by_id" validate:"required"`
}

type CalculatePlanningRequest struct {
	StoreID     uuid.UUID `json:"store_id" validate:"required"`
	Date       string    `json:"date"`
	ForceRecal bool      `json:"force_recal"`
}