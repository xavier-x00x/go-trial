package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type MonthlyInventoryStockRepository interface {
	FindByPeriodWarehouseProduct(ctx context.Context, periodMonth, warehouseID, productID string) (*entity.MonthlyInventoryStock, error)
	FindByPeriodMonth(ctx context.Context, periodMonth string) ([]entity.MonthlyInventoryStock, error)
	FindByWarehouseID(ctx context.Context, warehouseID string) ([]entity.MonthlyInventoryStock, error)
	FindByProductID(ctx context.Context, productID string) ([]entity.MonthlyInventoryStock, error)
	Create(ctx context.Context, mis *entity.MonthlyInventoryStock) error
	Update(ctx context.Context, mis *entity.MonthlyInventoryStock) error
}