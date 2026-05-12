package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type PriceListRepository interface {
	Create(ctx context.Context, pl *entity.PriceList) error
	FindByID(ctx context.Context, id string) (*entity.PriceList, error)
	FindByCode(ctx context.Context, code string) (*entity.PriceList, error)
	FindAll(ctx context.Context) ([]entity.PriceList, error)
	FindActive(ctx context.Context) ([]entity.PriceList, error)
	Update(ctx context.Context, pl *entity.PriceList) error
	Delete(ctx context.Context, id string) error
}