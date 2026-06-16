package usecase

import (
	"context"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"

	"github.com/google/uuid"
)

type PriceListUsecase struct {
	repo repository.PriceListRepository
}

func NewPriceListUsecase(repo repository.PriceListRepository) *PriceListUsecase {
	return &PriceListUsecase{
		repo: repo,
	}
}

func (u *PriceListUsecase) Create(ctx context.Context, storeIDStr *string, req dto.CreatePriceListRequest) (*entity.PriceList, error) {
	var storeID *uuid.UUID
	if storeIDStr != nil && *storeIDStr != "" {
		parsed, err := uuid.Parse(*storeIDStr)
		if err == nil {
			storeID = &parsed
		}
	}

	pl := &entity.PriceList{
		Code:         req.Code,
		Name:         req.Name,
		CurrencyCode: req.CurrencyCode,
		StoreID:      storeID,
		IsActive:     true,
	}
	if err := pl.GenerateID(); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, pl); err != nil {
		return nil, err
	}
	return pl, nil
}

func (u *PriceListUsecase) Update(ctx context.Context, req dto.UpdatePriceListRequest, id string) (*entity.PriceList, error) {
	pl, err := u.repo.FindByID(ctx, id)
	if err != nil || pl == nil {
		return nil, err
	}
	if req.Code != nil {
		pl.Code = *req.Code
	}
	if req.Name != nil {
		pl.Name = *req.Name
	}
	if req.CurrencyCode != nil {
		pl.CurrencyCode = *req.CurrencyCode
	}
	if req.IsActive != nil {
		pl.IsActive = *req.IsActive
	}
	if err := u.repo.Update(ctx, pl); err != nil {
		return nil, err
	}
	return pl, nil
}

func (u *PriceListUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}