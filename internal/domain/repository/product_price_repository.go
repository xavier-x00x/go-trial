package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type ProductPriceRepository interface {
	Create(ctx context.Context, pp *entity.ProductPrice) error
	FindByID(ctx context.Context, id string) (*entity.ProductPrice, error)
	FindByPriceListAndProduct(ctx context.Context, priceListID, productID string) (*entity.ProductPrice, error)
	FindByPriceListProductAndUOM(ctx context.Context, priceListID, productID, uomID string) (*entity.ProductPrice, error)
	FindByProductID(ctx context.Context, productID string) ([]entity.ProductPrice, error)
	FindByPriceListID(ctx context.Context, priceListID string) ([]entity.ProductPrice, error)
	FindAll(ctx context.Context) ([]entity.ProductPrice, error)
	Update(ctx context.Context, pp *entity.ProductPrice) error
	Delete(ctx context.Context, id string) error
	ExistsByProductIDAndUOMID(ctx context.Context, productID, uomID string) (bool, error)
}