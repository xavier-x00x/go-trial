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
	"gorm.io/gorm"
)

var (
	ErrPurchaseReturnNotFound = errors.New("purchase return not found")
	ErrPRInvalidStatus        = errors.New("purchase return status is not DRAFT")
)

type PurchaseReturnUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreatePurchaseReturnRequest) (*dto.PurchaseReturnDetailResponse, error)
	Post(ctx context.Context, userID string, id string) error
	GetByID(ctx context.Context, id string) (*dto.PurchaseReturnDetailResponse, error)
	GetAllWithPagination(ctx context.Context, req *dto.MetaRequest) ([]dto.PurchaseReturnListResponse, *entity.Meta, error)
}

type PurchaseReturnConfig struct {
	PRRepo               domainRepo.PurchaseReturnRepository
	PIRepo               domainRepo.PurchaseInvoiceRepository
	UserRepo             domainRepo.UserRepository
	StoreRepo            domainRepo.StoreRepository
	InventoryStockRepo   domainRepo.InventoryStockRepository
	MonthlyAPBalanceRepo domainRepo.MonthlyAPBalanceRepository
	ChartOfAccountRepo   domainRepo.ChartOfAccountRepository
	NumberSequenceRepo   domainRepo.NumberSequenceRepository
	DB                   *gorm.DB
	Uow                  uow.UnitOfWork
}

type purchaseReturnUseCaseImpl struct {
	prRepo               domainRepo.PurchaseReturnRepository
	piRepo               domainRepo.PurchaseInvoiceRepository
	userRepo             domainRepo.UserRepository
	storeRepo            domainRepo.StoreRepository
	inventoryStockRepo   domainRepo.InventoryStockRepository
	monthlyAPBalanceRepo domainRepo.MonthlyAPBalanceRepository
	coaRepo              domainRepo.ChartOfAccountRepository
	numberSequenceRepo   domainRepo.NumberSequenceRepository
	db                   *gorm.DB
	uow                  uow.UnitOfWork
}

func NewPurchaseReturnUseCase(cfg PurchaseReturnConfig) PurchaseReturnUseCase {
	return &purchaseReturnUseCaseImpl{
		prRepo:               cfg.PRRepo,
		piRepo:               cfg.PIRepo,
		userRepo:             cfg.UserRepo,
		storeRepo:            cfg.StoreRepo,
		inventoryStockRepo:   cfg.InventoryStockRepo,
		monthlyAPBalanceRepo: cfg.MonthlyAPBalanceRepo,
		coaRepo:              cfg.ChartOfAccountRepo,
		numberSequenceRepo:   cfg.NumberSequenceRepo,
		db:                   cfg.DB,
		uow:                  cfg.Uow,
	}
}

func (u *purchaseReturnUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreatePurchaseReturnRequest) (*dto.PurchaseReturnDetailResponse, error) {
	var fe FieldErrors

	pi, err := u.piRepo.FindByID(ctx, req.PurchaseInvoiceID.String())
	if err != nil || pi == nil {
		fe.Add("purchase_invoice_id", "purchase invoice tidak ditemukan")
		return nil, &fe
	}

	if pi.Status != entity.PurchaseInvoiceStatusPosted {
		fe.Add("purchase_invoice_id", "hanya invoice yang sudah diposting yang dapat diretur")
		return nil, &fe
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		fe.Add("created_by", "user tidak ditemukan")
		return nil, &fe
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
			fe.Add("purchase_invoice_item_id", "item invoice tidak ditemukan")
			return nil, &fe
		}

		if item.QtyReturn.GreaterThan(piItem.QtyInvoiced) {
			fe.Add("qty_return", "jumlah retur melebihi jumlah invoice")
			return nil, &fe
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

	// Load Invoice for accounts
	pi, err := u.piRepo.FindByID(ctx, pr.PurchaseInvoiceID.String())
	if err != nil || pi == nil {
		return errors.New("purchase invoice not found")
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

	// Update Monthly AP Balance
	periodMonth := pr.ReturnDate.Format("2006-01")
	apBalance, err := u.monthlyAPBalanceRepo.FindByPeriodSupplier(txCtx, periodMonth, pr.SupplierID.String())
	if err != nil {
		return err
	}

	if apBalance == nil {
		apBalance = &entity.MonthlyAPBalance{
			PeriodMonth:      periodMonth,
			SupplierID:       pr.SupplierID,
			BeginningBalance: decimal.Zero,
			TotalDebit:       decimal.Zero,
			TotalCredit:      decimal.Zero,
			EndingBalance:    decimal.Zero,
		}
		if err := u.monthlyAPBalanceRepo.Create(txCtx, apBalance); err != nil {
			return err
		}
	}

	// Purchase Return decreases the payable balance, so we subtract from TotalCredit or add to TotalDebit.
	// In this system, we'll treat it as a reduction of Credit (Invoice).
	apBalance.TotalCredit = apBalance.TotalCredit.Sub(pr.GrandTotal)
	apBalance.EndingBalance = apBalance.EndingBalance.Sub(pr.GrandTotal)
	if err := u.monthlyAPBalanceRepo.Update(txCtx, apBalance); err != nil {
		return err
	}

	// Create Journal Entry
	if err := u.createJournalEntry(txCtx, pr, pi, userID); err != nil {
		return err
	}

	if err := u.prRepo.Update(txCtx, pr); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseReturnUseCaseImpl) createJournalEntry(ctx context.Context, pr *entity.PurchaseReturn, pi *entity.PurchaseInvoice, userID string) error {
	tx := uow.GetTx(ctx, u.db)
	entryDate := pr.ReturnDate
	period := entryDate.Format("2006-01")
	description := fmt.Sprintf("Retur Pembelian %s - %s (Ref: %s)", pr.ReturnNumber, pi.SupplierName, pi.InvoiceNumber)

	seqNum, _ := u.numberSequenceRepo.GetNextNumber(ctx, "JE", period)
	entryNumber := fmt.Sprintf("JE/%s/%05d", period, seqNum)

	userUUID, _ := uuid.Parse(userID)

	je := &entity.JournalEntry{
		EntryNumber:        entryNumber,
		SourceDocumentType: entity.JournalSourcePurchaseReturn,
		SourceDocumentID:   &pr.ID,
		SourceDocumentNo:   &pr.ReturnNumber,
		EntryDate:           entryDate,
		Period:             period,
		TotalDebit:         pr.GrandTotal,
		TotalCredit:        pr.GrandTotal,
		Description:        description,
		Status:             entity.JournalStatusPosted,
		PostedByID:         userUUID,
		Lines: []entity.JournalEntryLine{
			{
				SeqNo:        1,
				AccountID:    pi.APAccountID,
				DebitAmount:  pr.GrandTotal, // Reduce Payable
				CreditAmount: decimal.Zero,
				Description:  &description,
			},
			{
				SeqNo:        2,
				AccountID:    pi.InventoryAccountID,
				DebitAmount:  decimal.Zero,
				CreditAmount: pr.GrandTotal, // Reduce Inventory
				Description:  &description,
			},
		},
	}

	if err := je.GenerateID(); err != nil {
		return err
	}

	for i := range je.Lines {
		je.Lines[i].JournalEntryID = je.ID
	}

	return tx.Create(je).Error
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
