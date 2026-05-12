package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"

	"gorm.io/gorm"
)

type monthlyAPBalanceRepository struct {
	db *gorm.DB
}

func NewMonthlyAPBalanceRepository(db *gorm.DB) domainRepo.MonthlyAPBalanceRepository {
	return &monthlyAPBalanceRepository{db: db}
}

func (r *monthlyAPBalanceRepository) FindByPeriodSupplier(ctx context.Context, periodMonth, supplierID string) (*entity.MonthlyAPBalance, error) {
	var mab entity.MonthlyAPBalance
	err := r.db.WithContext(ctx).Preload("Supplier").Where("period_month = ? AND supplier_id = ?", periodMonth, supplierID).First(&mab).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find monthly AP balance: %w", err)
	}
	return &mab, nil
}

func (r *monthlyAPBalanceRepository) FindByPeriodMonth(ctx context.Context, periodMonth string) ([]entity.MonthlyAPBalance, error) {
	var mabList []entity.MonthlyAPBalance
	err := r.db.WithContext(ctx).Preload("Supplier").Where("period_month = ?", periodMonth).Find(&mabList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find monthly AP balances: %w", err)
	}
	return mabList, nil
}

func (r *monthlyAPBalanceRepository) FindBySupplierID(ctx context.Context, supplierID string) ([]entity.MonthlyAPBalance, error) {
	var mabList []entity.MonthlyAPBalance
	err := r.db.WithContext(ctx).Where("supplier_id = ?", supplierID).Find(&mabList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find monthly AP balances: %w", err)
	}
	return mabList, nil
}

func (r *monthlyAPBalanceRepository) Create(ctx context.Context, mab *entity.MonthlyAPBalance) error {
	return r.db.Create(mab).Error
}

func (r *monthlyAPBalanceRepository) Update(ctx context.Context, mab *entity.MonthlyAPBalance) error {
	return r.db.Save(mab).Error
}