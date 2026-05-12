package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type InventoryStockRepository interface {
	FindByWarehouseAndProduct(ctx context.Context, warehouseID, productID string) (*entity.InventoryStock, error)
	FindByWarehouseID(ctx context.Context, warehouseID string) ([]entity.InventoryStock, error)
	FindByProductID(ctx context.Context, productID string) ([]entity.InventoryStock, error)
	FindAll(ctx context.Context) ([]entity.InventoryStock, error)
	Update(ctx context.Context, is *entity.InventoryStock) error
}