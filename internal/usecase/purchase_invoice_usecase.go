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
	ErrPurchaseInvoiceNotFound = errors.New("purchase invoice not found")
	ErrPIInvalidStatus         = errors.New("invalid purchase invoice status transition")
	ErrPIIsAlreadyVerified     = errors.New("invoice already verified")
	ErrPIIsAlreadyPosted       = errors.New("invoice already posted")
)

type PurchaseInvoiceUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreatePurchaseInvoiceRequest) (*dto.PurchaseInvoiceDetailResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PurchaseInvoiceDetailResponse, error)
	GetByInvoiceNumber(ctx context.Context, invoiceNum string) (*dto.PurchaseInvoiceDetailResponse, error)
	GetByStoreID(ctx context.Context, storeID string, status string) ([]dto.PurchaseInvoiceListResponse, error)
	GetPendingByStoreID(ctx context.Context, storeID string) ([]dto.PurchaseInvoiceListResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PurchaseInvoiceListResponse, *entity.Meta, error)
	Submit(ctx context.Context, userID string, id string) error
	Approve(ctx context.Context, userID string, id string) error
	Verify(ctx context.Context, userID string, id string) error
	Post(ctx context.Context, userID string, id string) error
	Update(ctx context.Context, userID string, id string, req dto.UpdatePurchaseInvoiceRequest) (*dto.PurchaseInvoiceDetailResponse, error)
	Cancel(ctx context.Context, userID string, id string, reason string) error
}

