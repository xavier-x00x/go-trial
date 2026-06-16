package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
)

type ProductPriceQueryService struct {
	db *gorm.DB
}

func NewProductPriceQueryService(db *gorm.DB) *ProductPriceQueryService {
	return &ProductPriceQueryService{db: db}
}

func (s *ProductPriceQueryService) GetByID(ctx context.Context, id string) (*entity.ProductPrice, error) {
	var pp entity.ProductPrice
	err := s.db.WithContext(ctx).Preload("PriceList").Preload("Product").Preload("UOM").Where("id = ? AND deleted_at IS NULL", id).First(&pp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pp, nil
}

func (s *ProductPriceQueryService) GetAll(ctx context.Context) ([]entity.ProductPrice, error) {
	var pps []entity.ProductPrice
	err := s.db.WithContext(ctx).Preload("PriceList").Preload("Product").Preload("UOM").Where("deleted_at IS NULL").Find(&pps).Error
	if err != nil {
		return nil, err
	}
	return pps, nil
}

func (s *ProductPriceQueryService) GetByProductID(ctx context.Context, productID string) ([]entity.ProductPrice, error) {
	var pps []entity.ProductPrice
	err := s.db.WithContext(ctx).Preload("PriceList").Preload("UOM").Where("product_id = ? AND deleted_at IS NULL", productID).Find(&pps).Error
	if err != nil {
		return nil, err
	}
	return pps, nil
}

func (s *ProductPriceQueryService) GetByPriceListID(ctx context.Context, priceListID string) ([]entity.ProductPrice, error) {
	var pps []entity.ProductPrice
	err := s.db.WithContext(ctx).Preload("Product").Preload("UOM").Where("price_list_id = ? AND deleted_at IS NULL", priceListID).Find(&pps).Error
	if err != nil {
		return nil, err
	}
	return pps, nil
}
