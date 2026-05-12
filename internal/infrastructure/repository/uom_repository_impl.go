package repository

import (
	"context"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type uomRepository struct {
	db *gorm.DB
}

func NewUOMRepository(db *gorm.DB) domainRepo.UOMRepository {
	return &uomRepository{db: db}
}

func (r *uomRepository) Create(ctx context.Context, uom *entity.UOM) error {
	return uow.GetTx(ctx, r.db).Create(uom).Error
}

func (r *uomRepository) FindByID(ctx context.Context, id string) (*entity.UOM, error) {
	var uom entity.UOM
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&uom).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find uom: %w", err)
	}
	return &uom, nil
}

func (r *uomRepository) FindByCode(ctx context.Context, code string) (*entity.UOM, error) {
	var uom entity.UOM
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&uom).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find uom by code: %w", err)
	}
	return &uom, nil
}

func (r *uomRepository) FindAll(ctx context.Context) ([]entity.UOM, error) {
	var uoms []entity.UOM
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&uoms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find uoms: %w", err)
	}
	return uoms, nil
}

func (r *uomRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.UOM, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.UOM{})
	return PaginateAndFilter[entity.UOM](r.db, baseQuery, filter)
}

func (r *uomRepository) Update(ctx context.Context, uom *entity.UOM) error {
	return uow.GetTx(ctx, r.db).Save(uom).Error
}

func (r *uomRepository) Delete(ctx context.Context, id string) error {
	var uom entity.UOM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&uom).Error; err != nil {
		return err
	}
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	return r.db.WithContext(ctx).Model(&uom).Updates(map[string]interface{}{
		"code":      uom.Code + suffix,
		"deleted_at": time.Now(),
	}).Error
}