package repository

import (
	"context"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) domainRepo.SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, s *entity.Supplier) error {
	return uow.GetTx(ctx, r.db).Create(s).Error
}

func (r *supplierRepository) FindByID(ctx context.Context, id string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&supplier).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find supplier: %w", err)
	}
	return &supplier, nil
}

func (r *supplierRepository) FindByCode(ctx context.Context, code string) (*entity.Supplier, error) {
	var supplier entity.Supplier
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&supplier).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find supplier by code: %w", err)
	}
	return &supplier, nil
}

func (r *supplierRepository) FindAll(ctx context.Context) ([]entity.Supplier, error) {
	var suppliers []entity.Supplier
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&suppliers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find suppliers: %w", err)
	}
	return suppliers, nil
}

func (r *supplierRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Supplier, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Supplier{})
	return PaginateAndFilter[entity.Supplier](r.db, baseQuery, filter)
}

func (r *supplierRepository) Update(ctx context.Context, s *entity.Supplier) error {
	return uow.GetTx(ctx, r.db).Save(s).Error
}

func (r *supplierRepository) Delete(ctx context.Context, id string) error {
	var supplier entity.Supplier
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&supplier).Error; err != nil {
		return err
	}
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	return r.db.WithContext(ctx).Model(&supplier).Updates(map[string]interface{}{
		"code":       supplier.Code + suffix,
		"deleted_at": time.Now(),
	}).Error
}