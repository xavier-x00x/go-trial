package usecase

import (
	"context"
	"errors"
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"
)

var (
	ErrSupplierCategoryNotFound = errors.New("supplier category not found")
)

type SupplierCategoryUseCase interface {
	Create(ctx context.Context, req dto.CreateSupplierCategoryRequest) (*dto.SupplierCategoryResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateSupplierCategoryRequest) (*dto.SupplierCategoryResponse, error)
	Delete(ctx context.Context, id string) error
}

type supplierCategoryUseCase struct {
	categoryRepo repository.SupplierCategoryRepository
	uow          uow.UnitOfWork
}

func NewSupplierCategoryUseCase(categoryRepo repository.SupplierCategoryRepository, uow uow.UnitOfWork) SupplierCategoryUseCase {
	return &supplierCategoryUseCase{
		categoryRepo: categoryRepo,
		uow:          uow,
	}
}

func (u *supplierCategoryUseCase) Create(ctx context.Context, req dto.CreateSupplierCategoryRequest) (*dto.SupplierCategoryResponse, error) {
	var fe FieldErrors

	existing, err := u.categoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("name", "nama kategori pemasok sudah digunakan")
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	cat := &entity.SupplierCategory{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := cat.GenerateID(); err != nil {
		return nil, err
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.categoryRepo.Create(txCtx, cat); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toSupplierCategoryResponse(cat)
	return &resp, nil
}

func (u *supplierCategoryUseCase) Update(ctx context.Context, id string, req dto.UpdateSupplierCategoryRequest) (*dto.SupplierCategoryResponse, error) {
	cat, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrSupplierCategoryNotFound
	}

	var fe FieldErrors

	if req.Name != nil {
		if *req.Name != cat.Name {
			existing, err := u.categoryRepo.FindByName(ctx, *req.Name)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				fe.Add("name", "nama kategori pemasok sudah digunakan")
			}
			cat.Name = *req.Name
		}
	}
	if req.Description != nil {
		cat.Description = *req.Description
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.categoryRepo.Update(txCtx, cat); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toSupplierCategoryResponse(cat)
	return &resp, nil
}

func (u *supplierCategoryUseCase) Delete(ctx context.Context, id string) error {
	cat, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrSupplierCategoryNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.categoryRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toSupplierCategoryResponse(cat *entity.SupplierCategory) dto.SupplierCategoryResponse {
	return dto.SupplierCategoryResponse{
		ID:          cat.ID.String(),
		Name:        cat.Name,
		Description: cat.Description,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}
}
