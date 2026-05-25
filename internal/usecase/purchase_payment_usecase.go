package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrPurchasePaymentNotFound = errors.New("purchase payment not found")
	ErrPPInvalidStatus         = errors.New("invalid purchase payment status transition")
	ErrPPAlreadyPosted         = errors.New("payment already posted")
	ErrPPAlreadyVoided         = errors.New("payment already voided")
)

type PurchasePaymentUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreatePurchasePaymentRequest) (*dto.PurchasePaymentDetailResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PurchasePaymentDetailResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PurchasePaymentListResponse, *entity.Meta, error)
	Post(ctx context.Context, userID string, id string, req dto.PostPurchasePaymentRequest) error
	Void(ctx context.Context, userID string, id string, reason string) error
}

type PurchasePaymentConfig struct {
	Repo                 repository.PurchasePaymentRepository
	PurchaseInvoiceRepo  repository.PurchaseInvoiceRepository
	PurchaseReturnRepo   repository.PurchaseReturnRepository
	SupplierRepo         repository.SupplierRepository
	UserRepo             repository.UserRepository
	ChartOfAccountRepo   repository.ChartOfAccountRepository
	MonthlyAPBalanceRepo repository.MonthlyAPBalanceRepository
	NumberSequenceRepo   repository.NumberSequenceRepository
	DB                   *gorm.DB
	Uow                  interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

type purchasePaymentUseCaseImpl struct {
	repo                 repository.PurchasePaymentRepository
	purchaseInvoiceRepo  repository.PurchaseInvoiceRepository
	purchaseReturnRepo   repository.PurchaseReturnRepository
	supplierRepo         repository.SupplierRepository
	userRepo             repository.UserRepository
	chartOfAccountRepo   repository.ChartOfAccountRepository
	monthlyAPBalanceRepo repository.MonthlyAPBalanceRepository
	numberSequenceRepo   repository.NumberSequenceRepository
	db                   *gorm.DB
	uow                  interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

func NewPurchasePaymentUseCase(cfg PurchasePaymentConfig) PurchasePaymentUseCase {
	return &purchasePaymentUseCaseImpl{
		repo:                 cfg.Repo,
		purchaseInvoiceRepo:  cfg.PurchaseInvoiceRepo,
		purchaseReturnRepo:   cfg.PurchaseReturnRepo,
		supplierRepo:         cfg.SupplierRepo,
		userRepo:             cfg.UserRepo,
		chartOfAccountRepo:   cfg.ChartOfAccountRepo,
		monthlyAPBalanceRepo: cfg.MonthlyAPBalanceRepo,
		numberSequenceRepo:   cfg.NumberSequenceRepo,
		db:                   cfg.DB,
		uow:                  cfg.Uow,
	}
}

func (u *purchasePaymentUseCaseImpl) generatePaymentNumber(date time.Time) (string, error) {
	prefix := "PAY"
	period := date.Format("0601")

	seqNum, err := u.numberSequenceRepo.GetNextNumber(context.Background(), prefix, period)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%05d", prefix, period, seqNum), nil
}

func (u *purchasePaymentUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreatePurchasePaymentRequest) (*dto.PurchasePaymentDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var fe FieldErrors

	if len(req.Items) == 0 {
		fe.Add("items", "items harus diisi")
		return nil, &fe
	}

	supplier, err := u.supplierRepo.FindByID(ctx, req.SupplierID.String())
	if err != nil || supplier == nil {
		fe.Add("supplier_id", "supplier tidak ditemukan")
		return nil, &fe
	}

	creator, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || creator == nil {
		fe.Add("created_by", "user tidak ditemukan")
		return nil, &fe
	}

	paymentDate := req.PaymentDate
	if paymentDate.IsZero() {
		paymentDate = time.Now()
	}

	paymentNum, err := u.generatePaymentNumber(paymentDate)
	if err != nil {
		return nil, err
	}

	var totalAmount decimal.Decimal
	var totalInvoice decimal.Decimal
	var totalReturn decimal.Decimal

	items := make([]entity.PurchasePaymentItem, len(req.Items))
	for i, item := range req.Items {
		var docAmount decimal.Decimal
		if item.PurchaseInvoiceID != nil {
			pi, err := u.purchaseInvoiceRepo.FindByID(ctx, item.PurchaseInvoiceID.String())
			if err != nil || pi == nil {
				return nil, fmt.Errorf("purchase invoice %s not found", item.PurchaseInvoiceID)
			}

			if pi.Status != entity.PurchaseInvoiceStatusPosted && pi.Status != entity.PurchaseInvoiceStatusPartiallyPaid {
				return nil, fmt.Errorf("invoice %s is not payable (status: %s)", pi.InvoiceNumber, pi.Status)
			}

			if item.PaidAmount.GreaterThan(pi.RemainingAmount) {
				return nil, fmt.Errorf("paid amount exceeds remaining amount for invoice %s", pi.InvoiceNumber)
			}
			docAmount = pi.GrandTotal
			totalInvoice = totalInvoice.Add(item.PaidAmount)
		} else if item.PurchaseReturnID != nil {
			pr, err := u.purchaseReturnRepo.FindByID(ctx, item.PurchaseReturnID.String())
			if err != nil || pr == nil {
				return nil, fmt.Errorf("purchase return %s not found", item.PurchaseReturnID)
			}

			if pr.Status != entity.PRStatusPosted {
				return nil, fmt.Errorf("return %s is not posted", pr.ReturnNumber)
			}

			if item.PaidAmount.Abs().GreaterThan(pr.RemainingAmount) {
				return nil, fmt.Errorf("offset amount exceeds remaining amount for return %s", pr.ReturnNumber)
			}
			docAmount = pr.GrandTotal
			totalReturn = totalReturn.Add(item.PaidAmount.Abs())
		} else {
			fe.Add("items", "purchase_invoice_id atau purchase_return_id harus diisi")
			return nil, &fe
		}

		totalAmount = totalAmount.Add(item.PaidAmount)

		items[i] = entity.PurchasePaymentItem{
			SeqNo:             i + 1,
			PurchaseInvoiceID: item.PurchaseInvoiceID,
			PurchaseReturnID:  item.PurchaseReturnID,
			DocumentAmount:    docAmount,
			PaidAmount:        item.PaidAmount,
		}
	}

	if totalReturn.GreaterThan(totalInvoice) {
		fe.Add("items", "total return tidak boleh melebihi total invoice")
		return nil, &fe
	}

	if req.PaymentMode == "GIRO" && req.GiroNumber == nil {
		fe.Add("giro_number", "giro_number wajib diisi jika payment_mode adalah GIRO")
		return nil, &fe
	}

	// Final Total Amount (Cash/Bank Out)
	totalAmount = totalInvoice.Sub(totalReturn).Add(req.AdminFeeAmount).Sub(req.DiscountAmount).Sub(req.WHTAmount)

	pp := &entity.PurchasePayment{
		PaymentNumber:     paymentNum,
		ReferenceNo:       req.ReferenceNo,
		SupplierID:        req.SupplierID,
		PaymentAccountID:  req.PaymentAccountID,
		APAccountID:       req.APAccountID,
		PaymentDate:       paymentDate,
		PaymentMode:       req.PaymentMode,
		GiroNumber:        req.GiroNumber,
		GiroDueDate:       req.GiroDueDate,
		TotalAmount:       totalAmount,
		AdminFeeAmount:    req.AdminFeeAmount,
		AdminFeeAccountID: req.AdminFeeAccountID,
		DiscountAmount:    req.DiscountAmount,
		DiscountAccountID: req.DiscountAccountID,
		WHTAmount:         req.WHTAmount,
		WHTAccountID:      req.WHTAccountID,
		Status:            entity.PPStatusDraft,
		CreatedByID:       userUUID,
		Notes:             req.Notes,
		Items:             items,
	}

	if err := pp.GenerateID(); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].PurchasePaymentID = pp.ID
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

	if err := u.repo.Create(txCtx, pp); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toPPDetailResponse(pp), nil
}

func (u *purchasePaymentUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.PurchasePaymentDetailResponse, error) {
	pp, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pp == nil {
		return nil, ErrPurchasePaymentNotFound
	}

	return toPPDetailResponse(pp), nil
}

func (u *purchasePaymentUseCaseImpl) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PurchasePaymentListResponse, *entity.Meta, error) {
	filter := &repository.QueryFilter{
		Page:          meta.Page,
		Limit:         meta.Limit,
		Search:        meta.Search,
		OrderBy:       meta.OrderColumn,
		OrderDir:      meta.OrderDir,
		SearchColumns: []string{"payment_number", "reference_no"},
		Conditions:    meta.Conditions,
	}

	data, resMeta, err := u.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return toPPListResponses(data), resMeta, nil
}

