package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrGoodsReceiptNotFound = errors.New("goods receipt not found")
	ErrGRInvalidStatus      = errors.New("invalid goods receipt status transition")
	ErrGRItemNotFound       = errors.New("gr item not found")
	ErrGRItemMismatch       = errors.New("item does not match purchase order item")
	ErrGRNoItemsReceived    = errors.New("no items received")
	ErrGRNoRejectReason     = errors.New("reject reason required when qty rejected > 0")
)

type GoodsReceiptUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreateGoodsReceiptRequest) (*dto.GoodsReceiptDetailResponse, error)
	Confirm(ctx context.Context, userID string, id string, req dto.ConfirmGoodsReceiptRequest) error
	Cancel(ctx context.Context, userID string, id string, reason string) error
	GetByID(ctx context.Context, id string) (*dto.GoodsReceiptDetailResponse, error)
	GetByPurchaseOrderID(ctx context.Context, poID string) ([]dto.GoodsReceiptListResponse, error)
	GetByWarehouseID(ctx context.Context, warehouseID string, status string) ([]dto.GoodsReceiptListResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.GoodsReceiptListResponse, *entity.Meta, error)
	CreateWithInvoice(ctx context.Context, userID string, req dto.CreateGoodsReceiptWithInvoiceRequest) (*dto.GoodsReceiptDetailResponse, error)
}

