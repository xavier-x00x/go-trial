package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"

	"gorm.io/gorm"
)

type inventoryStockRepository struct {
	db *gorm.DB
}

func NewInventoryStockRepository(db *gorm.DB) domainRepo.InventoryStockRepository {
	return &inventoryStockRepository{db: db}
}

func (r *inventoryStockRepository) FindByWarehouseAndProduct(ctx context.Context, warehouseID, productID string) (*entity.InventoryStock, error) {
	var is entity.InventoryStock
	err := r.db.WithContext(ctx).Preload("Warehouse").Preload("Product").Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).First(&is).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find inventory stock: %w", err)
	}
	return &is, nil
}

func (r *inventoryStockRepository) FindByWarehouseID(ctx context.Context, warehouseID string) ([]entity.InventoryStock, error) {
	var stocks []entity.InventoryStock
	err := r.db.WithContext(ctx).Preload("Product").Where("warehouse_id = ?", warehouseID).Find(&stocks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory stocks: %w", err)
	}
	return stocks, nil
}

func (r *inventoryStockRepository) FindByProductID(ctx context.Context, productID string) ([]entity.InventoryStock, error) {
	var stocks []entity.InventoryStock
	err := r.db.WithContext(ctx).Preload("Warehouse").Where("product_id = ?", productID).Find(&stocks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory stocks: %w", err)
	}
	return stocks, nil
}

func (r *inventoryStockRepository) FindAll(ctx context.Context) ([]entity.InventoryStock, error) {
	var stocks []entity.InventoryStock
	err := r.db.WithContext(ctx).Preload("Warehouse").Preload("Product").Find(&stocks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory stocks: %w", err)
	}
	return stocks, nil
}

func (r *inventoryStockRepository) Update(ctx context.Context, is *entity.InventoryStock) error {
	return r.db.Save(is).Error
}