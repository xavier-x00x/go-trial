package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type PurchaseReturnRepository interface {
	Create(ctx context.Context, pr *entity.PurchaseReturn) error
	Update(ctx context.Context, pr *entity.PurchaseReturn) error
	FindByID(ctx context.Context, id string) (*entity.PurchaseReturn, error)
	FindByReturnNumber(ctx context.Context, returnNum string) (*entity.PurchaseReturn, error)
	FindAllWithPagination(ctx context.Context, filter *QueryFilter) ([]entity.PurchaseReturn, *entity.Meta, error)
}