type GoodsReceiptConfig struct {
	GRRepo             repository.GoodsReceiptRepository
	PurchaseOrderRepo  repository.PurchaseOrderRepository
	InventoryStockRepo repository.InventoryStockRepository
	ProductRepo        repository.ProductRepository
	ProductUOMRepo     repository.ProductUOMConversionRepository
	UserRepo           repository.UserRepository
	WarehouseRepo      repository.WarehouseRepository
	NumberSequenceRepo repository.NumberSequenceRepository
	Uow                interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

type goodsReceiptUseCaseImpl struct {
	grRepo             repository.GoodsReceiptRepository
	purchaseOrderRepo  repository.PurchaseOrderRepository
	inventoryStockRepo repository.InventoryStockRepository
	productRepo        repository.ProductRepository
	productUOMRepo     repository.ProductUOMConversionRepository
	userRepo           repository.UserRepository
	warehouseRepo      repository.WarehouseRepository
	numberSequenceRepo repository.NumberSequenceRepository
	uow                interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

func NewGoodsReceiptUseCase(cfg GoodsReceiptConfig) GoodsReceiptUseCase {
	return &goodsReceiptUseCaseImpl{
		grRepo:             cfg.GRRepo,
		purchaseOrderRepo:  cfg.PurchaseOrderRepo,
		inventoryStockRepo: cfg.InventoryStockRepo,
		productRepo:        cfg.ProductRepo,
		productUOMRepo:     cfg.ProductUOMRepo,
		userRepo:           cfg.UserRepo,
		warehouseRepo:      cfg.WarehouseRepo,
		numberSequenceRepo: cfg.NumberSequenceRepo,
		uow:                cfg.Uow,
	}
}

func (u *goodsReceiptUseCaseImpl) generateGRNumber(date time.Time) (string, error) {
	prefix := "GR"
	period := date.Format("0601")

	seqNum, err := u.numberSequenceRepo.GetNextNumber(context.Background(), prefix, period)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%05d", prefix, period, seqNum), nil
}

func (u *goodsReceiptUseCaseImpl) convertToBaseUOM(ctx context.Context, productID, fromUOMID uuid.UUID, qty decimal.Decimal) (decimal.Decimal, error) {
	if fromUOMID.String() == "" {
		return qty, nil
	}

	conversions, err := u.productUOMRepo.FindByProductID(ctx, productID.String())
	if err != nil {
		return decimal.Zero, err
	}

	for _, conv := range conversions {
		if conv.UOMID == fromUOMID {
			return qty.Mul(conv.ConversionRate), nil
		}
	}

	return qty, nil
}

func (u *goodsReceiptUseCaseImpl) calculateWeightedAverage(existingQty, existingAvg, newQty, newPrice decimal.Decimal) decimal.Decimal {
	totalExistingValue := existingQty.Mul(existingAvg)
	totalNewValue := newQty.Mul(newPrice)
	totalQty := existingQty.Add(newQty)

	if totalQty.IsZero() {
		return newPrice
	}

	return totalExistingValue.Add(totalNewValue).Div(totalQty)
}

func (u *goodsReceiptUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreateGoodsReceiptRequest) (*dto.GoodsReceiptDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	if len(req.Items) == 0 {
		return nil, ErrGRNoItemsReceived
	}

	po, err := u.purchaseOrderRepo.FindByID(ctx, req.PurchaseOrderID.String())
	if err != nil || po == nil {
		return nil, ErrPurchaseOrderNotFound
	}

	if po.Status != entity.POStatusApproved && po.Status != entity.POStatusPartiallyReceived {
		return nil, ErrPOInvalidStatus
	}

	warehouse, err := u.warehouseRepo.FindByID(ctx, req.WarehouseID.String())
	if err != nil || warehouse == nil {
		return nil, errors.New("warehouse not found")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	supplier, err := u.purchaseOrderRepo.FindByIDWithSupplier(ctx, po.ID.String())
	if err != nil || supplier == nil {
		return nil, errors.New("supplier not found")
	}

	poItemsMap := make(map[string]entity.PurchaseOrderItem)
	for _, item := range po.Items {
		poItemsMap[item.ID.String()] = item
	}

	grNum, err := u.generateGRNumber(time.Now())
	if err != nil {
		return nil, err
	}

	items := make([]entity.GoodsReceiptItem, len(req.Items))
	for i, item := range req.Items {
		poItem, exists := poItemsMap[item.PurchaseOrderItemID.String()]
		if !exists {
			return nil, ErrGRItemMismatch
		}

		if item.QtyRejected.GreaterThan(decimal.Zero) && item.RejectReason == nil {
			return nil, ErrGRNoRejectReason
		}

		items[i] = entity.GoodsReceiptItem{
			SeqNo:               i + 1,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			ProductID:           item.ProductID,
			UOMID:               item.UOMID,
			QtyOrdered:          poItem.QtyOrdered,
			QtyReceived:         item.QtyReceived,
			QtyRejected:         item.QtyRejected,
			UnitPrice:           poItem.UnitPrice,
			Discount1Pct:        decimal.Zero,
			Discount2Pct:        decimal.Zero,
			Discount3Pct:        decimal.Zero,
			DiscountAmount:      decimal.Zero,
			TotalDiscountAmount: decimal.Zero,
			TaxPct:              decimal.Zero,
			TaxAmount:           decimal.Zero,
			LandedCostAmount:    decimal.Zero,
			NetUnitPrice:        poItem.UnitPrice, // Default to gross price since PO has no financials
			RejectReason:        item.RejectReason,
			Notes:               item.Notes,
		}
	}

	gr := &entity.GoodsReceipt{
		GRNumber:        grNum,
		PurchaseOrderID: req.PurchaseOrderID,
		WarehouseID:     req.WarehouseID,
		ReceiptDate:     req.ReceiptDate,
		DeliveryNoteNo:  req.DeliveryNoteNo,
		Status:          entity.GRStatusDraft,
		ReceivedByID:    userUUID,
		Notes:           req.Notes,
		Items:           items,
		WarehouseName:   warehouse.Name,
		SupplierName:    po.SupplierName,
		SupplierCode:    po.SupplierCode,
		SupplierAddress: po.SupplierAddress,
		StoreName:       po.StoreName,
		ReceivedByName:  user.Name,
		Subtotal:        po.TotalAmount,
		DiscountAmount:  decimal.Zero,
		TaxAmount:       decimal.Zero,
		FreightAmount:   decimal.Zero,
		OtherCostAmount: decimal.Zero,
		GrandTotal:      po.TotalAmount,
		IsTaxInclusive:  false,
	}

	if err := gr.GenerateID(); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].GoodsReceiptID = gr.ID
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.grRepo.Create(txCtx, gr); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toGRDetailResponse(gr), nil
}

func (u *goodsReceiptUseCaseImpl) Confirm(ctx context.Context, userID string, id string, req dto.ConfirmGoodsReceiptRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	gr, err := u.grRepo.FindByID(ctx, id)
	if err != nil || gr == nil {
		return ErrGoodsReceiptNotFound
	}

	if gr.Status != entity.GRStatusDraft {
		return ErrGRInvalidStatus
	}

	confirmer, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || confirmer == nil {
		return errors.New("user not found")
	}

	gr.Status = entity.GRStatusConfirmed
	gr.ConfirmedByID = &userUUID
	now := time.Now()
	gr.ConfirmedAt = &now
	gr.ConfirmedByName = &confirmer.Name
	if req.Notes != nil {
		gr.Notes = req.Notes
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	for _, item := range gr.Items {
		baseQty, err := u.convertToBaseUOM(txCtx, item.ProductID, item.UOMID, item.QtyReceived)
		if err != nil {
			return err
		}

		stock, err := u.inventoryStockRepo.FindByWarehouseAndProduct(txCtx, gr.WarehouseID.String(), item.ProductID.String())
		if err != nil {
			return err
		}

		if stock == nil {
			stock = &entity.InventoryStock{
				WarehouseID:     gr.WarehouseID,
				ProductID:       item.ProductID,
				Quantity:        decimal.Zero,
				ReservedQty:     decimal.Zero,
				AverageBuyPrice: decimal.Zero,
			}
			if err := stock.GenerateID(); err != nil {
				return err
			}
		}

		oldQty := stock.Quantity
		oldAvg := stock.AverageBuyPrice
		newPrice := item.NetUnitPrice // Use NetUnitPrice (after discounts and taxes) for HPP

		stock.Quantity = stock.Quantity.Add(baseQty)
		stock.AverageBuyPrice = u.calculateWeightedAverage(oldQty, oldAvg, baseQty, newPrice)

		if err := u.inventoryStockRepo.Update(txCtx, stock); err != nil {
			return err
		}

		poItem, err := u.purchaseOrderRepo.FindByID(txCtx, gr.PurchaseOrderID.String())
		if err != nil || poItem == nil {
			return ErrPurchaseOrderNotFound
		}

		for i := range poItem.Items {
			if poItem.Items[i].ID == item.PurchaseOrderItemID {
				poItem.Items[i].QtyReceived = poItem.Items[i].QtyReceived.Add(item.QtyReceived)
				break
			}
		}

		if err := u.purchaseOrderRepo.Update(txCtx, poItem); err != nil {
			return err
		}
	}

	totalReceived := decimal.Zero
	totalOrdered := decimal.Zero
	for _, item := range gr.Items {
		totalReceived = totalReceived.Add(item.QtyReceived)
		totalOrdered = totalOrdered.Add(item.QtyOrdered)
	}

	if totalReceived.GreaterThanOrEqual(totalOrdered) {
		po, err := u.purchaseOrderRepo.FindByID(txCtx, gr.PurchaseOrderID.String())
		if err == nil && po != nil {
			po.Status = entity.POStatusReceived
			if err := u.purchaseOrderRepo.Update(txCtx, po); err != nil {
				return err
			}
		}
	} else {
		po, err := u.purchaseOrderRepo.FindByID(txCtx, gr.PurchaseOrderID.String())
		if err == nil && po != nil {
			po.Status = entity.POStatusPartiallyReceived
			if err := u.purchaseOrderRepo.Update(txCtx, po); err != nil {
				return err
			}
		}
	}

	if err := u.grRepo.Update(txCtx, gr); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *goodsReceiptUseCaseImpl) Cancel(ctx context.Context, userID string, id string, reason string) error {
	gr, err := u.grRepo.FindByID(ctx, id)
	if err != nil || gr == nil {
		return ErrGoodsReceiptNotFound
	}

	if gr.Status != entity.GRStatusDraft {
		return ErrGRInvalidStatus
	}

	gr.Status = entity.GRStatusCancelled

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.grRepo.Update(txCtx, gr); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *goodsReceiptUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.GoodsReceiptDetailResponse, error) {
	gr, err := u.grRepo.FindByID(ctx, id)
	if err != nil || gr == nil {
		return nil, ErrGoodsReceiptNotFound
	}

	return toGRDetailResponse(gr), nil
}

func (u *goodsReceiptUseCaseImpl) GetByPurchaseOrderID(ctx context.Context, poID string) ([]dto.GoodsReceiptListResponse, error) {
	receipts, err := u.grRepo.FindByPurchaseOrderID(ctx, poID)
	if err != nil {
		return nil, err
	}

	return toGRListResponses(receipts), nil
}

func (u *goodsReceiptUseCaseImpl) GetByWarehouseID(ctx context.Context, warehouseID string, status string) ([]dto.GoodsReceiptListResponse, error) {
	receipts, err := u.grRepo.FindByWarehouseID(ctx, warehouseID, status)
	if err != nil {
		return nil, err
	}

	return toGRListResponses(receipts), nil
}

func (u *goodsReceiptUseCaseImpl) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.GoodsReceiptListResponse, *entity.Meta, error) {
	allowedOrder := []string{"created_at", "updated_at", "gr_number", "receipt_date", "status"}
	searchColumns := []string{"gr_number", "delivery_note_no"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	if meta.Conditions != nil {
		filter.Conditions = meta.Conditions
	}

	data, resMeta, err := u.grRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return toGRListResponses(data), resMeta, nil
}

func (u *goodsReceiptUseCaseImpl) CreateWithInvoice(ctx context.Context, userID string, req dto.CreateGoodsReceiptWithInvoiceRequest) (*dto.GoodsReceiptDetailResponse, error) {
	createReq := dto.CreateGoodsReceiptRequest{
		PurchaseOrderID: req.PurchaseOrderID,
		WarehouseID:     req.WarehouseID,
		ReceiptDate:     req.ReceiptDate,
		DeliveryNoteNo:  req.DeliveryNoteNo,
		Notes:           req.Notes,
		Items:           req.Items,
	}

	gr, err := u.Create(ctx, userID, createReq)
	if err != nil {
		return nil, err
	}

	confirmReq := dto.ConfirmGoodsReceiptRequest{
		Notes: nil,
	}

	if err := u.Confirm(ctx, userID, gr.ID.String(), confirmReq); err != nil {
		return nil, err
	}

	return u.GetByID(ctx, gr.ID.String())
}

func toGRListResponses(grs []entity.GoodsReceipt) []dto.GoodsReceiptListResponse {
	responses := make([]dto.GoodsReceiptListResponse, len(grs))
	for i, gr := range grs {
		responses[i] = dto.GoodsReceiptListResponse{
			ID:              gr.ID,
			GRNumber:        gr.GRNumber,
			PurchaseOrderID: gr.PurchaseOrderID,
			PONumber:        gr.PurchaseOrder.PONumber,
			WarehouseID:     gr.WarehouseID,
			WarehouseName:   gr.Warehouse.Name,
			ReceiptDate:     gr.ReceiptDate,
			DeliveryNoteNo:  gr.DeliveryNoteNo,
			Status:          gr.Status,
			ReceivedByID:    gr.ReceivedByID,
			SupplierName:    gr.SupplierName,
			StoreName:       gr.StoreName,
			CreatedAt:       gr.CreatedAt,
		}
	}
	return responses
}

func toGRDetailResponse(gr *entity.GoodsReceipt) *dto.GoodsReceiptDetailResponse {
	items := make([]dto.GoodsReceiptItemResponse, len(gr.Items))
	for i, item := range gr.Items {
		items[i] = dto.GoodsReceiptItemResponse{
			ID:                  item.ID,
			SeqNo:               item.SeqNo,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			ProductID:           item.ProductID,
			ProductName:         item.Product.Name,
			ProductSKU:          item.Product.SKU,
			UOMID:               item.UOMID,
			UOMCode:             item.UOM.Code,
			QtyOrdered:          item.QtyOrdered,
			QtyReceived:         item.QtyReceived,
			QtyRejected:         item.QtyRejected,
			UnitPrice:           item.UnitPrice,
			Discount1Pct:        item.Discount1Pct,
			Discount2Pct:        item.Discount2Pct,
			Discount3Pct:        item.Discount3Pct,
			DiscountAmount:      item.DiscountAmount,
			TotalDiscountAmount: item.TotalDiscountAmount,
			TaxPct:              item.TaxPct,
			TaxAmount:           item.TaxAmount,
			LandedCostAmount:    item.LandedCostAmount,
			NetUnitPrice:        item.NetUnitPrice,
			RejectReason:        item.RejectReason,
			Notes:               item.Notes,
		}
	}

	response := &dto.GoodsReceiptDetailResponse{
		ID:              gr.ID,
		GRNumber:        gr.GRNumber,
		PurchaseOrderID: gr.PurchaseOrderID,
		PONumber:        gr.PurchaseOrder.PONumber,
		WarehouseID:     gr.WarehouseID,
		WarehouseName:   gr.Warehouse.Name,
		ReceiptDate:     gr.ReceiptDate,
		DeliveryNoteNo:  gr.DeliveryNoteNo,
		Status:          gr.Status,
		ReceivedByID:    gr.ReceivedByID,
		ConfirmedByID:   gr.ConfirmedByID,
		ConfirmedAt:     gr.ConfirmedAt,
		Notes:           gr.Notes,
		SupplierName:    gr.SupplierName,
		SupplierCode:    gr.SupplierCode,
		SupplierAddress: gr.SupplierAddress,
		StoreName:       gr.StoreName,
		Subtotal:        gr.Subtotal,
		DiscountAmount:  gr.DiscountAmount,
		TaxAmount:       gr.TaxAmount,
		FreightAmount:   gr.FreightAmount,
		OtherCostAmount: gr.OtherCostAmount,
		GrandTotal:      gr.GrandTotal,
		IsTaxInclusive:  gr.IsTaxInclusive,
		CreatedAt:       gr.CreatedAt,
		UpdatedAt:       gr.UpdatedAt,
		Items:           items,
	}

	return response
}
