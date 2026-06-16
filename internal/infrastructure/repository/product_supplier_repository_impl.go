package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productSupplierRepository struct {
	db *gorm.DB
}

func NewProductSupplierRepository(db *gorm.DB) domainRepo.ProductSupplierRepository {
	return &productSupplierRepository{db: db}
}

func (r *productSupplierRepository) Create(ctx context.Context, ps *entity.ProductSupplier) error {
	return uow.GetTx(ctx, r.db).Create(ps).Error
}

func (r *productSupplierRepository) FindByID(ctx context.Context, id string) (*entity.ProductSupplier, error) {
	var ps entity.ProductSupplier
	err := r.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&ps).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product supplier: %w", err)
	}
	return &ps, nil
}

func (r *productSupplierRepository) FindByProductAndSupplier(ctx context.Context, productID, supplierID string) (*entity.ProductSupplier, error) {
	var ps entity.ProductSupplier
	err := r.db.WithContext(ctx).Preload("Supplier").Preload("Store").Where("product_id = ? AND supplier_id = ?", productID, supplierID).First(&ps).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product supplier: %w", err)
	}
	return &ps, nil
}

func (r *productSupplierRepository) FindByProductID(ctx context.Context, productID string) ([]entity.ProductSupplier, error) {
	var pss []entity.ProductSupplier
	err := r.db.WithContext(ctx).Preload("Supplier").Preload("Store").Where("product_id = ?", productID).Find(&pss).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product suppliers: %w", err)
	}
	return pss, nil
}

func (r *productSupplierRepository) FindBySupplierID(ctx context.Context, supplierID string) ([]entity.ProductSupplier, error) {
	var pss []entity.ProductSupplier
	err := r.db.WithContext(ctx).Preload("Product").Preload("Store").Where("supplier_id = ?", supplierID).Find(&pss).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product suppliers: %w", err)
	}
	return pss, nil
}

func (r *productSupplierRepository) FindAll(ctx context.Context) ([]entity.ProductSupplier, error) {
	var pss []entity.ProductSupplier
	err := r.db.WithContext(ctx).Preload("Product").Preload("Supplier").Preload("Store").Find(&pss).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product suppliers: %w", err)
	}
	return pss, nil
}

func (r *productSupplierRepository) FindByStoreID(ctx context.Context, storeID string) ([]entity.ProductSupplier, error) {
	var pss []entity.ProductSupplier
	err := r.db.WithContext(ctx).Preload("Product").Preload("Supplier").Preload("Store").Where("store_id = ?", storeID).Find(&pss).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product suppliers by store: %w", err)
	}
	return pss, nil
}

func (r *productSupplierRepository) Update(ctx context.Context, ps *entity.ProductSupplier) error {
	return uow.GetTx(ctx, r.db).Save(ps).Error
}

func (r *productSupplierRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.ProductSupplier{}).Error
}