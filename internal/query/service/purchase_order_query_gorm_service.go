package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"go-trial/internal/query/row"
)

type PurchaseOrderQueryService struct {
	db *gorm.DB
}

func NewPurchaseOrderQueryService(db *gorm.DB) *PurchaseOrderQueryService {
	return &PurchaseOrderQueryService{db: db}
}

func (s *PurchaseOrderQueryService) getBaseSelectAndJoins(ctx context.Context) *gorm.DB {
	return s.db.WithContext(ctx).
		Table("purchase_orders p").
		Select(`
			p.id,
			p.po_number,
			p.reference_no,
			p.supplier_id,
			p.supplier_name,
			p.supplier_code,
			p.store_id,
			p.store_name,
			p.warehouse_id,
			p.warehouse_name,
			p.order_date,
			p.expected_delivery,
			p.total_amount,
			p.status,
			p.approved_by_id,
			p.approved_at,
			p.created_by_id,
			p.created_by_name,
			p.created_at,
			p.updated_at
		`)
}

func (s *PurchaseOrderQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.PurchaseOrderListRow, *entity.Meta, error) {

	allowedOrder := []string{"id", "po_number", "order_date", "total_amount", "status", "created_at", "updated_at"}
	searchColumns := []string{"p.po_number", "p.reference_no", "p.supplier_name", "p.store_name"}

	// Map order column to use table alias to prevent ambiguous column error
	switch param.OrderColumn {
	case "id", "po_number", "order_date", "total_amount", "status", "created_at", "updated_at":
		param.OrderColumn = "p." + param.OrderColumn
	}
	allowedOrder = append(allowedOrder, "p.id", "p.po_number", "p.order_date", "p.total_amount", "p.status", "p.created_at", "p.updated_at")

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["p.deleted_at"] = nil

	baseQuery := s.getBaseSelectAndJoins(ctx)

	return PaginateAndFilter[row.PurchaseOrderListRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}

func (s *PurchaseOrderQueryService) GetByID(ctx context.Context, id string) (*row.PurchaseOrderDetailRow, error) {
	var detail row.PurchaseOrderDetailRow

	err := s.db.WithContext(ctx).
		Table("purchase_orders p").
		Select(`
			p.id,
			p.po_number,
			p.reference_no,
			p.supplier_id,
			p.supplier_name as supplier_name_snapshot,
			p.supplier_code as supplier_code_snapshot,
			p.supplier_address as supplier_address_snapshot,
			p.store_id,
			p.store_code as store_code_snapshot,
			p.store_name as store_name_snapshot,
			p.store_address as store_address_snapshot,
			p.warehouse_id,
			p.warehouse_name as warehouse_name_snapshot,
			p.order_date,
			p.expected_delivery,
			p.payment_term_days,
			p.payment_mode,
			p.total_amount,
			p.status,
			p.approved_by_id,
			p.approved_at,
			p.approved_by_name,
			p.created_by_id,
			p.created_by_name,
			p.notes,
			p.supplier_notes,
			p.created_at,
			p.updated_at
		`).
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&detail).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var items []row.PurchaseOrderItemRow
	err = s.db.WithContext(ctx).
		Table("purchase_order_items i").
		Select(`
			i.id,
			i.purchase_order_id,
			i.seq_no,
			i.product_id,
			i.product_sku,
			i.product_name,
			i.uom_id,
			i.uom_name,
			i.qty_ordered,
			i.qty_received,
			i.unit_price,
			i.subtotal,
			i.product_supplier_id,
			i.planning_id,
			i.notes
		`).
		Where("i.purchase_order_id = ? AND i.deleted_at IS NULL", id).
		Order("i.seq_no ASC").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	detail.Items = items
	return &detail, nil
}

func (s *PurchaseOrderQueryService) GetByPONumber(ctx context.Context, poNum string) (*row.PurchaseOrderDetailRow, error) {
	var detail row.PurchaseOrderDetailRow

	err := s.db.WithContext(ctx).
		Table("purchase_orders p").
		Select(`
			p.id,
			p.po_number,
			p.reference_no,
			p.supplier_id,
			p.supplier_name as supplier_name_snapshot,
			p.supplier_code as supplier_code_snapshot,
			p.supplier_address as supplier_address_snapshot,
			p.store_id,
			p.store_code as store_code_snapshot,
			p.store_name as store_name_snapshot,
			p.store_address as store_address_snapshot,
			p.warehouse_id,
			p.warehouse_name as warehouse_name_snapshot,
			p.order_date,
			p.expected_delivery,
			p.payment_term_days,
			p.payment_mode,
			p.total_amount,
			p.status,
			p.approved_by_id,
			p.approved_at,
			p.approved_by_name,
			p.created_by_id,
			p.created_by_name,
			p.notes,
			p.supplier_notes,
			p.created_at,
			p.updated_at
		`).
		Where("p.po_number = ? AND p.deleted_at IS NULL", poNum).
		First(&detail).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var items []row.PurchaseOrderItemRow
	err = s.db.WithContext(ctx).
		Table("purchase_order_items i").
		Select(`
			i.id,
			i.purchase_order_id,
			i.seq_no,
			i.product_id,
			i.product_sku,
			i.product_name,
			i.uom_id,
			i.uom_name,
			i.qty_ordered,
			i.qty_received,
			i.unit_price,
			i.subtotal,
			i.product_supplier_id,
			i.planning_id,
			i.notes
		`).
		Where("i.purchase_order_id = ? AND i.deleted_at IS NULL", detail.ID).
		Order("i.seq_no ASC").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	detail.Items = items
	return &detail, nil
}

func (s *PurchaseOrderQueryService) GetByStoreID(ctx context.Context, storeID string, status string) ([]row.PurchaseOrderListRow, error) {
	var rows []row.PurchaseOrderListRow
	query := s.getBaseSelectAndJoins(ctx).
		Where("p.deleted_at IS NULL").
		Where("p.store_id = ?", storeID)

	if status != "" {
		query = query.Where("p.status = ?", status)
	}

	err := query.Order("p.created_at DESC").Find(&rows).Error
	return rows, err
}

func (s *PurchaseOrderQueryService) GetPendingByStoreID(ctx context.Context, storeID string) ([]row.PurchaseOrderListRow, error) {
	var rows []row.PurchaseOrderListRow
	err := s.getBaseSelectAndJoins(ctx).
		Where("p.deleted_at IS NULL").
		Where("p.store_id = ?", storeID).
		Where("p.status IN ?", []string{entity.POStatusDraft, entity.POStatusSubmitted, entity.POStatusApproved}). // adjust pending status logically if needed
		Order("p.created_at ASC").
		Find(&rows).Error
	return rows, err
}
