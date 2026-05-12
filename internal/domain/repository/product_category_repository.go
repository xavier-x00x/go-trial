package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, cat *entity.ProductCategory) error
	FindByID(ctx context.Context, id string) (*entity.ProductCategory, error)
	FindBySlug(ctx context.Context, slug string) (*entity.ProductCategory, error)
	FindAll(ctx context.Context) ([]entity.ProductCategory, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.ProductCategory, *entity.Meta, error)
	FindByParentID(ctx context.Context, parentID string) ([]entity.ProductCategory, error)
	Update(ctx context.Context, cat *entity.ProductCategory) error
	Delete(ctx context.Context, id string) error
}