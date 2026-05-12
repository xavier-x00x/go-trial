package usecase

import (
	"context"
	"errors"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
)

var (
	ErrUOMNotFound   = errors.New("uom not found")
	ErrUOMCodeExists = errors.New("uom code already exists")
)

type UOMUseCase interface {
	Create(ctx context.Context, req dto.CreateUOMRequest) (*dto.UOMResponse, error)
	GetByID(ctx context.Context, id string) (*dto.UOMResponse, error)
	GetAll(ctx context.Context) ([]dto.UOMResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.UOMResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateUOMRequest) (*dto.UOMResponse, error)
}

type uomUseCase struct {
	uomRepo repository.UOMRepository
	uow     uow.UnitOfWork
}

func NewUOMUseCase(uomRepo repository.UOMRepository, uow uow.UnitOfWork) UOMUseCase {
	return &uomUseCase{
		uomRepo: uomRepo,
		uow:     uow,
	}
}

func (u *uomUseCase) Create(ctx context.Context, req dto.CreateUOMRequest) (*dto.UOMResponse, error) {
	existing, err := u.uomRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUOMCodeExists
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	uom := &entity.UOM{
		BaseModel: entity.BaseModel{ID: id},
		Code: req.Code,
		Name: req.Name,
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.uomRepo.Create(txCtx, uom); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toUOMResponse(uom)
	return &resp, nil
}

func (u *uomUseCase) GetByID(ctx context.Context, id string) (*dto.UOMResponse, error) {
	uom, err := u.uomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if uom == nil {
		return nil, ErrUOMNotFound
	}

	resp := toUOMResponse(uom)
	return &resp, nil
}

func (u *uomUseCase) GetAll(ctx context.Context) ([]dto.UOMResponse, error) {
	uoms, err := u.uomRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.UOMResponse
	for _, uom := range uoms {
		resp = append(resp, toUOMResponse(&uom))
	}
	return resp, nil
}

func (u *uomUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.UOMResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "created_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.uomRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.UOMResponse
	for _, uom := range data {
		resp = append(resp, toUOMResponse(&uom))
	}
	return resp, resMeta, nil
}

func (u *uomUseCase) Update(ctx context.Context, id string, req dto.UpdateUOMRequest) (*dto.UOMResponse, error) {
	uom, err := u.uomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if uom == nil {
		return nil, ErrUOMNotFound
	}

	if req.Code != nil && *req.Code != uom.Code {
		existing, err := u.uomRepo.FindByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrUOMCodeExists
		}
		uom.Code = *req.Code
	}
	if req.Name != nil {
		uom.Name = *req.Name
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.uomRepo.Update(txCtx, uom); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toUOMResponse(uom)
	return &resp, nil
}

func toUOMResponse(uom *entity.UOM) dto.UOMResponse {
	return dto.UOMResponse{
		ID:        uom.ID.String(),
		Code:     uom.Code,
		Name:     uom.Name,
		CreatedAt: uom.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: uom.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}