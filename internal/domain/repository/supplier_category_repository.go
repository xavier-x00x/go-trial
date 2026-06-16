package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type SupplierCategoryRepository interface {
	Create(ctx context.Context, cat *entity.SupplierCategory) error
	FindByID(ctx context.Context, id string) (*entity.SupplierCategory, error)
	FindByName(ctx context.Context, name string) (*entity.SupplierCategory, error)
	Update(ctx context.Context, cat *entity.SupplierCategory) error
	Delete(ctx context.Context, id string) error
}
