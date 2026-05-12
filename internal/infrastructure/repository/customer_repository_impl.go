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

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) domainRepo.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, c *entity.Customer) error {
	return uow.GetTx(ctx, r.db).Create(c).Error
}

func (r *customerRepository) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).Preload("ARAccount").Where("id = ?", id).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}
	return &customer, nil
}

func (r *customerRepository) FindByCode(ctx context.Context, code string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find customer by code: %w", err)
	}
	return &customer, nil
}

func (r *customerRepository) FindByPhone(ctx context.Context, phone string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).Where("phone_number = ? AND deleted_at IS NULL", phone).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find customer by phone: %w", err)
	}
	return &customer, nil
}

func (r *customerRepository) FindAll(ctx context.Context) ([]entity.Customer, error) {
	var customers []entity.Customer
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&customers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find customers: %w", err)
	}
	return customers, nil
}

func (r *customerRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Customer, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Customer{})
	return PaginateAndFilter[entity.Customer](r.db, baseQuery, filter)
}

func (r *customerRepository) Update(ctx context.Context, c *entity.Customer) error {
	return uow.GetTx(ctx, r.db).Save(c).Error
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	var customer entity.Customer
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&customer).Error; err != nil {
		return err
	}
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	return r.db.WithContext(ctx).Model(&customer).Updates(map[string]interface{}{
		"code":       customer.Code + suffix,
		"deleted_at": time.Now(),
	}).Error
}