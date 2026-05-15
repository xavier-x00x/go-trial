package repository

import (
	"context"

	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type expenseVoucherRepositoryImpl struct {
	db *gorm.DB
}

func NewExpenseVoucherRepository(db *gorm.DB) repository.ExpenseVoucherRepository {
	return &expenseVoucherRepositoryImpl{db: db}
}

func (r *expenseVoucherRepositoryImpl) Create(ctx context.Context, ev *entity.ExpenseVoucher) error {
	tx := uow.GetTx(ctx, r.db)
	return tx.Create(ev).Error
}

func (r *expenseVoucherRepositoryImpl) Update(ctx context.Context, ev *entity.ExpenseVoucher) error {
	tx := uow.GetTx(ctx, r.db)
	return tx.Save(ev).Error
}

func (r *expenseVoucherRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.ExpenseVoucher, error) {
	var ev entity.ExpenseVoucher
	err := r.db.WithContext(ctx).
		Preload("Items.ExpenseAccount").
		Preload("CreditAccount").
		Preload("CreatedBy").
		Preload("PostedBy").
		First(&ev, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *expenseVoucherRepositoryImpl) FindAllWithPagination(ctx context.Context, filter *repository.QueryFilter) ([]entity.ExpenseVoucher, *entity.Meta, error) {
	var evs []entity.ExpenseVoucher
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ExpenseVoucher{})

	if filter.Search != "" {
		s := "%" + filter.Search + "%"
		query = query.Where("voucher_number LIKE ? OR vendor_name LIKE ?", s, s)
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.Limit
	err := query.Order(filter.OrderBy + " " + filter.OrderDir).
		Limit(filter.Limit).
		Offset(offset).
		Find(&evs).Error

	if err != nil {
		return nil, nil, err
	}

	meta := &entity.Meta{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: int(total),
	}

	return evs, meta, nil
}

func (r *expenseVoucherRepositoryImpl) DeleteItemsByVoucherID(ctx context.Context, voucherID string) error {
	tx := uow.GetTx(ctx, r.db)
	return tx.Where("expense_voucher_id = ?", voucherID).Delete(&entity.ExpenseVoucherItem{}).Error
}
