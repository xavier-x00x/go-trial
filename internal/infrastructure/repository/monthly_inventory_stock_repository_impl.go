package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"

	"gorm.io/gorm"
)

type monthlyInventoryStockRepository struct {
	db *gorm.DB
}

func NewMonthlyInventoryStockRepository(db *gorm.DB) domainRepo.MonthlyInventoryStockRepository {
	return &monthlyInventoryStockRepository{db: db}
}

func (r *monthlyInventoryStockRepository) FindByPeriodWarehouseProduct(ctx context.Context, periodMonth, warehouseID, productID string) (*entity.MonthlyInventoryStock, error) {
	var mis entity.MonthlyInventoryStock
	err := r.db.WithContext(ctx).Preload("Warehouse").Preload("Product").Where("period_month = ? AND warehouse_id = ? AND product_id = ?", periodMonth, warehouseID, productID).First(&mis).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find monthly inventory stock: %w", err)
	}
	return &mis, nil
}

func (r *monthlyInventoryStockRepository) FindByPeriodMonth(ctx context.Context, periodMonth string) ([]entity.MonthlyInventoryStock, error) {
	var misList []entity.MonthlyInventoryStock
	err := r.db.WithContext(ctx).Preload("Warehouse").Preload("Product").Where("period_month = ?", periodMonth).Find(&misList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find monthly inventory stocks: %w", err)
	}
	return misList, nil
}

func (r *monthlyInventoryStockRepository) FindByWarehouseID(ctx context.Context, warehouseID string) ([]entity.MonthlyInventoryStock, error) {
	var misList []entity.MonthlyInventoryStock
	err := r.db.WithContext(ctx).Preload("Product").Where("warehouse_id = ?", warehouseID).Find(&misList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find monthly inventory stocks: %w", err)
	}
	return misList, nil
}

func (r *monthlyInventoryStockRepository) FindByProductID(ctx context.Context, productID string) ([]entity.MonthlyInventoryStock, error) {
	var misList []entity.MonthlyInventoryStock
	err := r.db.WithContext(ctx).Preload("Warehouse").Where("product_id = ?", productID).Find(&misList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find monthly inventory stocks: %w", err)
	}
	return misList, nil
}

func (r *monthlyInventoryStockRepository) Create(ctx context.Context, mis *entity.MonthlyInventoryStock) error {
	return r.db.Create(mis).Error
}

func (r *monthlyInventoryStockRepository) Update(ctx context.Context, mis *entity.MonthlyInventoryStock) error {
	return r.db.Save(mis).Error
}