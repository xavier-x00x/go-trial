package repository

import (
	"context"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domainRepo.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *entity.Product) error {
	return uow.GetTx(ctx, r.db).Create(p).Error
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", sku).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product by sku: %w", err)
	}
	return &product, nil
}

func (r *productRepository) FindByBarcode(ctx context.Context, barcode string) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Where("barcode = ? AND deleted_at IS NULL", barcode).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find product by barcode: %w", err)
	}
	return &product, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Preload("Category").Preload("BaseUOM").Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}
	return products, nil
}

func (r *productRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Product, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Product{})
	return PaginateAndFilter[entity.Product](r.db, baseQuery, filter)
}

func (r *productRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.WithContext(ctx).Where("category_id = ? AND deleted_at IS NULL", categoryID).Find(&products).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find products by category: %w", err)
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, p *entity.Product) error {
	return uow.GetTx(ctx, r.db).Save(p).Error
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	var product entity.Product
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
		return err
	}
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	updates := map[string]interface{}{
		"sku":       product.SKU + suffix,
		"deleted_at": time.Now(),
	}
	if product.Barcode != nil {
		updates["barcode"] = *product.Barcode + suffix
	}
	return r.db.WithContext(ctx).Model(&product).Updates(updates).Error
}