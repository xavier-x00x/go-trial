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

type ExpenseVoucherUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreateExpenseVoucherRequest) (*entity.ExpenseVoucher, error)
	Update(ctx context.Context, userID string, id string, req dto.UpdateExpenseVoucherRequest) error
	Post(ctx context.Context, userID string, id string) error
	Cancel(ctx context.Context, userID string, id string) error
	GetByID(ctx context.Context, id string) (*dto.ExpenseVoucherDetailResponse, error)
	GetAllWithPagination(ctx context.Context, req *dto.MetaRequest) ([]dto.ExpenseVoucherListResponse, *entity.Meta, error)
}

type expenseVoucherUseCaseImpl struct {
	repo               repository.ExpenseVoucherRepository
	coaRepo            repository.ChartOfAccountRepository
	userRepo           repository.UserRepository
	numberSequenceRepo repository.NumberSequenceRepository
	uow                uow.UnitOfWork
	db                 *gorm.DB
}

func NewExpenseVoucherUseCase(
	repo repository.ExpenseVoucherRepository,
	coaRepo repository.ChartOfAccountRepository,
	userRepo repository.UserRepository,
	numberSequenceRepo repository.NumberSequenceRepository,
	uow uow.UnitOfWork,
	db *gorm.DB,
) ExpenseVoucherUseCase {
	return &expenseVoucherUseCaseImpl{
		repo:               repo,
		coaRepo:            coaRepo,
		userRepo:           userRepo,
		numberSequenceRepo: numberSequenceRepo,
		uow:                uow,
		db:                 db,
	}
}

func (u *expenseVoucherUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreateExpenseVoucherRequest) (*entity.ExpenseVoucher, error) {
	userUUID, _ := uuid.Parse(userID)

	var grandTotal decimal.Decimal
	items := make([]entity.ExpenseVoucherItem, len(req.Items))
	for i, item := range req.Items {
		grandTotal = grandTotal.Add(item.Amount)
		items[i] = entity.ExpenseVoucherItem{
			SeqNo:            i + 1,
			Description:      item.Description,
			ExpenseAccountID: item.ExpenseAccountID,
			Amount:           item.Amount,
		}
	}

	ev := &entity.ExpenseVoucher{
		VoucherDate:     req.VoucherDate,
		VendorName:      req.VendorName,
		PaymentType:     req.PaymentType,
		CreditAccountID: req.CreditAccountID,
		GrandTotal:      grandTotal,
		Status:          entity.EVStatusDraft,
		Notes:           req.Notes,
		CreatedByID:     userUUID,
		Items:           items,
	}

	if err := ev.GenerateID(); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].ExpenseVoucherID = ev.ID
	}

	ev.VoucherNumber = "EV/DRAFT/" + time.Now().Format("20060102150405")

	if err := u.repo.Create(ctx, ev); err != nil {
		return nil, err
	}

	return ev, nil
}

func (u *expenseVoucherUseCaseImpl) Update(ctx context.Context, userID string, id string, req dto.UpdateExpenseVoucherRequest) error {
	return u.uow.Do(ctx, func(txCtx context.Context) error {
		ev, err := u.repo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		if ev.Status != entity.EVStatusDraft {
			return errors.New("only DRAFT vouchers can be updated")
		}

		if err := u.repo.DeleteItemsByVoucherID(txCtx, id); err != nil {
			return err
		}

		var grandTotal decimal.Decimal
		items := make([]entity.ExpenseVoucherItem, len(req.Items))
		for i, item := range req.Items {
			grandTotal = grandTotal.Add(item.Amount)
			items[i] = entity.ExpenseVoucherItem{
				ExpenseVoucherID: ev.ID,
				SeqNo:            i + 1,
				Description:      item.Description,
				ExpenseAccountID: item.ExpenseAccountID,
				Amount:           item.Amount,
			}
		}

		ev.VoucherDate = req.VoucherDate
		ev.VendorName = req.VendorName
		ev.PaymentType = req.PaymentType
		ev.CreditAccountID = req.CreditAccountID
		ev.GrandTotal = grandTotal
		ev.Notes = req.Notes
		ev.Items = items

		return u.repo.Update(txCtx, ev)
	})
}

func (u *expenseVoucherUseCaseImpl) Post(ctx context.Context, userID string, id string) error {
	return u.uow.Do(ctx, func(txCtx context.Context) error {
		ev, err := u.repo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		if ev.Status != entity.EVStatusDraft {
			return errors.New("only DRAFT vouchers can be posted")
		}

		userUUID, _ := uuid.Parse(userID)
		now := time.Now()

		period := ev.VoucherDate.Format("2006-01")
		nextNum, err := u.numberSequenceRepo.GetNextNumber(txCtx, "EV", period)
		if err != nil {
			return err
		}
		docNo := fmt.Sprintf("EV/%s/%04d", ev.VoucherDate.Format("2006/01"), nextNum)

		ev.VoucherNumber = docNo
		ev.Status = entity.EVStatusPosted
		ev.PostedByID = &userUUID
		ev.PostedAt = &now

		if err := u.repo.Update(txCtx, ev); err != nil {
			return err
		}

		journal := &entity.JournalEntry{
			EntryDate:          ev.VoucherDate,
			EntryNumber:        docNo, // Often journals use doc no as ref
			SourceDocumentType: entity.JournalSourceExpenseVoucher,
			SourceDocumentID:   &ev.ID,
			SourceDocumentNo:   &ev.VoucherNumber,
			Description:        fmt.Sprintf("Expense Voucher: %s - %s", ev.VoucherNumber, ev.VendorName),
			Period:             period,
			TotalDebit:         ev.GrandTotal,
			TotalCredit:        ev.GrandTotal,
			Status:             entity.JournalStatusPosted,
			PostedByID:         userUUID,
		}

		if err := journal.GenerateID(); err != nil {
			return err
		}

		lines := make([]entity.JournalEntryLine, 0)
		for i, item := range ev.Items {
			lines = append(lines, entity.JournalEntryLine{
				JournalEntryID: journal.ID,
				SeqNo:          i + 1,
				AccountID:      item.ExpenseAccountID,
				DebitAmount:    item.Amount,
				CreditAmount:   decimal.Zero,
			})
		}
		lines = append(lines, entity.JournalEntryLine{
			JournalEntryID: journal.ID,
			SeqNo:          len(ev.Items) + 1,
			AccountID:      ev.CreditAccountID,
			DebitAmount:    decimal.Zero,
			CreditAmount:   ev.GrandTotal,
		})
		journal.Lines = lines

		tx := uow.GetTx(txCtx, u.db)
		return tx.Create(journal).Error
	})
}

