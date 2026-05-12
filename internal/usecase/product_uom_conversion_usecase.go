package usecase

import (
	"context"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
)

type ProductUOMConversionUsecase struct {
	repo repository.ProductUOMConversionRepository
}

func NewProductUOMConversionUsecase(repo repository.ProductUOMConversionRepository) *ProductUOMConversionUsecase {
	return &ProductUOMConversionUsecase{repo: repo}
}

func (u *ProductUOMConversionUsecase) Create(ctx context.Context, req dto.CreateProductUOMConversionRequest) (*entity.ProductUOMConversion, error) {
	puc := &entity.ProductUOMConversion{
		ProductID:      req.ProductID,
		UOMID:          req.UOMID,
		ConversionRate: req.ConversionRate,
		Barcode:        req.Barcode,
		Length:        req.Length,
		Width:         req.Width,
		Height:        req.Height,
		Weight:        req.Weight,
		IsStackable:   req.IsStackable,
		MaxStackLayer: req.MaxStackLayer,
	}
	if err := puc.GenerateID(); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, puc); err != nil {
		return nil, err
	}
	return puc, nil
}

func (u *ProductUOMConversionUsecase) GetByID(ctx context.Context, id string) (*entity.ProductUOMConversion, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *ProductUOMConversionUsecase) GetByProductID(ctx context.Context, productID string) ([]entity.ProductUOMConversion, error) {
	return u.repo.FindByProductID(ctx, productID)
}

func (u *ProductUOMConversionUsecase) GetAll(ctx context.Context) ([]entity.ProductUOMConversion, error) {
	return u.repo.FindAll(ctx)
}

func (u *ProductUOMConversionUsecase) Update(ctx context.Context, req dto.UpdateProductUOMConversionRequest, id string) (*entity.ProductUOMConversion, error) {
	puc, err := u.repo.FindByID(ctx, id)
	if err != nil || puc == nil {
		return nil, err
	}
	if req.UOMID != nil {
		puc.UOMID = *req.UOMID
	}
	if req.ConversionRate != nil {
		puc.ConversionRate = *req.ConversionRate
	}
	if req.Barcode != nil {
		puc.Barcode = req.Barcode
	}
	if req.Length != nil {
		puc.Length = *req.Length
	}
	if req.Width != nil {
		puc.Width = *req.Width
	}
	if req.Height != nil {
		puc.Height = *req.Height
	}
	if req.Weight != nil {
		puc.Weight = *req.Weight
	}
	if req.IsStackable != nil {
		puc.IsStackable = *req.IsStackable
	}
	if req.MaxStackLayer != nil {
		puc.MaxStackLayer = *req.MaxStackLayer
	}
	if err := u.repo.Update(ctx, puc); err != nil {
		return nil, err
	}
	return puc, nil
}

func (u *ProductUOMConversionUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}