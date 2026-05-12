package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type SupplierRepository interface {
	Create(ctx context.Context, s *entity.Supplier) error
	FindByID(ctx context.Context, id string) (*entity.Supplier, error)
	FindByCode(ctx context.Context, code string) (*entity.Supplier, error)
	FindAll(ctx context.Context) ([]entity.Supplier, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Supplier, *entity.Meta, error)
	Update(ctx context.Context, s *entity.Supplier) error
	Delete(ctx context.Context, id string) error
}