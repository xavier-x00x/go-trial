package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type StoreProductAssortmentRepository interface {
	Create(ctx context.Context, spa *entity.StoreProductAssortment) error
	FindByID(ctx context.Context, id string) (*entity.StoreProductAssortment, error)
	FindByStoreAndProduct(ctx context.Context, storeID, productID string) (*entity.StoreProductAssortment, error)
	FindByStoreID(ctx context.Context, storeID string) ([]entity.StoreProductAssortment, error)
	FindByProductID(ctx context.Context, productID string) ([]entity.StoreProductAssortment, error)
	FindAll(ctx context.Context) ([]entity.StoreProductAssortment, error)
	FindForPlanning(ctx context.Context, storeID string) ([]entity.PlanningData, error)
	Update(ctx context.Context, spa *entity.StoreProductAssortment) error
	Delete(ctx context.Context, id string) error
}