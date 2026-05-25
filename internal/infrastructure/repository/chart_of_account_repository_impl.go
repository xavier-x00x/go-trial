package repository

import (
	"context"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type chartOfAccountRepository struct {
	db *gorm.DB
}

func NewChartOfAccountRepository(db *gorm.DB) domainRepo.ChartOfAccountRepository {
	return &chartOfAccountRepository{db: db}
}

func (r *chartOfAccountRepository) Create(ctx context.Context, coa *entity.ChartOfAccount) error {
	return uow.GetTx(ctx, r.db).Create(coa).Error
}

func (r *chartOfAccountRepository) FindByID(ctx context.Context, id string) (*entity.ChartOfAccount, error) {
	var coa entity.ChartOfAccount
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&coa).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find chart of account: %w", err)
	}
	return &coa, nil
}

func (r *chartOfAccountRepository) FindByCode(ctx context.Context, code string) (*entity.ChartOfAccount, error) {
	var coa entity.ChartOfAccount
	err := r.db.WithContext(ctx).Where("account_code = ? AND deleted_at IS NULL", code).First(&coa).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find chart of account by code: %w", err)
	}
	return &coa, nil
}

func (r *chartOfAccountRepository) FindAll(ctx context.Context) ([]entity.ChartOfAccount, error) {
	var coas []entity.ChartOfAccount
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&coas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find chart of accounts: %w", err)
	}
	return coas, nil
}

func (r *chartOfAccountRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.ChartOfAccount, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.ChartOfAccount{})
	return PaginateAndFilter[entity.ChartOfAccount](r.db, baseQuery, filter)
}

func (r *chartOfAccountRepository) FindByType(ctx context.Context, accountType string) ([]entity.ChartOfAccount, error) {
	var coas []entity.ChartOfAccount
	err := r.db.WithContext(ctx).Where("account_type = ? AND deleted_at IS NULL", accountType).Find(&coas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find chart of accounts by type: %w", err)
	}
	return coas, nil
}

func (r *chartOfAccountRepository) Update(ctx context.Context, coa *entity.ChartOfAccount) error {
	return uow.GetTx(ctx, r.db).Save(coa).Error
}

func (r *chartOfAccountRepository) FindByParentID(ctx context.Context, parentID *string) ([]entity.ChartOfAccount, error) {
	var coas []entity.ChartOfAccount
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
	if parentID == nil || *parentID == "" {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Order("account_code ASC").Find(&coas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find chart of accounts by parent: %w", err)
	}
	return coas, nil
}

func (r *chartOfAccountRepository) BulkCreate(ctx context.Context, coas []entity.ChartOfAccount) error {
	return uow.GetTx(ctx, r.db).Create(&coas).Error
}

func (r *chartOfAccountRepository) Delete(ctx context.Context, id string) error {
	var coa entity.ChartOfAccount
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&coa).Error; err != nil {
		return err
	}
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	return r.db.WithContext(ctx).Model(&coa).Updates(map[string]interface{}{
		"account_code": coa.AccountCode + suffix,
		"deleted_at":   time.Now(),
	}).Error
}