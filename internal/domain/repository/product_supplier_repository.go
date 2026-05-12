package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type ProductSupplierRepository interface {
	Create(ctx context.Context, ps *entity.ProductSupplier) error
	FindByID(ctx context.Context, id string) (*entity.ProductSupplier, error)
	FindByProductAndSupplier(ctx context.Context, productID, supplierID string) (*entity.ProductSupplier, error)
	FindByProductID(ctx context.Context, productID string) ([]entity.ProductSupplier, error)
	FindBySupplierID(ctx context.Context, supplierID string) ([]entity.ProductSupplier, error)
	FindByStoreID(ctx context.Context, storeID string) ([]entity.ProductSupplier, error)
	FindAll(ctx context.Context) ([]entity.ProductSupplier, error)
	Update(ctx context.Context, ps *entity.ProductSupplier) error
	Delete(ctx context.Context, id string) error
}