type PurchaseInvoiceConfig struct {
	Repo               repository.PurchaseInvoiceRepository
	PurchaseOrderRepo  repository.PurchaseOrderRepository
	SupplierRepo       repository.SupplierRepository
	StoreRepo          repository.StoreRepository
	WarehouseRepo      repository.WarehouseRepository
	UserRepo           repository.UserRepository
	NumberSequenceRepo repository.NumberSequenceRepository
	Uow                interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

type purchaseInvoiceUseCaseImpl struct {
	repo               repository.PurchaseInvoiceRepository
	purchaseOrderRepo  repository.PurchaseOrderRepository
	supplierRepo       repository.SupplierRepository
	storeRepo          repository.StoreRepository
	warehouseRepo      repository.WarehouseRepository
	userRepo           repository.UserRepository
	numberSequenceRepo repository.NumberSequenceRepository
	uow                interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

func NewPurchaseInvoiceUseCase(cfg PurchaseInvoiceConfig) PurchaseInvoiceUseCase {
	return &purchaseInvoiceUseCaseImpl{
		repo:               cfg.Repo,
		purchaseOrderRepo:  cfg.PurchaseOrderRepo,
		supplierRepo:       cfg.SupplierRepo,
		storeRepo:          cfg.StoreRepo,
		warehouseRepo:      cfg.WarehouseRepo,
		userRepo:           cfg.UserRepo,
		numberSequenceRepo: cfg.NumberSequenceRepo,
		uow:                cfg.Uow,
	}
}

func (u *purchaseInvoiceUseCaseImpl) generateInvoiceNumber(date time.Time) (string, error) {
	prefix := "PI"
	period := date.Format("0601")

	seqNum, err := u.numberSequenceRepo.GetNextNumber(context.Background(), prefix, period)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%05d", prefix, period, seqNum), nil
}

func (u *purchaseInvoiceUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreatePurchaseInvoiceRequest) (*dto.PurchaseInvoiceDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	if len(req.Items) == 0 {
		return nil, errors.New("items is required")
	}

	supplier, err := u.supplierRepo.FindByID(ctx, req.SupplierID.String())
	if err != nil || supplier == nil {
		return nil, errors.New("supplier not found")
	}

	store, err := u.storeRepo.FindByID(ctx, req.StoreID.String())
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	warehouse, err := u.warehouseRepo.FindByID(ctx, req.WarehouseID.String())
	if err != nil || warehouse == nil {
		return nil, errors.New("warehouse not found")
	}

	po, err := u.purchaseOrderRepo.FindByID(ctx, req.PurchaseOrderID.String())
	if err != nil || po == nil {
		return nil, errors.New("purchase order not found")
	}

	creator, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || creator == nil {
		return nil, errors.New("creator user not found")
	}

	invoiceDate := req.InvoiceDate
	if invoiceDate.IsZero() {
		invoiceDate = time.Now()
	}

	invoiceNum, err := u.generateInvoiceNumber(invoiceDate)
	if err != nil {
		return nil, err
	}

	dueDate := invoiceDate.AddDate(0, 0, req.PaymentTermDays)

	subtotal := decimal.Zero
	taxAmount := decimal.Zero

	items := make([]entity.PurchaseInvoiceItem, len(req.Items))
	for i, item := range req.Items {
		qty := item.QtyInvoiced
		unitPrice := item.UnitPrice
		lineSubtotal := qty.Mul(unitPrice)

		discount1 := lineSubtotal.Mul(item.Discount1Pct.Div(decimal.NewFromInt(100)))
		discount2 := lineSubtotal.Sub(discount1).Mul(item.Discount2Pct.Div(decimal.NewFromInt(100)))
		discount3 := lineSubtotal.Sub(discount1).Sub(discount2).Mul(item.Discount3Pct.Div(decimal.NewFromInt(100)))
		totalDiscount := discount1.Add(discount2).Add(discount3).Add(item.DiscountAmount)

		afterDiscount := lineSubtotal.Sub(totalDiscount)
		itemTaxAmount := afterDiscount.Mul(item.TaxPct.Div(decimal.NewFromInt(100)))

		subtotal = subtotal.Add(afterDiscount)
		taxAmount = taxAmount.Add(itemTaxAmount)

		items[i] = entity.PurchaseInvoiceItem{
			SeqNo:               i + 1,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			GoodsReceiptItemID:  item.GoodsReceiptItemID,
			ProductID:           item.ProductID,
			UOMID:               item.UOMID,
			QtyInvoiced:         qty,
			UnitPrice:           unitPrice,
			Discount1Pct:        item.Discount1Pct,
			Discount2Pct:        item.Discount2Pct,
			Discount3Pct:        item.Discount3Pct,
			DiscountAmount:      item.DiscountAmount,
			TotalDiscountAmount: totalDiscount,
			TaxPct:              item.TaxPct,
			TaxAmount:           itemTaxAmount,
			Subtotal:            afterDiscount,
			Notes:               item.Notes,
		}
	}

	grandTotal := subtotal.Add(taxAmount).Add(req.FreightAmount).Add(req.OtherCostAmount).Sub(req.DiscountAmount)

	pi := &entity.PurchaseInvoice{
		InvoiceNumber:         invoiceNum,
		SupplierInvoiceNumber: req.SupplierInvoiceNumber,
		ReferenceNo:           req.ReferenceNo,
		PurchaseOrderID:       req.PurchaseOrderID,
		SupplierID:            req.SupplierID,
		StoreID:               req.StoreID,
		WarehouseID:           req.WarehouseID,
		APAccountID:           req.APAccountID,
		InvoiceDate:           invoiceDate,
		ReceivedDate:          req.ReceivedDate,
		DueDate:               dueDate,
		PaymentTermDays:       req.PaymentTermDays,
		PaymentMode:           req.PaymentMode,
		Subtotal:              subtotal,
		DiscountAmount:        req.DiscountAmount,
		TaxAmount:             taxAmount,
		FreightAmount:         req.FreightAmount,
		OtherCostAmount:       req.OtherCostAmount,
		GrandTotal:            grandTotal,
		IsTaxInclusive:        req.IsTaxInclusive,
		PaidAmount:            decimal.Zero,
		RemainingAmount:       grandTotal,
		Status:                entity.PurchaseInvoiceStatusDraft,
		CreatedByID:           userUUID,
		Notes:                 req.Notes,
		SupplierName:          supplier.Name,
		SupplierCode:          supplier.Code,
		SupplierAddress:       supplier.Address,
		StoreCode:             store.Code,
		StoreName:             store.Name,
		StoreAddress:          store.Address,
		WarehouseName:         warehouse.Name,
		CreatedByName:         creator.Name,
		Items:                 items,
	}

	if err := pi.GenerateID(); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].PurchaseInvoiceID = pi.ID
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

	if err := u.repo.Create(txCtx, pi); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toPIDetailResponse(pi), nil
}

func (u *purchaseInvoiceUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.PurchaseInvoiceDetailResponse, error) {
	pi, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pi == nil {
		return nil, ErrPurchaseInvoiceNotFound
	}

	return toPIDetailResponse(pi), nil
}

func (u *purchaseInvoiceUseCaseImpl) GetByInvoiceNumber(ctx context.Context, invoiceNum string) (*dto.PurchaseInvoiceDetailResponse, error) {
	pi, err := u.repo.FindByInvoiceNumber(ctx, invoiceNum)
	if err != nil {
		return nil, err
	}
	if pi == nil {
		return nil, ErrPurchaseInvoiceNotFound
	}

	return toPIDetailResponse(pi), nil
}

func (u *purchaseInvoiceUseCaseImpl) GetByStoreID(ctx context.Context, storeID string, status string) ([]dto.PurchaseInvoiceListResponse, error) {
	pis, err := u.repo.FindByStoreID(ctx, storeID, status)
	if err != nil {
		return nil, err
	}

	return toPIListResponses(pis), nil
}

func (u *purchaseInvoiceUseCaseImpl) GetPendingByStoreID(ctx context.Context, storeID string) ([]dto.PurchaseInvoiceListResponse, error) {
	pis, err := u.repo.FindPendingByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return toPIListResponses(pis), nil
}

func (u *purchaseInvoiceUseCaseImpl) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PurchaseInvoiceListResponse, *entity.Meta, error) {
	filter := &repository.QueryFilter{
		Page:          meta.Page,
		Limit:         meta.Limit,
		Search:        meta.Search,
		OrderBy:       meta.OrderColumn,
		OrderDir:      meta.OrderDir,
		SearchColumns: []string{"invoice_number", "supplier_invoice_number"},
		Conditions:    meta.Conditions,
	}

	data, resMeta, err := u.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return toPIListResponses(data), resMeta, nil
}

