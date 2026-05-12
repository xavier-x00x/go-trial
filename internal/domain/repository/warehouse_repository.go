package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type WarehouseRepository interface {
	Create(ctx context.Context, w *entity.Warehouse) error
	FindByID(ctx context.Context, id string) (*entity.Warehouse, error)
	FindByCode(ctx context.Context, code string) (*entity.Warehouse, error)
	FindByStoreID(ctx context.Context, storeID string) ([]entity.Warehouse, error)
	FindAll(ctx context.Context) ([]entity.Warehouse, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Warehouse, *entity.Meta, error)
	Update(ctx context.Context, w *entity.Warehouse) error
	Delete(ctx context.Context, id string) error
}