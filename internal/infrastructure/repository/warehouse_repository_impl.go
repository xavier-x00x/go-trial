package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) domainRepo.WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(ctx context.Context, w *entity.Warehouse) error {
	return uow.GetTx(ctx, r.db).Create(w).Error
}

func (r *warehouseRepository) FindByID(ctx context.Context, id string) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.WithContext(ctx).Preload("Store").Where("id = ?", id).First(&warehouse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find warehouse: %w", err)
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindByCode(ctx context.Context, code string) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&warehouse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find warehouse by code: %w", err)
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindByStoreID(ctx context.Context, storeID string) ([]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	err := r.db.WithContext(ctx).Where("store_id = ? AND deleted_at IS NULL", storeID).Find(&warehouses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find warehouses: %w", err)
	}
	return warehouses, nil
}

func (r *warehouseRepository) FindAll(ctx context.Context) ([]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&warehouses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find warehouses: %w", err)
	}
	return warehouses, nil
}

func (r *warehouseRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Warehouse, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Warehouse{})
	return PaginateAndFilter[entity.Warehouse](r.db, baseQuery, filter)
}

func (r *warehouseRepository) Update(ctx context.Context, w *entity.Warehouse) error {
	return uow.GetTx(ctx, r.db).Save(w).Error
}

func (r *warehouseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.Warehouse{}).Error
}