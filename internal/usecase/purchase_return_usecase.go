package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrPurchaseReturnNotFound = errors.New("purchase return not found")
	ErrPRInvalidStatus        = errors.New("purchase return status is not DRAFT")
	ErrPRQtyExceed            = errors.New("return quantity exceeds invoice quantity")
)

type PurchaseReturnUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreatePurchaseReturnRequest) (*dto.PurchaseReturnDetailResponse, error)
	Post(ctx context.Context, userID string, id string) error
	GetByID(ctx context.Context, id string) (*dto.PurchaseReturnDetailResponse, error)
	GetAllWithPagination(ctx context.Context, req *dto.MetaRequest) ([]dto.PurchaseReturnListResponse, *entity.Meta, error)
}

type purchaseReturnUseCaseImpl struct {
	prRepo             domainRepo.PurchaseReturnRepository
	piRepo             domainRepo.PurchaseInvoiceRepository
	userRepo           domainRepo.UserRepository
	storeRepo          domainRepo.StoreRepository
	inventoryStockRepo domainRepo.InventoryStockRepository
	uow                uow.UnitOfWork
}

func NewPurchaseReturnUseCase(
	prRepo domainRepo.PurchaseReturnRepository,
	piRepo domainRepo.PurchaseInvoiceRepository,
	userRepo domainRepo.UserRepository,
	storeRepo domainRepo.StoreRepository,
	inventoryStockRepo domainRepo.InventoryStockRepository,
	uow uow.UnitOfWork,
) PurchaseReturnUseCase {
	return &purchaseReturnUseCaseImpl{
		prRepo:             prRepo,
		piRepo:             piRepo,
		userRepo:           userRepo,
		storeRepo:          storeRepo,
		inventoryStockRepo: inventoryStockRepo,
		uow:                uow,
	}
}

func (u *purchaseReturnUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreatePurchaseReturnRequest) (*dto.PurchaseReturnDetailResponse, error) {
	pi, err := u.piRepo.FindByID(ctx, req.PurchaseInvoiceID.String())
	if err != nil || pi == nil {
		return nil, errors.New("purchase invoice not found")
	}

	if pi.Status != entity.PurchaseInvoiceStatusPosted {
		return nil, errors.New("only posted invoices can be returned")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	piItemsMap := make(map[string]entity.PurchaseInvoiceItem)
	for _, item := range pi.Items {
		piItemsMap[item.ID.String()] = item
	}

	subtotal := decimal.Zero
	taxAmount := decimal.Zero
	items := make([]entity.PurchaseReturnItem, len(req.Items))

	for i, item := range req.Items {
		piItem, exists := piItemsMap[item.PurchaseInvoiceItemID.String()]
		if !exists {
			return nil, errors.New("invoice item not found")
		}

		if item.QtyReturn.GreaterThan(piItem.QtyInvoiced) {
			return nil, ErrPRQtyExceed
		}

		lineSubtotal := item.QtyReturn.Mul(piItem.UnitPrice)
		lineTax := lineSubtotal.Mul(piItem.TaxPct.Div(decimal.NewFromInt(100)))

		subtotal = subtotal.Add(lineSubtotal)
		taxAmount = taxAmount.Add(lineTax)

		items[i] = entity.PurchaseReturnItem{
			SeqNo:               i + 1,
			PurchaseInvoiceItemID: item.PurchaseInvoiceItemID,
			ProductID:           item.ProductID,
			UOMID:               item.UOMID,
			QtyReturn:           item.QtyReturn,
			UnitPrice:           piItem.UnitPrice,
			Discount1Pct:        piItem.Discount1Pct,
			Discount2Pct:        piItem.Discount2Pct,
			Discount3Pct:        piItem.Discount3Pct,
			DiscountAmount:      piItem.DiscountAmount,
			TaxPct:              piItem.TaxPct,
			TaxAmount:           lineTax,
			Subtotal:            lineSubtotal,
		}
		if err := items[i].GenerateID(); err != nil {
			return nil, err
		}
	}

	grandTotal := subtotal.Add(taxAmount)
	returnNum := fmt.Sprintf("PR/%s/%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	pr := &entity.PurchaseReturn{
		ReturnNumber:      returnNum,
		ReturnDate:        req.ReturnDate,
		PurchaseInvoiceID: req.PurchaseInvoiceID,
		SupplierID:        pi.SupplierID,
		StoreID:           pi.StoreID,
		WarehouseID:       pi.WarehouseID,
		Subtotal:          subtotal,
		TaxAmount:         taxAmount,
		GrandTotal:        grandTotal,
		RemainingAmount:   grandTotal,
		Status:            entity.PRStatusDraft,
		Notes:             req.Notes,
		CreatedByID:       userUUID,
		Items:             items,
	}

	if err := pr.GenerateID(); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].PurchaseReturnID = pr.ID
	}

	if err := u.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return u.GetByID(ctx, pr.ID.String())
}