func (u *purchasePaymentUseCaseImpl) Post(ctx context.Context, userID string, id string, req dto.PostPurchasePaymentRequest) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	pp, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if pp == nil {
		return ErrPurchasePaymentNotFound
	}

	if pp.Status != entity.PPStatusDraft {
		return ErrPPInvalidStatus
	}

	poster, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || poster == nil {
		return errors.New("poster user not found")
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

	tx := uow.GetTx(txCtx, u.db)

	for _, item := range pp.Items {
		if item.PurchaseInvoiceID != nil {
			pi, err := u.purchaseInvoiceRepo.FindByID(txCtx, item.PurchaseInvoiceID.String())
			if err != nil || pi == nil {
				return fmt.Errorf("purchase invoice not found: %s", item.PurchaseInvoiceID)
			}

			pi.PaidAmount = pi.PaidAmount.Add(item.PaidAmount)
			pi.RemainingAmount = pi.RemainingAmount.Sub(item.PaidAmount)

			if pi.RemainingAmount.IsZero() || pi.RemainingAmount.LessThan(decimal.Zero) {
				pi.Status = entity.PurchaseInvoiceStatusPaid
			} else {
				pi.Status = entity.PurchaseInvoiceStatusPartiallyPaid
			}

			if err := u.purchaseInvoiceRepo.Update(txCtx, pi); err != nil {
				return err
			}
		} else if item.PurchaseReturnID != nil {
			pr, err := u.purchaseReturnRepo.FindByID(txCtx, item.PurchaseReturnID.String())
			if err != nil || pr == nil {
				return fmt.Errorf("purchase return not found: %s", item.PurchaseReturnID)
			}

			pr.RemainingAmount = pr.RemainingAmount.Sub(item.PaidAmount.Abs())
			if err := u.purchaseReturnRepo.Update(txCtx, pr); err != nil {
				return err
			}
		}

		periodMonth := pp.PaymentDate.Format("2006-01")
		apBalance, err := u.monthlyAPBalanceRepo.FindByPeriodSupplier(txCtx, periodMonth, pp.SupplierID.String())
		if err != nil {
			return err
		}

		if apBalance == nil {
			apBalance = &entity.MonthlyAPBalance{
				PeriodMonth:      periodMonth,
				SupplierID:       pp.SupplierID,
				BeginningBalance: decimal.Zero,
				TotalDebit:       decimal.Zero,
				TotalCredit:      decimal.Zero,
				EndingBalance:    decimal.Zero,
			}
			if err := u.monthlyAPBalanceRepo.Create(txCtx, apBalance); err != nil {
				return err
			}
		}

		apBalance.TotalDebit = apBalance.TotalDebit.Add(item.PaidAmount)
		apBalance.EndingBalance = apBalance.EndingBalance.Sub(item.PaidAmount)
		if err := u.monthlyAPBalanceRepo.Update(txCtx, apBalance); err != nil {
			return err
		}
	}

	journalEntry := u.createJournalEntry(txCtx, pp, poster, false)
	if err := tx.Create(journalEntry).Error; err != nil {
		return err
	}

	now := time.Now()
	pp.Status = entity.PPStatusPosted
	pp.PostedByID = &userUUID
	pp.PostedAt = &now
	if req.Notes != nil {
		pp.Notes = req.Notes
	}

	if err := u.repo.Update(txCtx, pp); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchasePaymentUseCaseImpl) Void(ctx context.Context, userID string, id string, reason string) error {
	pp, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if pp == nil {
		return ErrPurchasePaymentNotFound
	}

	if pp.Status != entity.PPStatusPosted {
		return ErrPPInvalidStatus
	}

	poster, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || poster == nil {
		return errors.New("poster user not found")
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

	tx := uow.GetTx(txCtx, u.db)

	for _, item := range pp.Items {
		if item.PurchaseInvoiceID != nil {
			pi, err := u.purchaseInvoiceRepo.FindByID(txCtx, item.PurchaseInvoiceID.String())
			if err != nil || pi == nil {
				return fmt.Errorf("purchase invoice not found: %s", item.PurchaseInvoiceID)
			}

			pi.PaidAmount = pi.PaidAmount.Sub(item.PaidAmount)
			pi.RemainingAmount = pi.RemainingAmount.Add(item.PaidAmount)

			if pi.RemainingAmount.GreaterThanOrEqual(pi.GrandTotal) {
				pi.Status = entity.PurchaseInvoiceStatusPosted
			} else {
				pi.Status = entity.PurchaseInvoiceStatusPartiallyPaid
			}

			if err := u.purchaseInvoiceRepo.Update(txCtx, pi); err != nil {
				return err
			}
		} else if item.PurchaseReturnID != nil {
			pr, err := u.purchaseReturnRepo.FindByID(txCtx, item.PurchaseReturnID.String())
			if err != nil || pr == nil {
				return fmt.Errorf("purchase return not found: %s", item.PurchaseReturnID)
			}

			pr.RemainingAmount = pr.RemainingAmount.Add(item.PaidAmount.Abs())
			if err := u.purchaseReturnRepo.Update(txCtx, pr); err != nil {
				return err
			}
		}

		periodMonth := pp.PaymentDate.Format("2006-01")
		apBalance, err := u.monthlyAPBalanceRepo.FindByPeriodSupplier(txCtx, periodMonth, pp.SupplierID.String())
		if err != nil {
			return err
		}

		if apBalance != nil {
			apBalance.TotalDebit = apBalance.TotalDebit.Sub(item.PaidAmount)
			apBalance.EndingBalance = apBalance.EndingBalance.Add(item.PaidAmount)
			if err := u.monthlyAPBalanceRepo.Update(txCtx, apBalance); err != nil {
				return err
			}
		}
	}

	reversalJournal := u.createJournalEntry(txCtx, pp, poster, true)
	if err := tx.Create(reversalJournal).Error; err != nil {
		return err
	}

	pp.Status = entity.PPStatusVoided
	pp.Notes = &reason

	if err := u.repo.Update(txCtx, pp); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchasePaymentUseCaseImpl) createJournalEntry(ctx context.Context, pp *entity.PurchasePayment, poster *entity.User, isReversal bool) *entity.JournalEntry {
	entryDate := pp.PaymentDate
	period := entryDate.Format("2006-01")
	description := fmt.Sprintf("Pembayaran ke %s", pp.Supplier.Name)

	if isReversal {
		description = fmt.Sprintf("VOID - %s", description)
	}

	seqNum, _ := u.numberSequenceRepo.GetNextNumber(context.Background(), "JE", period)
	entryNumber := fmt.Sprintf("JE/%s/%05d", period, seqNum)

	totalAP := decimal.Zero
	for _, item := range pp.Items {
		totalAP = totalAP.Add(item.PaidAmount)
	}

	posterID, _ := uuid.Parse(poster.ID)

	je := &entity.JournalEntry{
		EntryNumber:        entryNumber,
		SourceDocumentType: entity.JournalSourcePurchasePayment,
		SourceDocumentID:   &pp.ID,
		SourceDocumentNo:   &pp.PaymentNumber,
		EntryDate:          entryDate,
		Period:             period,
		Description:        description,
		Status:             entity.JournalStatusPosted,
		PostedByID:         posterID,
	}

	lines := []entity.JournalEntryLine{}
	seq := 1

	// 1. Debit: Accounts Payable (Total allocation)
	apAmount := totalAP
	if isReversal {
		apAmount = apAmount.Neg()
	}
	lines = append(lines, entity.JournalEntryLine{
		SeqNo:        seq,
		AccountID:    pp.APAccountID,
		DebitAmount:  apAmount,
		CreditAmount: decimal.Zero,
		Description:  &description,
	})
	seq++

	// 2. Debit: Admin Fee (if any)
	if pp.AdminFeeAmount.IsPositive() && pp.AdminFeeAccountID != nil {
		feeAmount := pp.AdminFeeAmount
		if isReversal {
			feeAmount = feeAmount.Neg()
		}
		lines = append(lines, entity.JournalEntryLine{
			SeqNo:        seq,
			AccountID:    *pp.AdminFeeAccountID,
			DebitAmount:  feeAmount,
			CreditAmount: decimal.Zero,
			Description:  &description,
		})
		seq++
	}

	// 3. Credit: Kas/Bank (Actual Cash Out)
	cashAmount := pp.TotalAmount
	if isReversal {
		cashAmount = cashAmount.Neg()
	}
	lines = append(lines, entity.JournalEntryLine{
		SeqNo:        seq,
		AccountID:    pp.PaymentAccountID,
		DebitAmount:  decimal.Zero,
		CreditAmount: cashAmount,
		Description:  &description,
	})
	seq++

	// 4. Credit: Discount (if any)
	if pp.DiscountAmount.IsPositive() && pp.DiscountAccountID != nil {
		discAmount := pp.DiscountAmount
		if isReversal {
			discAmount = discAmount.Neg()
		}
		lines = append(lines, entity.JournalEntryLine{
			SeqNo:        seq,
			AccountID:    *pp.DiscountAccountID,
			DebitAmount:  decimal.Zero,
			CreditAmount: discAmount,
			Description:  &description,
		})
		seq++
	}

	// 5. Credit: WHT/PPh (if any)
	if pp.WHTAmount.IsPositive() && pp.WHTAccountID != nil {
		whtAmt := pp.WHTAmount
		if isReversal {
			whtAmt = whtAmt.Neg()
		}
		lines = append(lines, entity.JournalEntryLine{
			SeqNo:        seq,
			AccountID:    *pp.WHTAccountID,
			DebitAmount:  decimal.Zero,
			CreditAmount: whtAmt,
			Description:  &description,
		})
		seq++
	}

	je.Lines = lines
	je.TotalDebit = totalAP.Add(pp.AdminFeeAmount)
	je.TotalCredit = pp.TotalAmount.Add(pp.DiscountAmount).Add(pp.WHTAmount)

	if isReversal {
		je.TotalDebit = je.TotalDebit.Neg()
		je.TotalCredit = je.TotalCredit.Neg()
	}

	je.GenerateID()

	for i := range je.Lines {
		je.Lines[i].JournalEntryID = je.ID
	}

	return je
}

func toPPListResponses(pps []entity.PurchasePayment) []dto.PurchasePaymentListResponse {
	responses := make([]dto.PurchasePaymentListResponse, len(pps))
	for i, pp := range pps {
		responses[i] = dto.PurchasePaymentListResponse{
			ID:            pp.ID,
			PaymentNumber: pp.PaymentNumber,
			SupplierID:    pp.SupplierID,
			SupplierName:  pp.Supplier.Name,
			PaymentDate:   pp.PaymentDate,
			PaymentMode:   pp.PaymentMode,
			TotalAmount:   pp.TotalAmount,
			Status:        pp.Status,
			CreatedAt:     pp.CreatedAt,
		}
	}
	return responses
}

func toPPDetailResponse(pp *entity.PurchasePayment) *dto.PurchasePaymentDetailResponse {
	items := make([]dto.PurchasePaymentItemResponse, len(pp.Items))
	for i, item := range pp.Items {
		invoiceNum := ""
		if item.PurchaseInvoiceID != nil {
			invoiceNum = item.PurchaseInvoice.InvoiceNumber
		}
		returnNum := ""
		if item.PurchaseReturnID != nil {
			returnNum = item.PurchaseReturn.ReturnNumber
		}
		items[i] = dto.PurchasePaymentItemResponse{
			ID:                item.ID,
			SeqNo:             item.SeqNo,
			PurchaseInvoiceID: item.PurchaseInvoiceID,
			InvoiceNumber:     invoiceNum,
			PurchaseReturnID:  item.PurchaseReturnID,
			ReturnNumber:      returnNum,
			DocumentAmount:    item.DocumentAmount,
			PaidAmount:        item.PaidAmount,
		}
	}

	resp := &dto.PurchasePaymentDetailResponse{
		ID:                pp.ID,
		PaymentNumber:     pp.PaymentNumber,
		ReferenceNo:       pp.ReferenceNo,
		SupplierID:        pp.SupplierID,
		SupplierName:      pp.Supplier.Name,
		PaymentAccountID:  pp.PaymentAccountID,
		APAccountID:       pp.APAccountID,
		PaymentDate:       pp.PaymentDate,
		PaymentMode:       pp.PaymentMode,
		GiroNumber:        pp.GiroNumber,
		GiroDueDate:       pp.GiroDueDate,
		TotalAmount:       pp.TotalAmount,
		AdminFeeAmount:    pp.AdminFeeAmount,
		AdminFeeAccountID: pp.AdminFeeAccountID,
		DiscountAmount:    pp.DiscountAmount,
		DiscountAccountID: pp.DiscountAccountID,
		WHTAmount:         pp.WHTAmount,
		WHTAccountID:      pp.WHTAccountID,
		Status:            pp.Status,
		CreatedByID:       pp.CreatedByID,
		PostedByID:        pp.PostedByID,
		PostedAt:          pp.PostedAt,
		Notes:             pp.Notes,
		CreatedAt:         pp.CreatedAt,
		UpdatedAt:         pp.UpdatedAt,
		Items:             items,
	}

	return resp
}
