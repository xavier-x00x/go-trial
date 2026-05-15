package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type purchaseReturnRepository struct {
	db *gorm.DB
}

func NewPurchaseReturnRepository(db *gorm.DB) domainRepo.PurchaseReturnRepository {
	return &purchaseReturnRepository{db: db}
}

func (r *purchaseReturnRepository) Create(ctx context.Context, pr *entity.PurchaseReturn) error {
	return uow.GetTx(ctx, r.db).Create(pr).Error
}

func (r *purchaseReturnRepository) Update(ctx context.Context, pr *entity.PurchaseReturn) error {
	return uow.GetTx(ctx, r.db).Save(pr).Error
}

func (r *purchaseReturnRepository) FindByID(ctx context.Context, id string) (*entity.PurchaseReturn, error) {
	var pr entity.PurchaseReturn
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Preload("PurchaseInvoice").
		Preload("Items.Product").
		Preload("Items.UOM").
		Preload("Items.PurchaseInvoiceItem").
		Where("id = ?", id).First(&pr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase return: %w", err)
	}
	return &pr, nil
}

func (r *purchaseReturnRepository) FindByReturnNumber(ctx context.Context, returnNum string) (*entity.PurchaseReturn, error) {
	var pr entity.PurchaseReturn
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Items").
		Where("return_number = ?", returnNum).First(&pr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase return by number: %w", err)
	}
	return &pr, nil
}

func (r *purchaseReturnRepository) FindAllWithPagination(ctx context.Context, filter *domainRepo.QueryFilter) ([]entity.PurchaseReturn, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.PurchaseReturn{}).
		Preload("Supplier").
		Preload("Store").
		Preload("PurchaseInvoice")

	return PaginateAndFilter[entity.PurchaseReturn](r.db, baseQuery, entity.QueryFilter{
		Page:         filter.Page,
		Limit:        filter.Limit,
		OrderColumn:  filter.OrderBy,
		OrderDir:     filter.OrderDir,
		Search:       filter.Search,
		SearchColumn: filter.SearchColumns,
		Conditions:   filter.Conditions,
	})
}
