package repository

import (
	"context"

	"go-trial/internal/domain/entity"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type PurchaseOrderPlanningRepository interface {
	Create(ctx context.Context, p *entity.PurchaseOrderPlanning) error
	CreateBatch(ctx context.Context, plannings []entity.PurchaseOrderPlanning) error
	FindByID(ctx context.Context, id string) (*entity.PurchaseOrderPlanning, error)
	FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseOrderPlanning, error)
	FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseOrderPlanning, error)
	FindAll(ctx context.Context) ([]entity.PurchaseOrderPlanning, error)
	Update(ctx context.Context, p *entity.PurchaseOrderPlanning) error
	Delete(ctx context.Context, id string) error
	DeleteByDate(ctx context.Context, storeID string, date string) error
}

type purchaseOrderPlanningRepository struct {
	db *gorm.DB
}

func NewPurchaseOrderPlanningRepository(db *gorm.DB) PurchaseOrderPlanningRepository {
	return &purchaseOrderPlanningRepository{db: db}
}

func (r *purchaseOrderPlanningRepository) Create(ctx context.Context, p *entity.PurchaseOrderPlanning) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *purchaseOrderPlanningRepository) CreateBatch(ctx context.Context, plannings []entity.PurchaseOrderPlanning) error {
	if len(plannings) == 0 {
		return nil
	}
	db := uow.GetTx(ctx, r.db)
	return db.WithContext(ctx).Create(&plannings).Error
}

func (r *purchaseOrderPlanningRepository) FindByID(ctx context.Context, id string) (*entity.PurchaseOrderPlanning, error) {
	var p entity.PurchaseOrderPlanning
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *purchaseOrderPlanningRepository) FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseOrderPlanning, error) {
	var plans []entity.PurchaseOrderPlanning
	query := r.db.WithContext(ctx).Where("store_id = ?", storeID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("calculated_date DESC").Find(&plans).Error
	return plans, err
}

func (r *purchaseOrderPlanningRepository) FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseOrderPlanning, error) {
	var plans []entity.PurchaseOrderPlanning
	err := r.db.WithContext(ctx).
		Where("store_id = ? AND status = ?", storeID, entity.PlanningStatusPending).
		Order("calculated_date DESC").
		Find(&plans).Error
	return plans, err
}

func (r *purchaseOrderPlanningRepository) FindAll(ctx context.Context) ([]entity.PurchaseOrderPlanning, error) {
	var plans []entity.PurchaseOrderPlanning
	err := r.db.WithContext(ctx).Order("calculated_date DESC").Find(&plans).Error
	return plans, err
}

func (r *purchaseOrderPlanningRepository) Update(ctx context.Context, p *entity.PurchaseOrderPlanning) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *purchaseOrderPlanningRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.PurchaseOrderPlanning{}).Error
}

func (r *purchaseOrderPlanningRepository) DeleteByDate(ctx context.Context, storeID string, date string) error {
	return r.db.WithContext(ctx).
		Where("store_id = ? AND DATE(calculated_date) = ?", storeID, date).
		Delete(&entity.PurchaseOrderPlanning{}).Error
}