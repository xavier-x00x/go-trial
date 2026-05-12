package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type ProductUOMConversionRepository interface {
	Create(ctx context.Context, puc *entity.ProductUOMConversion) error
	FindByID(ctx context.Context, id string) (*entity.ProductUOMConversion, error)
	FindByProductID(ctx context.Context, productID string) ([]entity.ProductUOMConversion, error)
	FindByBarcode(ctx context.Context, barcode string) (*entity.ProductUOMConversion, error)
	FindAll(ctx context.Context) ([]entity.ProductUOMConversion, error)
	Update(ctx context.Context, puc *entity.ProductUOMConversion) error
	Delete(ctx context.Context, id string) error
}