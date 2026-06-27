package repository

import (
	"context"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type goodsReceiptRepository struct {
	db *gorm.DB
}

func NewGoodsReceiptRepository(db *gorm.DB) domainRepo.GoodsReceiptRepository {
	return &goodsReceiptRepository{db: db}
}

func (r *goodsReceiptRepository) Create(ctx context.Context, gr *entity.GoodsReceipt) error {
	return uow.GetTx(ctx, r.db).Create(gr).Error
}

func (r *goodsReceiptRepository) FindByID(ctx context.Context, id string) (*entity.GoodsReceipt, error) {
	var gr entity.GoodsReceipt
	err := r.db.WithContext(ctx).
		Preload("PurchaseOrder").
		Preload("Warehouse").
		Preload("ReceivedBy").
		Preload("ConfirmedBy").
		Preload("Items.Product").
		Preload("Items.UOM").
		Preload("Items.PurchaseOrderItem").
		Where("id = ?", id).First(&gr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find goods receipt: %w", err)
	}
	return &gr, nil
}

func (r *goodsReceiptRepository) FindByGRNumber(ctx context.Context, grNum string) (*entity.GoodsReceipt, error) {
	var gr entity.GoodsReceipt
	err := r.db.WithContext(ctx).
		Preload("PurchaseOrder").
		Preload("Warehouse").
		Preload("Items.Product").
		Preload("Items.UOM").
		Where("gr_number = ?", grNum).First(&gr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find goods receipt by number: %w", err)
	}
	return &gr, nil
}

func (r *goodsReceiptRepository) FindByPurchaseOrderID(ctx context.Context, poID string) ([]entity.GoodsReceipt, error) {
	var receipts []entity.GoodsReceipt
	err := r.db.WithContext(ctx).
		Preload("PurchaseOrder").
		Preload("Warehouse").
		Preload("Items.Product").
		Preload("Items.UOM").
		Where("purchase_order_id = ?", poID).
		Order("created_at DESC").Find(&receipts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find goods receipts: %w", err)
	}
	return receipts, nil
}

func (r *goodsReceiptRepository) FindByWarehouseID(ctx context.Context, warehouseID string, status string) ([]entity.GoodsReceipt, error) {
	var receipts []entity.GoodsReceipt
	query := r.db.WithContext(ctx).
		Preload("PurchaseOrder").
		Preload("Warehouse").
		Where("warehouse_id = ?", warehouseID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Find(&receipts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find goods receipts: %w", err)
	}
	return receipts, nil
}

func (r *goodsReceiptRepository) Update(ctx context.Context, gr *entity.GoodsReceipt) error {
	return uow.GetTx(ctx, r.db).Save(gr).Error
}

func (r *goodsReceiptRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.GoodsReceipt, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.GoodsReceipt{}).Preload("PurchaseOrder").Preload("Warehouse").Preload("ReceivedBy").Preload("ConfirmedBy")
	return PaginateAndFilter[entity.GoodsReceipt](r.db, baseQuery, filter)
}

func (r *goodsReceiptRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.GoodsReceipt{}).Error
}

func (r *goodsReceiptRepository) DeleteItemsByGoodsReceiptID(ctx context.Context, grID string) error {
	return uow.GetTx(ctx, r.db).Where("goods_receipt_id = ?", grID).Delete(&entity.GoodsReceiptItem{}).Error
}

func (r *goodsReceiptRepository) GetTotalDraftQtyByPOItemID(ctx context.Context, poItemID string, excludeGRID *string) (decimal.Decimal, error) {
	var result struct {
		Total decimal.Decimal
	}

	query := r.db.WithContext(ctx).
		Table("goods_receipt_items").
		Select("SUM(qty_received) as total").
		Joins("JOIN goods_receipts ON goods_receipts.id = goods_receipt_items.goods_receipt_id").
		Where("goods_receipt_items.purchase_order_item_id = ?", poItemID).
		Where("goods_receipts.status = ?", entity.GRStatusDraft)

	if excludeGRID != nil {
		query = query.Where("goods_receipts.id <> ?", *excludeGRID)
	}

	err := query.Scan(&result).Error
	if err != nil {
		return decimal.Zero, err
	}

	return result.Total, nil
}

func (r *goodsReceiptRepository) FindPostedItemsByDate(ctx context.Context, date time.Time) ([]entity.GoodsReceiptItem, error) {
	var items []entity.GoodsReceiptItem
	
	// Query to find goods receipt items that were received on the given date (comparing only the DATE part)
	// and the GoodsReceipt is in CONFIRMED status.
	err := r.db.WithContext(ctx).
		Preload("GoodsReceipt").
		Preload("Product").
		Preload("UOM").
		Joins("JOIN goods_receipts ON goods_receipts.id = goods_receipt_items.goods_receipt_id").
		Where("DATE(goods_receipts.receipt_date) = DATE(?)", date).
		Where("goods_receipts.status = ?", entity.GRStatusConfirmed).
		Find(&items).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find posted goods receipt items by date: %w", err)
	}

	return items, nil
}

func (r *goodsReceiptRepository) FindLastPriceBeforeDate(ctx context.Context, productID, uomID string, date time.Time) (*decimal.Decimal, error) {
	var items []entity.GoodsReceiptItem
	err := r.db.WithContext(ctx).
		Joins("JOIN goods_receipts ON goods_receipts.id = goods_receipt_items.goods_receipt_id").
		Where("goods_receipts.status = ?", entity.GRStatusConfirmed).
		Where("goods_receipt_items.product_id = ? AND goods_receipt_items.uom_id = ?", productID, uomID).
		Where("DATE(goods_receipts.receipt_date) < DATE(?)", date).
		Order("goods_receipts.receipt_date DESC, goods_receipts.created_at DESC").
		Limit(1).
		Find(&items).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find last purchase price: %w", err)
	}

	if len(items) == 0 {
		return nil, nil
	}

	return &items[0].NetUnitPrice, nil
}