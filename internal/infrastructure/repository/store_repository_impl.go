package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) domainRepo.StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *entity.Store) error {
	return uow.GetTx(ctx, r.db).Create(store).Error
}

func (r *storeRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Store, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Store{})
	return PaginateAndFilter[entity.Store](r.db, baseQuery, filter)
}

func (r *storeRepository) FindByID(ctx context.Context, id string) (*entity.Store, error) {
	var store entity.Store
	err := uow.GetTx(ctx, r.db).Where("id = ?", id).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindByCode(ctx context.Context, code string) (*entity.Store, error) {
	var store entity.Store
	err := uow.GetTx(ctx, r.db).Where("code = ? AND deleted_at IS NULL", code).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindAll(ctx context.Context) ([]entity.Store, error) {
	var stores []entity.Store
	err := uow.GetTx(ctx, r.db).Where("deleted_at IS NULL").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) Update(ctx context.Context, store *entity.Store) error {
	return uow.GetTx(ctx, r.db).Save(store).Error
}

func (r *storeRepository) Delete(ctx context.Context, id string) error {
	var store entity.Store
	if err := uow.GetTx(ctx, r.db).Where("id = ?", id).First(&store).Error; err != nil {
		return err
	}
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	return uow.GetTx(ctx, r.db).Model(&store).Updates(map[string]interface{}{
		"code":       store.Code + suffix,
		"deleted_at": time.Now(),
	}).Error
}
