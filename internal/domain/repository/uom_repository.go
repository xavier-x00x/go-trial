package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type UOMRepository interface {
	Create(ctx context.Context, uom *entity.UOM) error
	FindByID(ctx context.Context, id string) (*entity.UOM, error)
	FindByCode(ctx context.Context, code string) (*entity.UOM, error)
	FindAll(ctx context.Context) ([]entity.UOM, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.UOM, *entity.Meta, error)
	Update(ctx context.Context, uom *entity.UOM) error
	Delete(ctx context.Context, id string) error
}