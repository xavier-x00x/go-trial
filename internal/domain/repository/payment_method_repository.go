package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type PaymentMethodRepository interface {
	Create(ctx context.Context, pm *entity.PaymentMethod) error
	FindByID(ctx context.Context, id string) (*entity.PaymentMethod, error)
	FindByCode(ctx context.Context, code string) (*entity.PaymentMethod, error)
	FindAll(ctx context.Context) ([]entity.PaymentMethod, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.PaymentMethod, *entity.Meta, error)
	Update(ctx context.Context, pm *entity.PaymentMethod) error
	Delete(ctx context.Context, id string) error
}