func (u *purchaseInvoiceUseCaseImpl) Update(ctx context.Context, userID string, id string, req dto.UpdatePurchaseInvoiceRequest) (*dto.PurchaseInvoiceDetailResponse, error) {
	pi, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pi == nil {
		return nil, ErrPurchaseInvoiceNotFound
	}

	if pi.Status != entity.PurchaseInvoiceStatusDraft {
		return nil, ErrPIInvalidStatus
	}

	invoiceDate := req.InvoiceDate
	if invoiceDate.IsZero() {
		invoiceDate = time.Now()
	}
	dueDate := invoiceDate.AddDate(0, 0, req.PaymentTermDays)

	subtotal := decimal.Zero
	taxAmount := decimal.Zero

	items := make([]entity.PurchaseInvoiceItem, len(req.Items))
	for i, item := range req.Items {
		lineSubtotal := item.QtyInvoiced.Mul(item.UnitPrice)

		// Simple discount calculation (can be improved to match Create logic if needed)
		d1 := lineSubtotal.Mul(item.Discount1Pct.Div(decimal.NewFromInt(100)))
		afterD1 := lineSubtotal.Sub(d1)
		d2 := afterD1.Mul(item.Discount2Pct.Div(decimal.NewFromInt(100)))
		afterD2 := afterD1.Sub(d2)
		d3 := afterD2.Mul(item.Discount3Pct.Div(decimal.NewFromInt(100)))
		
		totalDiscount := d1.Add(d2).Add(d3).Add(item.DiscountAmount)
		lineNet := lineSubtotal.Sub(totalDiscount)
		lineTax := lineNet.Mul(item.TaxPct.Div(decimal.NewFromInt(100)))

		subtotal = subtotal.Add(lineSubtotal)
		taxAmount = taxAmount.Add(lineTax)

		items[i] = entity.PurchaseInvoiceItem{
			PurchaseInvoiceID:   pi.ID,
			SeqNo:               i + 1,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			GoodsReceiptItemID:  item.GoodsReceiptItemID,
			ProductID:           item.ProductID,
			UOMID:               item.UOMID,
			QtyInvoiced:         item.QtyInvoiced,
			UnitPrice:           item.UnitPrice,
			Discount1Pct:        item.Discount1Pct,
			Discount2Pct:        item.Discount2Pct,
			Discount3Pct:        item.Discount3Pct,
			DiscountAmount:      item.DiscountAmount,
			TotalDiscountAmount: totalDiscount,
			TaxPct:              item.TaxPct,
			TaxAmount:           lineTax,
			Subtotal:            lineSubtotal,
		}
		if err := items[i].GenerateID(); err != nil {
			return nil, err
		}
	}

	grandTotal := subtotal.Add(taxAmount).Add(req.FreightAmount).Add(req.OtherCostAmount).Sub(req.DiscountAmount)

	pi.SupplierInvoiceNumber = req.SupplierInvoiceNumber
	pi.ReferenceNo = req.ReferenceNo
	pi.APAccountID = req.APAccountID
	pi.InvoiceDate = invoiceDate
	pi.ReceivedDate = req.ReceivedDate
	pi.DueDate = dueDate
	pi.PaymentTermDays = req.PaymentTermDays
	pi.PaymentMode = req.PaymentMode
	pi.Subtotal = subtotal
	pi.DiscountAmount = req.DiscountAmount
	pi.TaxAmount = taxAmount
	pi.FreightAmount = req.FreightAmount
	pi.OtherCostAmount = req.OtherCostAmount
	pi.GrandTotal = grandTotal
	pi.IsTaxInclusive = req.IsTaxInclusive
	pi.RemainingAmount = grandTotal
	pi.Notes = req.Notes
	pi.Items = items

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	// Delete old items
	if err := u.repo.DeleteItemsByPurchaseInvoiceID(txCtx, id); err != nil {
		return nil, err
	}

	// Update header and save new items
	if err := u.repo.Update(txCtx, pi); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toPIDetailResponse(pi), nil
}

