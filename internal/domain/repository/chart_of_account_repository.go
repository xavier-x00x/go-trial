package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type ChartOfAccountRepository interface {
	Create(ctx context.Context, coa *entity.ChartOfAccount) error
	FindByID(ctx context.Context, id string) (*entity.ChartOfAccount, error)
	FindByCode(ctx context.Context, code string) (*entity.ChartOfAccount, error)
	FindAll(ctx context.Context) ([]entity.ChartOfAccount, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.ChartOfAccount, *entity.Meta, error)
	FindByType(ctx context.Context, accountType string) ([]entity.ChartOfAccount, error)
	Update(ctx context.Context, coa *entity.ChartOfAccount) error
	Delete(ctx context.Context, id string) error
	BulkCreate(ctx context.Context, coas []entity.ChartOfAccount) error
	FindByParentID(ctx context.Context, parentID *string) ([]entity.ChartOfAccount, error)
}