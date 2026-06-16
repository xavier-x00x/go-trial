package repository

import (
	"context"
	"fmt"
	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type supplierCategoryRepository struct {
	db *gorm.DB
}

func NewSupplierCategoryRepository(db *gorm.DB) domainRepo.SupplierCategoryRepository {
	return &supplierCategoryRepository{db: db}
}

func (r *supplierCategoryRepository) Create(ctx context.Context, cat *entity.SupplierCategory) error {
	return uow.GetTx(ctx, r.db).Create(cat).Error
}

func (r *supplierCategoryRepository) FindByID(ctx context.Context, id string) (*entity.SupplierCategory, error) {
	var cat entity.SupplierCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find supplier category: %w", err)
	}
	return &cat, nil
}

func (r *supplierCategoryRepository) FindByName(ctx context.Context, name string) (*entity.SupplierCategory, error) {
	var cat entity.SupplierCategory
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find supplier category by name: %w", err)
	}
	return &cat, nil
}

func (r *supplierCategoryRepository) Update(ctx context.Context, cat *entity.SupplierCategory) error {
	return uow.GetTx(ctx, r.db).Save(cat).Error
}

func (r *supplierCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.SupplierCategory{}).Error
}
