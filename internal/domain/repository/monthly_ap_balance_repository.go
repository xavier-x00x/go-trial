package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type MonthlyAPBalanceRepository interface {
	FindByPeriodSupplier(ctx context.Context, periodMonth, supplierID string) (*entity.MonthlyAPBalance, error)
	FindByPeriodMonth(ctx context.Context, periodMonth string) ([]entity.MonthlyAPBalance, error)
	FindBySupplierID(ctx context.Context, supplierID string) ([]entity.MonthlyAPBalance, error)
	Create(ctx context.Context, mab *entity.MonthlyAPBalance) error
	Update(ctx context.Context, mab *entity.MonthlyAPBalance) error
}