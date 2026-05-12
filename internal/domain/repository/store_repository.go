package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type StoreRepository interface {
	Create(ctx context.Context, store *entity.Store) error
	FindByID(ctx context.Context, id string) (*entity.Store, error)
	FindByCode(ctx context.Context, code string) (*entity.Store, error)
	FindAll(ctx context.Context) ([]entity.Store, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Store, *entity.Meta, error)
	Update(ctx context.Context, store *entity.Store) error
	Delete(ctx context.Context, id string) error
}
