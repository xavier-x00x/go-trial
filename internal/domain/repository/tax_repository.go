package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type TaxRepository interface {
	Create(ctx context.Context, t *entity.Tax) error
	FindByID(ctx context.Context, id string) (*entity.Tax, error)
	FindAll(ctx context.Context) ([]entity.Tax, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Tax, *entity.Meta, error)
	Update(ctx context.Context, t *entity.Tax) error
	Delete(ctx context.Context, id string) error
}