func (u *purchaseReturnUseCaseImpl) Post(ctx context.Context, userID string, id string) error {
	pr, err := u.prRepo.FindByID(ctx, id)
	if err != nil || pr == nil {
		return ErrPurchaseReturnNotFound
	}

	if pr.Status != entity.PRStatusDraft {
		return ErrPRInvalidStatus
	}

	userUUID, _ := uuid.Parse(userID)
	now := time.Now()
	pr.Status = entity.PRStatusPosted
	pr.PostedByID = &userUUID
	pr.PostedAt = &now

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	for _, item := range pr.Items {
		// Reduce Stock
		stock, err := u.inventoryStockRepo.FindByWarehouseAndProduct(txCtx, pr.WarehouseID.String(), item.ProductID.String())
		if err != nil {
			return err
		}
		if stock == nil {
			return errors.New("stock not found for product " + item.ProductID.String())
		}

		stock.Quantity = stock.Quantity.Sub(item.QtyReturn)
		if err := u.inventoryStockRepo.Update(txCtx, stock); err != nil {
			return err
		}
	}

	if err := u.prRepo.Update(txCtx, pr); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseReturnUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.PurchaseReturnDetailResponse, error) {
	pr, err := u.prRepo.FindByID(ctx, id)
	if err != nil || pr == nil {
		return nil, ErrPurchaseReturnNotFound
	}
	return toPRDetailResponse(pr), nil
}

func (u *purchaseReturnUseCaseImpl) GetAllWithPagination(ctx context.Context, req *dto.MetaRequest) ([]dto.PurchaseReturnListResponse, *entity.Meta, error) {
	filter := &domainRepo.QueryFilter{
		Page:          req.Page,
		Limit:         req.Limit,
		OrderBy:       req.OrderColumn,
		OrderDir:      req.OrderDir,
		Search:        req.Search,
		SearchColumns: []string{"return_number"},
		Conditions:    req.Conditions,
	}

	data, meta, err := u.prRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	res := make([]dto.PurchaseReturnListResponse, len(data))
	for i, d := range data {
		res[i] = dto.PurchaseReturnListResponse{
			ID:            d.ID,
			ReturnNumber:  d.ReturnNumber,
			ReturnDate:    d.ReturnDate,
			InvoiceNumber: d.PurchaseInvoice.InvoiceNumber,
			SupplierName:  d.Supplier.Name,
			GrandTotal:    d.GrandTotal,
			Status:        d.Status,
		}
	}

	return res, meta, nil
}

func toPRDetailResponse(pr *entity.PurchaseReturn) *dto.PurchaseReturnDetailResponse {
	items := make([]dto.PurchaseReturnItemResponse, len(pr.Items))
	for i, item := range pr.Items {
		items[i] = dto.PurchaseReturnItemResponse{
			ID:             item.ID,
			SeqNo:          item.SeqNo,
			ProductID:      item.ProductID,
			ProductName:    item.Product.Name,
			UOMCode:        item.UOM.Code,
			QtyReturn:      item.QtyReturn,
			UnitPrice:      item.UnitPrice,
			DiscountAmount: item.DiscountAmount,
			TaxAmount:      item.TaxAmount,
			Subtotal:       item.Subtotal,
		}
	}

	return &dto.PurchaseReturnDetailResponse{
		ID:                pr.ID,
		ReturnNumber:      pr.ReturnNumber,
		ReturnDate:        pr.ReturnDate,
		PurchaseInvoiceID: pr.PurchaseInvoiceID,
		InvoiceNumber:     pr.PurchaseInvoice.InvoiceNumber,
		SupplierName:      pr.Supplier.Name,
		StoreName:         pr.Store.Name,
		WarehouseName:     pr.Warehouse.Name,
		Subtotal:          pr.Subtotal,
		DiscountAmount:    pr.DiscountAmount,
		TaxAmount:         pr.TaxAmount,
		GrandTotal:        pr.GrandTotal,
		RemainingAmount:   pr.RemainingAmount,
		Status:            pr.Status,
		Notes:             pr.Notes,
		Items:             items,
	}
}
