package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type purchasePaymentRepository struct {
	db *gorm.DB
}

func NewPurchasePaymentRepository(db *gorm.DB) domainRepo.PurchasePaymentRepository {
	return &purchasePaymentRepository{db: db}
}

func (r *purchasePaymentRepository) Create(ctx context.Context, pp *entity.PurchasePayment) error {
	return uow.GetTx(ctx, r.db).Create(pp).Error
}

func (r *purchasePaymentRepository) Update(ctx context.Context, pp *entity.PurchasePayment) error {
	return uow.GetTx(ctx, r.db).Save(pp).Error
}

func (r *purchasePaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.PurchasePayment{}).Error
}

func (r *purchasePaymentRepository) FindByID(ctx context.Context, id string) (*entity.PurchasePayment, error) {
	var pp entity.PurchasePayment
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("PaymentAccount").
		Preload("APAccount").
		Preload("CreatedBy").
		Preload("PostedBy").
		Preload("Items").
		Preload("Items.PurchaseInvoice").
		Where("id = ?", id).First(&pp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase payment: %w", err)
	}
	return &pp, nil
}

func (r *purchasePaymentRepository) FindByPaymentNumber(ctx context.Context, paymentNum string) (*entity.PurchasePayment, error) {
	var pp entity.PurchasePayment
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("PaymentAccount").
		Preload("APAccount").
		Preload("Items").
		Where("payment_number = ?", paymentNum).First(&pp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase payment by number: %w", err)
	}
	return &pp, nil
}

func (r *purchasePaymentRepository) FindBySupplierID(ctx context.Context, supplierID string) ([]entity.PurchasePayment, error) {
	var pps []entity.PurchasePayment
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("PaymentAccount").
		Where("supplier_id = ?", supplierID).
		Order("created_at DESC").Find(&pps).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find purchase payments by supplier: %w", err)
	}
	return pps, nil
}

func (r *purchasePaymentRepository) FindAllWithPagination(ctx context.Context, filter *domainRepo.QueryFilter) ([]entity.PurchasePayment, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.PurchasePayment{}).
		Preload("Supplier").
		Preload("PaymentAccount").
		Preload("CreatedBy")
	return PaginateAndFilter[entity.PurchasePayment](r.db, baseQuery, entity.QueryFilter{
		Page:         filter.Page,
		Limit:        filter.Limit,
		OrderColumn:  filter.OrderBy,
		OrderDir:     filter.OrderDir,
		Search:       filter.Search,
		SearchColumn: filter.SearchColumns,
		Conditions:   filter.Conditions,
	})
}