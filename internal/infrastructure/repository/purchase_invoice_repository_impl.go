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

type purchaseInvoiceRepository struct {
	db *gorm.DB
}

func NewPurchaseInvoiceRepository(db *gorm.DB) domainRepo.PurchaseInvoiceRepository {
	return &purchaseInvoiceRepository{db: db}
}

func (r *purchaseInvoiceRepository) Create(ctx context.Context, inv *entity.PurchaseInvoice) error {
	return uow.GetTx(ctx, r.db).Create(inv).Error
}

func (r *purchaseInvoiceRepository) Update(ctx context.Context, inv *entity.PurchaseInvoice) error {
	return uow.GetTx(ctx, r.db).Save(inv).Error
}

func (r *purchaseInvoiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.PurchaseInvoice{}).Error
}

func (r *purchaseInvoiceRepository) DeleteItemsByPurchaseInvoiceID(ctx context.Context, invoiceID string) error {
	return uow.GetTx(ctx, r.db).Where("purchase_invoice_id = ?", invoiceID).Delete(&entity.PurchaseInvoiceItem{}).Error
}

func (r *purchaseInvoiceRepository) FindByID(ctx context.Context, id string) (*entity.PurchaseInvoice, error) {
	var inv entity.PurchaseInvoice
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Preload("APAccount").
		Preload("CreatedBy").
		Preload("VerifiedBy").
		Preload("PostedBy").
		Preload("Items.Product").
		Preload("Items.UOM").
		Preload("Items.PurchaseOrderItem").
		Preload("Items.GoodsReceiptItem").
		Where("id = ?", id).First(&inv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase invoice: %w", err)
	}
	return &inv, nil
}

func (r *purchaseInvoiceRepository) FindByInvoiceNumber(ctx context.Context, invoiceNum string) (*entity.PurchaseInvoice, error) {
	var inv entity.PurchaseInvoice
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Preload("Items").
		Where("invoice_number = ?", invoiceNum).First(&inv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find purchase invoice by number: %w", err)
	}
	return &inv, nil
}

func (r *purchaseInvoiceRepository) FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseInvoice, error) {
	var invs []entity.PurchaseInvoice
	query := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Where("store_id = ?", storeID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&invs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find purchase invoices: %w", err)
	}
	return invs, nil
}

func (r *purchaseInvoiceRepository) FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseInvoice, error) {
	var invs []entity.PurchaseInvoice
	err := r.db.WithContext(ctx).
		Preload("Supplier").
		Preload("Store").
		Where("store_id = ? AND status IN ?", storeID, []string{entity.PurchaseInvoiceStatusPosted}).
		Order("due_date ASC").Find(&invs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find pending purchase invoices: %w", err)
	}
	return invs, nil
}

func (r *purchaseInvoiceRepository) FindAllWithPagination(ctx context.Context, filter *domainRepo.QueryFilter) ([]entity.PurchaseInvoice, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.PurchaseInvoice{}).
		Preload("Supplier").
		Preload("Store").
		Preload("Warehouse").
		Preload("CreatedBy")
	return PaginateAndFilter[entity.PurchaseInvoice](r.db, baseQuery, entity.QueryFilter{
		Page:         filter.Page,
		Limit:        filter.Limit,
		OrderColumn:  filter.OrderBy,
		OrderDir:     filter.OrderDir,
		Search:       filter.Search,
		SearchColumn: filter.SearchColumns,
		Conditions:   filter.Conditions,
	})
}