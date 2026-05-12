package usecase

import (
	"context"
	"errors"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
)

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrWarehouseCodeExists = errors.New("warehouse code already exists")
)

type WarehouseUsecase struct {
	repo repository.WarehouseRepository
}

func NewWarehouseUsecase(repo repository.WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{repo: repo}
}

type WarehouseUsecaseInterface interface {
	Create(ctx context.Context, req dto.CreateWarehouseRequest) (*entity.Warehouse, error)
	GetByID(ctx context.Context, id string) (*entity.Warehouse, error)
	GetByStoreID(ctx context.Context, storeID string) ([]entity.Warehouse, error)
	GetAll(ctx context.Context) ([]entity.Warehouse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]entity.Warehouse, *entity.Meta, error)
	Update(ctx context.Context, req dto.UpdateWarehouseRequest, id string) (*entity.Warehouse, error)
	Delete(ctx context.Context, id string) error
}

func (u *WarehouseUsecase) Create(ctx context.Context, req dto.CreateWarehouseRequest) (*entity.Warehouse, error) {
	warehouse := &entity.Warehouse{
		StoreID:  req.StoreID,
		Code:     req.Code,
		Name:     req.Name,
		IsActive: true,
	}
	if err := warehouse.GenerateID(); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, warehouse); err != nil {
		return nil, err
	}
	return warehouse, nil
}

func (u *WarehouseUsecase) GetByID(ctx context.Context, id string) (*entity.Warehouse, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *WarehouseUsecase) GetByStoreID(ctx context.Context, storeID string) ([]entity.Warehouse, error) {
	return u.repo.FindByStoreID(ctx, storeID)
}

func (u *WarehouseUsecase) GetAll(ctx context.Context) ([]entity.Warehouse, error) {
	return u.repo.FindAll(ctx)
}

func (u *WarehouseUsecase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]entity.Warehouse, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	return u.repo.FindAllWithPagination(ctx, filter)
}

func (u *WarehouseUsecase) Update(ctx context.Context, req dto.UpdateWarehouseRequest, id string) (*entity.Warehouse, error) {
	warehouse, err := u.repo.FindByID(ctx, id)
	if err != nil || warehouse == nil {
		return nil, err
	}
	if req.StoreID != nil {
		warehouse.StoreID = *req.StoreID
	}
	if req.Code != nil {
		warehouse.Code = *req.Code
	}
	if req.Name != nil {
		warehouse.Name = *req.Name
	}
	if req.IsActive != nil {
		warehouse.IsActive = *req.IsActive
	}
	if err := u.repo.Update(ctx, warehouse); err != nil {
		return nil, err
	}
	return warehouse, nil
}

func (u *WarehouseUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}