package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type purchaseOrderRepository struct {
	db *gorm.DB
}

func NewPurchaseOrderRepository(db *gorm.DB) domainRepo.PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

func (r *purchaseOrderRepository) Create(ctx context.Context, po *entity.PurchaseOrder) error {
	return uow.GetTx(ctx, r.db).Create(po).Error
}

func (r *purchaseOrderRepository) FindByID(ctx context.Context, id string) (*entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Preload("Items.Product").
		Preload("Items.UOM").
		Preload("Items.ProductSupplier").
		Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase order: %w", err)
	}
	return &po, nil
}

func (r *purchaseOrderRepository) FindByIDWithSupplier(ctx context.Context, id string) (*entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Where("id = ?", id).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase order with supplier: %w", err)
	}
	return &po, nil
}

func (r *purchaseOrderRepository) FindByPONumber(ctx context.Context, poNum string) (*entity.PurchaseOrder, error) {
	var po entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Preload("Items").
		Where("po_number = ?", poNum).First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase order by number: %w", err)
	}
	return &po, nil
}

func (r *purchaseOrderRepository) FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseOrder, error) {
	var pos []entity.PurchaseOrder
	query := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Where("store_id = ?", storeID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find purchase orders: %w", err)
	}
	return pos, nil
}

func (r *purchaseOrderRepository) FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseOrder, error) {
	var pos []entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Where("store_id = ? AND status IN ?", storeID, []string{"DRAFT", "SUBMITTED"}).
		Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find pending purchase orders: %w", err)
	}
	return pos, nil
}

func (r *purchaseOrderRepository) FindBySupplierID(ctx context.Context, supplierID string) ([]entity.PurchaseOrder, error) {
	var pos []entity.PurchaseOrder
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Where("supplier_id = ?", supplierID).
		Order("created_at DESC").Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find purchase orders by supplier: %w", err)
	}
	return pos, nil
}

func (r *purchaseOrderRepository) Update(ctx context.Context, po *entity.PurchaseOrder) error {
	return uow.GetTx(ctx, r.db).Save(po).Error
}

func (r *purchaseOrderRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.PurchaseOrder, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.PurchaseOrder{}).Preload("Supplier").Preload("Store").Preload("Warehouse").Preload("CreatedBy").Preload("ApprovedBy")
	return PaginateAndFilter[entity.PurchaseOrder](r.db, baseQuery, filter)
}

func (r *purchaseOrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.PurchaseOrder{}).Error
}

func (r *purchaseOrderRepository) DeleteItemsByPurchaseOrderID(ctx context.Context, poID string) error {
	return uow.GetTx(ctx, r.db).Where("purchase_order_id = ?", poID).Delete(&entity.PurchaseOrderItem{}).Error
}