func (u *purchaseInvoiceUseCaseImpl) Submit(ctx context.Context, userID string, id string) error {
	pi, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if pi == nil {
		return ErrPurchaseInvoiceNotFound
	}

	if pi.Status != entity.PurchaseInvoiceStatusDraft {
		return ErrPIInvalidStatus
	}

	pi.Status = entity.PurchaseInvoiceStatusSubmitted

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, pi); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseInvoiceUseCaseImpl) Approve(ctx context.Context, userID string, id string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	pi, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if pi == nil {
		return ErrPurchaseInvoiceNotFound
	}

	if pi.Status != entity.PurchaseInvoiceStatusSubmitted {
		return ErrPIInvalidStatus
	}

	approver, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || approver == nil {
		return errors.New("approver user not found")
	}

	pi.Status = entity.PurchaseInvoiceStatusVerified
	pi.VerifiedByID = &userUUID
	now := time.Now()
	pi.VerifiedAt = &now

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, pi); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseInvoiceUseCaseImpl) Verify(ctx context.Context, userID string, id string) error {
	return u.Approve(ctx, userID, id)
}

func (u *purchaseInvoiceUseCaseImpl) Post(ctx context.Context, userID string, id string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	pi, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if pi == nil {
		return ErrPurchaseInvoiceNotFound
	}

	if pi.Status != entity.PurchaseInvoiceStatusVerified {
		return ErrPIInvalidStatus
	}

	poster, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || poster == nil {
		return errors.New("poster user not found")
	}

	pi.Status = entity.PurchaseInvoiceStatusPosted
	pi.PostedByID = &userUUID
	now := time.Now()
	pi.PostedAt = &now

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, pi); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseInvoiceUseCaseImpl) Cancel(ctx context.Context, userID string, id string, reason string) error {
	pi, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if pi == nil {
		return ErrPurchaseInvoiceNotFound
	}

	if pi.Status == entity.PurchaseInvoiceStatusPosted || pi.Status == entity.PurchaseInvoiceStatusPaid {
		return ErrPIInvalidStatus
	}

	pi.Status = entity.PurchaseInvoiceStatusCancelled
	if reason != "" {
		pi.Notes = &reason
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

	if err := u.repo.Update(txCtx, pi); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toPIListResponses(pis []entity.PurchaseInvoice) []dto.PurchaseInvoiceListResponse {
	responses := make([]dto.PurchaseInvoiceListResponse, len(pis))
	for i, pi := range pis {
		responses[i] = dto.PurchaseInvoiceListResponse{
			ID:                    pi.ID,
			InvoiceNumber:         pi.InvoiceNumber,
			SupplierInvoiceNumber: pi.SupplierInvoiceNumber,
			PurchaseOrderID:       pi.PurchaseOrderID,
			SupplierID:            pi.SupplierID,
			SupplierName:          pi.Supplier.Name,
			SupplierCode:          pi.Supplier.Code,
			StoreID:               pi.StoreID,
			StoreName:             pi.Store.Name,
			WarehouseID:           pi.WarehouseID,
			WarehouseName:         pi.Warehouse.Name,
			InvoiceDate:           pi.InvoiceDate,
			DueDate:               pi.DueDate,
			GrandTotal:            pi.GrandTotal,
			PaidAmount:            pi.PaidAmount,
			RemainingAmount:       pi.RemainingAmount,
			Status:                pi.Status,
			CreatedByID:           pi.CreatedByID,
			CreatedByName:         pi.CreatedBy.Name,
			CreatedAt:             pi.CreatedAt,
			UpdatedAt:             pi.UpdatedAt,
		}
	}
	return responses
}

func toPIDetailResponse(pi *entity.PurchaseInvoice) *dto.PurchaseInvoiceDetailResponse {
	items := make([]dto.PurchaseInvoiceItemResponse, len(pi.Items))
	for i, item := range pi.Items {
		items[i] = dto.PurchaseInvoiceItemResponse{
			ID:                  item.ID,
			SeqNo:               item.SeqNo,
			PurchaseOrderItemID: item.PurchaseOrderItemID,
			GoodsReceiptItemID:  item.GoodsReceiptItemID,
			ProductID:           item.ProductID,
			ProductName:         item.Product.Name,
			ProductSKU:          item.Product.SKU,
			UOMID:               item.UOMID,
			UOMCode:             item.UOM.Name,
			QtyInvoiced:         item.QtyInvoiced,
			UnitPrice:           item.UnitPrice,
			Discount1Pct:        item.Discount1Pct,
			Discount2Pct:        item.Discount2Pct,
			Discount3Pct:        item.Discount3Pct,
			DiscountAmount:      item.DiscountAmount,
			TotalDiscountAmount: item.TotalDiscountAmount,
			TaxPct:              item.TaxPct,
			TaxAmount:           item.TaxAmount,
			LandedCostAmount:    item.LandedCostAmount,
			Subtotal:            item.Subtotal,
			NetUnitPrice:        item.NetUnitPrice,
			Notes:               item.Notes,
		}
	}

	return &dto.PurchaseInvoiceDetailResponse{
		ID:                    pi.ID,
		InvoiceNumber:         pi.InvoiceNumber,
		SupplierInvoiceNumber: pi.SupplierInvoiceNumber,
		ReferenceNo:           pi.ReferenceNo,
		PurchaseOrderID:       pi.PurchaseOrderID,
		PONumber:              pi.PurchaseOrder.PONumber,
		Supplier: dto.SupplierResponse{
			ID:   pi.Supplier.ID,
			Code: pi.Supplier.Code,
			Name: pi.Supplier.Name,
		},
		Store: dto.StoreResponse{
			ID:   pi.StoreID.String(),
			Code: pi.StoreCode,
			Name: pi.StoreName,
		},
		Warehouse: dto.WarehouseResponse{
			ID:   pi.WarehouseID,
			Name: pi.WarehouseName,
		},
		APAccountID:             pi.APAccountID,
		InvoiceDate:             pi.InvoiceDate,
		ReceivedDate:            pi.ReceivedDate,
		DueDate:                 pi.DueDate,
		ExpectedDelivery:        pi.ExpectedDelivery,
		PaymentTermDays:         pi.PaymentTermDays,
		PaymentMode:             pi.PaymentMode,
		Subtotal:                pi.Subtotal,
		DiscountAmount:          pi.DiscountAmount,
		TaxAmount:               pi.TaxAmount,
		FreightAmount:           pi.FreightAmount,
		OtherCostAmount:         pi.OtherCostAmount,
		GrandTotal:              pi.GrandTotal,
		IsTaxInclusive:          pi.IsTaxInclusive,
		PaidAmount:              pi.PaidAmount,
		RemainingAmount:         pi.RemainingAmount,
		Status:                  pi.Status,
		VerifiedByID:            pi.VerifiedByID,
		VerifiedAt:              pi.VerifiedAt,
		PostedByID:              pi.PostedByID,
		PostedAt:                pi.PostedAt,
		CreatedByID:             pi.CreatedByID,
		CreatedByName:           pi.CreatedByName,
		Notes:                   pi.Notes,
		SupplierNameSnapshot:    pi.SupplierName,
		SupplierCodeSnapshot:    pi.SupplierCode,
		SupplierAddressSnapshot: pi.SupplierAddress,
		StoreCodeSnapshot:       pi.StoreCode,
		StoreNameSnapshot:       pi.StoreName,
		StoreAddressSnapshot:    pi.StoreAddress,
		WarehouseNameSnapshot:   pi.WarehouseName,
		CreatedAt:               pi.CreatedAt,
		UpdatedAt:               pi.UpdatedAt,
		Items:                   items,
	}
}
