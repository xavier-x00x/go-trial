package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type priceListRepository struct {
	db *gorm.DB
}

func NewPriceListRepository(db *gorm.DB) domainRepo.PriceListRepository {
	return &priceListRepository{db: db}
}

func (r *priceListRepository) Create(ctx context.Context, pl *entity.PriceList) error {
	return uow.GetTx(ctx, r.db).Create(pl).Error
}

func (r *priceListRepository) FindByID(ctx context.Context, id string) (*entity.PriceList, error) {
	var pl entity.PriceList
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&pl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find price list: %w", err)
	}
	return &pl, nil
}

func (r *priceListRepository) FindByCode(ctx context.Context, code string) (*entity.PriceList, error) {
	var pl entity.PriceList
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&pl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find price list by code: %w", err)
	}
	return &pl, nil
}

func (r *priceListRepository) FindAll(ctx context.Context) ([]entity.PriceList, error) {
	var pls []entity.PriceList
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&pls).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find price lists: %w", err)
	}
	return pls, nil
}

func (r *priceListRepository) FindActive(ctx context.Context) ([]entity.PriceList, error) {
	var pls []entity.PriceList
	err := r.db.WithContext(ctx).Where("is_active = ? AND deleted_at IS NULL", true).Find(&pls).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find active price lists: %w", err)
	}
	return pls, nil
}

func (r *priceListRepository) Update(ctx context.Context, pl *entity.PriceList) error {
	return uow.GetTx(ctx, r.db).Save(pl).Error
}

func (r *priceListRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.PriceList{}).Error
}