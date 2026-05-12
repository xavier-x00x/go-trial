package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type PurchaseOrderPlanningRepository interface {
	Create(ctx context.Context, p *entity.PurchaseOrderPlanning) error
	CreateBatch(ctx context.Context, plannings []entity.PurchaseOrderPlanning) error
	FindByID(ctx context.Context, id string) (*entity.PurchaseOrderPlanning, error)
	FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseOrderPlanning, error)
	FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseOrderPlanning, error)
	FindAll(ctx context.Context) ([]entity.PurchaseOrderPlanning, error)
	Update(ctx context.Context, p *entity.PurchaseOrderPlanning) error
	Delete(ctx context.Context, id string) error
	DeleteByDate(ctx context.Context, storeID string, date string) error
}