func (u *expenseVoucherUseCaseImpl) Cancel(ctx context.Context, userID string, id string) error {
	return u.uow.Do(ctx, func(txCtx context.Context) error {
		ev, err := u.repo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		if ev.Status == entity.EVStatusVoided {
			return errors.New("voucher is already voided")
		}

		if ev.Status == entity.EVStatusDraft {
			ev.Status = entity.EVStatusVoided
			return u.repo.Update(txCtx, ev)
		}

		userUUID, _ := uuid.Parse(userID)
		
		ev.Status = entity.EVStatusVoided
		if err := u.repo.Update(txCtx, ev); err != nil {
			return err
		}

		tx := uow.GetTx(txCtx, u.db)
		var originalJournal entity.JournalEntry
		if err := tx.Where("source_document_id = ?", ev.ID).Preload("Lines").First(&originalJournal).Error; err != nil {
			return err
		}

		reversal := &entity.JournalEntry{
			EntryDate:          time.Now(),
			EntryNumber:        "REV-" + originalJournal.EntryNumber,
			SourceDocumentType: entity.JournalSourceExpenseVoucher,
			SourceDocumentID:   &ev.ID,
			SourceDocumentNo:   &ev.VoucherNumber,
			Description:        "REVERSAL - " + originalJournal.Description,
			Period:             time.Now().Format("2006-01"),
			TotalDebit:         originalJournal.TotalCredit,
			TotalCredit:        originalJournal.TotalDebit,
			Status:             entity.JournalStatusPosted,
			ReversalOfID:       &originalJournal.ID,
			PostedByID:         userUUID,
		}

		if err := reversal.GenerateID(); err != nil {
			return err
		}

		revLines := make([]entity.JournalEntryLine, len(originalJournal.Lines))
		for i, line := range originalJournal.Lines {
			revLines[i] = entity.JournalEntryLine{
				JournalEntryID: reversal.ID,
				SeqNo:          i + 1,
				AccountID:      line.AccountID,
				DebitAmount:    line.CreditAmount,
				CreditAmount:   line.DebitAmount,
			}
		}
		reversal.Lines = revLines

		return tx.Create(reversal).Error
	})
}

func (u *expenseVoucherUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.ExpenseVoucherDetailResponse, error) {
	ev, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ExpenseVoucherItemResponse, len(ev.Items))
	for i, item := range ev.Items {
		items[i] = dto.ExpenseVoucherItemResponse{
			ID:                 item.ID,
			SeqNo:              item.SeqNo,
			Description:        item.Description,
			ExpenseAccountID:   item.ExpenseAccountID,
			ExpenseAccountName: item.ExpenseAccount.Name,
			Amount:             item.Amount,
		}
	}

	return &dto.ExpenseVoucherDetailResponse{
		ID:                ev.ID,
		VoucherNumber:     ev.VoucherNumber,
		VoucherDate:       ev.VoucherDate,
		VendorName:        ev.VendorName,
		PaymentType:       ev.PaymentType,
		CreditAccountID:   ev.CreditAccountID,
		CreditAccountName: ev.CreditAccount.Name,
		GrandTotal:        ev.GrandTotal,
		Status:            ev.Status,
		Notes:             ev.Notes,
		Items:             items,
	}, nil
}

func (u *expenseVoucherUseCaseImpl) GetAllWithPagination(ctx context.Context, req *dto.MetaRequest) ([]dto.ExpenseVoucherListResponse, *entity.Meta, error) {
	filter := &repository.QueryFilter{
		Page:     req.Page,
		Limit:    req.Limit,
		Search:   req.Search,
		OrderBy:  req.OrderColumn,
		OrderDir: req.OrderDir,
	}

	if filter.OrderBy == "" {
		filter.OrderBy = "created_at"
	}
	if filter.OrderDir == "" {
		filter.OrderDir = "DESC"
	}

	evs, meta, err := u.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	res := make([]dto.ExpenseVoucherListResponse, len(evs))
	for i, ev := range evs {
		res[i] = dto.ExpenseVoucherListResponse{
			ID:            ev.ID,
			VoucherNumber: ev.VoucherNumber,
			VoucherDate:   ev.VoucherDate,
			VendorName:    ev.VendorName,
			PaymentType:   ev.PaymentType,
			GrandTotal:    ev.GrandTotal,
			Status:        ev.Status,
		}
	}

	return res, meta, nil
}
