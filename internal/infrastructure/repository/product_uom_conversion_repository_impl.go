package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type productUOMConversionRepository struct {
	db *gorm.DB
}

func NewProductUOMConversionRepository(db *gorm.DB) domainRepo.ProductUOMConversionRepository {
	return &productUOMConversionRepository{db: db}
}

func (r *productUOMConversionRepository) Create(ctx context.Context, puc *entity.ProductUOMConversion) error {
	return uow.GetTx(ctx, r.db).Create(puc).Error
}

func (r *productUOMConversionRepository) FindByID(ctx context.Context, id string) (*entity.ProductUOMConversion, error) {
	var puc entity.ProductUOMConversion
	err := r.db.WithContext(ctx).Preload("Product").Preload("UOM").Where("id = ?", id).First(&puc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product UOM conversion: %w", err)
	}
	return &puc, nil
}

func (r *productUOMConversionRepository) FindByProductID(ctx context.Context, productID string) ([]entity.ProductUOMConversion, error) {
	var pucs []entity.ProductUOMConversion
	err := r.db.WithContext(ctx).Preload("UOM").Where("product_id = ?", productID).Find(&pucs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product UOM conversions: %w", err)
	}
	return pucs, nil
}

func (r *productUOMConversionRepository) FindByBarcode(ctx context.Context, barcode string) (*entity.ProductUOMConversion, error) {
	var puc entity.ProductUOMConversion
	err := r.db.WithContext(ctx).Preload("Product").Preload("UOM").Where("barcode = ?", barcode).First(&puc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product UOM conversion by barcode: %w", err)
	}
	return &puc, nil
}

func (r *productUOMConversionRepository) FindAll(ctx context.Context) ([]entity.ProductUOMConversion, error) {
	var pucs []entity.ProductUOMConversion
	err := r.db.WithContext(ctx).Preload("Product").Preload("UOM").Find(&pucs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product UOM conversions: %w", err)
	}
	return pucs, nil
}

func (r *productUOMConversionRepository) Update(ctx context.Context, puc *entity.ProductUOMConversion) error {
	return uow.GetTx(ctx, r.db).Save(puc).Error
}

func (r *productUOMConversionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.ProductUOMConversion{}).Error
}