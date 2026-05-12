package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) domainRepo.PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) Create(ctx context.Context, pm *entity.PaymentMethod) error {
	return uow.GetTx(ctx, r.db).Create(pm).Error
}

func (r *paymentMethodRepository) FindByID(ctx context.Context, id string) (*entity.PaymentMethod, error) {
	var pm entity.PaymentMethod
	err := r.db.WithContext(ctx).Preload("DepositAccount").Preload("ExpenseAccount").Where("id = ?", id).First(&pm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find payment method: %w", err)
	}
	return &pm, nil
}

func (r *paymentMethodRepository) FindByCode(ctx context.Context, code string) (*entity.PaymentMethod, error) {
	var pm entity.PaymentMethod
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&pm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find payment method by code: %w", err)
	}
	return &pm, nil
}

func (r *paymentMethodRepository) FindAll(ctx context.Context) ([]entity.PaymentMethod, error) {
	var pms []entity.PaymentMethod
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&pms).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find payment methods: %w", err)
	}
	return pms, nil
}

func (r *paymentMethodRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.PaymentMethod, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.PaymentMethod{})
	return PaginateAndFilter[entity.PaymentMethod](r.db, baseQuery, filter)
}

func (r *paymentMethodRepository) Update(ctx context.Context, pm *entity.PaymentMethod) error {
	return uow.GetTx(ctx, r.db).Save(pm).Error
}

func (r *paymentMethodRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?").Delete(&entity.PaymentMethod{}).Error
}