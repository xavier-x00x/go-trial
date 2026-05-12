package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type taxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) domainRepo.TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) Create(ctx context.Context, t *entity.Tax) error {
	return uow.GetTx(ctx, r.db).Create(t).Error
}

func (r *taxRepository) FindByID(ctx context.Context, id string) (*entity.Tax, error) {
	var tax entity.Tax
	err := r.db.WithContext(ctx).Preload("TaxAccount").Where("id = ?", id).First(&tax).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find tax: %w", err)
	}
	return &tax, nil
}

func (r *taxRepository) FindAll(ctx context.Context) ([]entity.Tax, error) {
	var taxes []entity.Tax
	err := r.db.WithContext(ctx).Preload("TaxAccount").Find(&taxes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find taxes: %w", err)
	}
	return taxes, nil
}

func (r *taxRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Tax, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Tax{})
	return PaginateAndFilter[entity.Tax](r.db, baseQuery, filter)
}

func (r *taxRepository) Update(ctx context.Context, t *entity.Tax) error {
	return uow.GetTx(ctx, r.db).Save(t).Error
}

func (r *taxRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.Tax{}).Error
}