package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/row"
)

type PurchaseOrderPlanningQueryService struct {
	db *gorm.DB
}

func NewPurchaseOrderPlanningQueryService(db *gorm.DB) *PurchaseOrderPlanningQueryService {
	return &PurchaseOrderPlanningQueryService{db: db}
}

func (s *PurchaseOrderPlanningQueryService) getBaseSelectAndJoins(ctx context.Context) *gorm.DB {
	return s.db.WithContext(ctx).
		Table("purchase_order_plannings p").
		Select(`
			p.id,
			p.store_id,
			p.product_id,
			p.product_supplier_id,
			pr.sku as product_sku,
			pr.name as product_name,
			sp.code as supplier_code,
			sp.name as supplier_name,
			p.current_stock,
			p.safety_stock,
			p.dynamic_safety_stock,
			p.reorder_point,
			p.average_daily_sales,
			p.lead_time_days,
			p.lead_time_demand,
			p.status,
			p.recommended_order_qty,
			p.calculated_date,
			p.processed_date,
			p.processed_by_id
		`).
		Joins("LEFT JOIN products pr ON pr.id = p.product_id").
		Joins("LEFT JOIN product_suppliers ps ON ps.id = p.product_supplier_id").
		Joins("LEFT JOIN suppliers sp ON sp.id = ps.supplier_id")
}

func (s *PurchaseOrderPlanningQueryService) GetPending(ctx context.Context, storeID string) ([]row.PurchaseOrderPlanningRow, error) {
	var rows []row.PurchaseOrderPlanningRow
	err := s.getBaseSelectAndJoins(ctx).
		Where("p.deleted_at IS NULL AND p.store_id = ? AND p.status = ?", storeID, entity.PlanningStatusPending).
		Order("p.calculated_date DESC").
		Find(&rows).Error
	return rows, err
}

func (s *PurchaseOrderPlanningQueryService) GetAll(ctx context.Context, storeID string, status string) ([]row.PurchaseOrderPlanningRow, error) {
	var rows []row.PurchaseOrderPlanningRow
	query := s.getBaseSelectAndJoins(ctx).
		Where("p.deleted_at IS NULL AND p.store_id = ?", storeID)
	
	if status != "" {
		query = query.Where("p.status = ?", status)
	}

	err := query.Order("p.calculated_date DESC").Find(&rows).Error
	return rows, err
}

func (s *PurchaseOrderPlanningQueryService) GetByID(ctx context.Context, id string) (*row.PurchaseOrderPlanningRow, error) {
	var detail row.PurchaseOrderPlanningRow

	err := s.getBaseSelectAndJoins(ctx).
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&detail).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}
