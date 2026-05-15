package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type ExpenseVoucherRepository interface {
	Create(ctx context.Context, ev *entity.ExpenseVoucher) error
	Update(ctx context.Context, ev *entity.ExpenseVoucher) error
	FindByID(ctx context.Context, id string) (*entity.ExpenseVoucher, error)
	FindAllWithPagination(ctx context.Context, filter *QueryFilter) ([]entity.ExpenseVoucher, *entity.Meta, error)
	DeleteItemsByVoucherID(ctx context.Context, voucherID string) error
}
