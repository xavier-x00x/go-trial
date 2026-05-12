package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindBySKU(ctx context.Context, sku string) (*entity.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*entity.Product, error)
	FindAll(ctx context.Context) ([]entity.Product, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Product, *entity.Meta, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]entity.Product, error)
	Update(ctx context.Context, p *entity.Product) error
	Delete(ctx context.Context, id string) error
}