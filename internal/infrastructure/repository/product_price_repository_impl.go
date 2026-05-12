package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type productPriceRepository struct {
	db *gorm.DB
}

func NewProductPriceRepository(db *gorm.DB) domainRepo.ProductPriceRepository {
	return &productPriceRepository{db: db}
}

func (r *productPriceRepository) Create(ctx context.Context, pp *entity.ProductPrice) error {
	return uow.GetTx(ctx, r.db).Create(pp).Error
}

func (r *productPriceRepository) FindByID(ctx context.Context, id string) (*entity.ProductPrice, error) {
	var pp entity.ProductPrice
	err := r.db.WithContext(ctx).Preload("PriceList").Preload("Product").Preload("UOM").Where("id = ?", id).First(&pp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product price: %w", err)
	}
	return &pp, nil
}

func (r *productPriceRepository) FindByPriceListAndProduct(ctx context.Context, priceListID, productID string) (*entity.ProductPrice, error) {
	var pp entity.ProductPrice
	err := r.db.WithContext(ctx).Where("price_list_id = ? AND product_id = ?", priceListID, productID).First(&pp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product price: %w", err)
	}
	return &pp, nil
}

func (r *productPriceRepository) FindByProductID(ctx context.Context, productID string) ([]entity.ProductPrice, error) {
	var pps []entity.ProductPrice
	err := r.db.WithContext(ctx).Preload("PriceList").Preload("UOM").Where("product_id = ?", productID).Find(&pps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product prices: %w", err)
	}
	return pps, nil
}

func (r *productPriceRepository) FindByPriceListID(ctx context.Context, priceListID string) ([]entity.ProductPrice, error) {
	var pps []entity.ProductPrice
	err := r.db.WithContext(ctx).Preload("Product").Preload("UOM").Where("price_list_id = ?", priceListID).Find(&pps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product prices: %w", err)
	}
	return pps, nil
}

func (r *productPriceRepository) FindAll(ctx context.Context) ([]entity.ProductPrice, error) {
	var pps []entity.ProductPrice
	err := r.db.WithContext(ctx).Preload("PriceList").Preload("Product").Preload("UOM").Find(&pps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find product prices: %w", err)
	}
	return pps, nil
}

func (r *productPriceRepository) Update(ctx context.Context, pp *entity.ProductPrice) error {
	return uow.GetTx(ctx, r.db).Save(pp).Error
}

func (r *productPriceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.ProductPrice{}).Error
}