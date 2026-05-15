package repository

import (
	"context"

	"go-trial/internal/domain/entity"

	"github.com/shopspring/decimal"
)

type GoodsReceiptRepository interface {
	Create(ctx context.Context, gr *entity.GoodsReceipt) error
	FindByID(ctx context.Context, id string) (*entity.GoodsReceipt, error)
	FindByGRNumber(ctx context.Context, grNum string) (*entity.GoodsReceipt, error)
	FindByPurchaseOrderID(ctx context.Context, poID string) ([]entity.GoodsReceipt, error)
	FindByWarehouseID(ctx context.Context, warehouseID string, status string) ([]entity.GoodsReceipt, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.GoodsReceipt, *entity.Meta, error)
	Update(ctx context.Context, gr *entity.GoodsReceipt) error
	Delete(ctx context.Context, id string) error
	DeleteItemsByGoodsReceiptID(ctx context.Context, grID string) error
	GetTotalDraftQtyByPOItemID(ctx context.Context, poItemID string, excludeGRID *string) (decimal.Decimal, error)
}