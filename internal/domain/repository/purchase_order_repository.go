package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type PurchaseOrderRepository interface {
	Create(ctx context.Context, po *entity.PurchaseOrder) error
	FindByID(ctx context.Context, id string) (*entity.PurchaseOrder, error)
	FindByIDWithSupplier(ctx context.Context, id string) (*entity.PurchaseOrder, error)
	FindByPONumber(ctx context.Context, poNum string) (*entity.PurchaseOrder, error)
	FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseOrder, error)
	FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseOrder, error)
	FindBySupplierID(ctx context.Context, supplierID string) ([]entity.PurchaseOrder, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.PurchaseOrder, *entity.Meta, error)
	Update(ctx context.Context, po *entity.PurchaseOrder) error
	Delete(ctx context.Context, id string) error
	DeleteItemsByPurchaseOrderID(ctx context.Context, poID string) error
}