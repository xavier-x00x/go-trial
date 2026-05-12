package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type storeProductAssortmentRepository struct {
	db *gorm.DB
}

func NewStoreProductAssortmentRepository(db *gorm.DB) domainRepo.StoreProductAssortmentRepository {
	return &storeProductAssortmentRepository{db: db}
}

func (r *storeProductAssortmentRepository) Create(ctx context.Context, spa *entity.StoreProductAssortment) error {
	return uow.GetTx(ctx, r.db).Create(spa).Error
}

func (r *storeProductAssortmentRepository) FindByID(ctx context.Context, id string) (*entity.StoreProductAssortment, error) {
	var spa entity.StoreProductAssortment
	err := r.db.WithContext(ctx).Preload("Store").Preload("Product").Where("id = ?", id).First(&spa).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find store product assortment: %w", err)
	}
	return &spa, nil
}

func (r *storeProductAssortmentRepository) FindByStoreAndProduct(ctx context.Context, storeID, productID string) (*entity.StoreProductAssortment, error) {
	var spa entity.StoreProductAssortment
	err := r.db.WithContext(ctx).Where("store_id = ? AND product_id = ?", storeID, productID).First(&spa).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find store product assortment: %w", err)
	}
	return &spa, nil
}

func (r *storeProductAssortmentRepository) FindByStoreID(ctx context.Context, storeID string) ([]entity.StoreProductAssortment, error) {
	var spas []entity.StoreProductAssortment
	err := r.db.WithContext(ctx).Preload("Product").Where("store_id = ?", storeID).Find(&spas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find store product assortments: %w", err)
	}
	return spas, nil
}

func (r *storeProductAssortmentRepository) FindByProductID(ctx context.Context, productID string) ([]entity.StoreProductAssortment, error) {
	var spas []entity.StoreProductAssortment
	err := r.db.WithContext(ctx).Preload("Store").Where("product_id = ?", productID).Find(&spas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find store product assortments: %w", err)
	}
	return spas, nil
}

func (r *storeProductAssortmentRepository) FindAll(ctx context.Context) ([]entity.StoreProductAssortment, error) {
	var spas []entity.StoreProductAssortment
	err := r.db.WithContext(ctx).Preload("Store").Preload("Product").Find(&spas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find store product assortments: %w", err)
	}
	return spas, nil
}

func (r *storeProductAssortmentRepository) FindForPlanning(ctx context.Context, storeID string) ([]entity.PlanningData, error) {
	var results []entity.PlanningData
	err := r.db.WithContext(ctx).
		Table("store_product_assortments spa").
		Select(`
			spa.store_id, 
			spa.product_id, 
			ps.id as product_supplier_id, 
			ps.default_lead_time_days, 
			spa.average_daily_sales, 
			spa.safety_stock_qty, 
			spa.max_stock_qty,
			COALESCE(SUM(is.quantity - is.reserved_qty), 0) as current_stock
		`).
		Joins("JOIN product_suppliers ps ON ps.product_id = spa.product_id").
		Joins("JOIN warehouses w ON w.store_id = spa.store_id AND w.is_active = true").
		Joins("LEFT JOIN inventory_stocks is ON is.warehouse_id = w.id AND is.product_id = spa.product_id").
		Where("spa.store_id = ?", storeID).
		Group("spa.store_id, spa.product_id, ps.id, ps.default_lead_time_days, spa.average_daily_sales, spa.safety_stock_qty, spa.max_stock_qty").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find planning data: %w", err)
	}
	return results, nil
}

func (r *storeProductAssortmentRepository) Update(ctx context.Context, spa *entity.StoreProductAssortment) error {
	return uow.GetTx(ctx, r.db).Save(spa).Error
}

func (r *storeProductAssortmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.StoreProductAssortment{}).Error
}