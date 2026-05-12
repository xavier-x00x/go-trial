package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *entity.Customer) error
	FindByID(ctx context.Context, id string) (*entity.Customer, error)
	FindByCode(ctx context.Context, code string) (*entity.Customer, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Customer, error)
	FindAll(ctx context.Context) ([]entity.Customer, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Customer, *entity.Meta, error)
	Update(ctx context.Context, c *entity.Customer) error
	Delete(ctx context.Context, id string) error
}