package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type productCategoryRepository struct {
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) domainRepo.ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) Create(ctx context.Context, cat *entity.ProductCategory) error {
	return uow.GetTx(ctx, r.db).Create(cat).Error
}

func (r *productCategoryRepository) FindByID(ctx context.Context, id string) (*entity.ProductCategory, error) {
	var cat entity.ProductCategory
	err := r.db.WithContext(ctx).Preload("Parent").Where("id = ?", id).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find category: %w", err)
	}
	return &cat, nil
}

func (r *productCategoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.ProductCategory, error) {
	var cat entity.ProductCategory
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find category by slug: %w", err)
	}
	return &cat, nil
}

func (r *productCategoryRepository) FindAll(ctx context.Context) ([]entity.ProductCategory, error) {
	var cats []entity.ProductCategory
	err := r.db.WithContext(ctx).Preload("Parent").Find(&cats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	return cats, nil
}

func (r *productCategoryRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.ProductCategory, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.ProductCategory{}).Preload("Parent")
	return PaginateAndFilter[entity.ProductCategory](r.db, baseQuery, filter)
}

func (r *productCategoryRepository) FindByParentID(ctx context.Context, parentID string) ([]entity.ProductCategory, error) {
	var cats []entity.ProductCategory
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&cats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find categories by parent: %w", err)
	}
	return cats, nil
}

func (r *productCategoryRepository) Update(ctx context.Context, cat *entity.ProductCategory) error {
	return uow.GetTx(ctx, r.db).Save(cat).Error
}

func (r *productCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.ProductCategory{}).Error
}