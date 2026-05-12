package usecase

import (
	"context"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
)

type ProductPriceUsecase struct {
	repo repository.ProductPriceRepository
}

func NewProductPriceUsecase(repo repository.ProductPriceRepository) *ProductPriceUsecase {
	return &ProductPriceUsecase{repo: repo}
}

func (u *ProductPriceUsecase) Create(ctx context.Context, req dto.CreateProductPriceRequest) (*entity.ProductPrice, error) {
	pp := &entity.ProductPrice{
		PriceListID: req.PriceListID,
		ProductID: req.ProductID,
		UOMID:     req.UOMID,
		MarkupPct: req.MarkupPct,
		SellPrice: req.SellPrice,
		DiscountPct: req.DiscountPct,
	}
	if err := pp.GenerateID(); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, pp); err != nil {
		return nil, err
	}
	return pp, nil
}

func (u *ProductPriceUsecase) GetByID(ctx context.Context, id string) (*entity.ProductPrice, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *ProductPriceUsecase) GetByProductID(ctx context.Context, productID string) ([]entity.ProductPrice, error) {
	return u.repo.FindByProductID(ctx, productID)
}

func (u *ProductPriceUsecase) GetByPriceListID(ctx context.Context, priceListID string) ([]entity.ProductPrice, error) {
	return u.repo.FindByPriceListID(ctx, priceListID)
}

func (u *ProductPriceUsecase) GetAll(ctx context.Context) ([]entity.ProductPrice, error) {
	return u.repo.FindAll(ctx)
}

func (u *ProductPriceUsecase) Update(ctx context.Context, req dto.UpdateProductPriceRequest, id string) (*entity.ProductPrice, error) {
	pp, err := u.repo.FindByID(ctx, id)
	if err != nil || pp == nil {
		return nil, err
	}
	if req.PriceListID != nil {
		pp.PriceListID = *req.PriceListID
	}
	if req.ProductID != nil {
		pp.ProductID = *req.ProductID
	}
	if req.UOMID != nil {
		pp.UOMID = *req.UOMID
	}
	if req.MarkupPct != nil {
		pp.MarkupPct = *req.MarkupPct
	}
	if req.SellPrice != nil {
		pp.SellPrice = *req.SellPrice
	}
	if req.DiscountPct != nil {
		pp.DiscountPct = *req.DiscountPct
	}
	if err := u.repo.Update(ctx, pp); err != nil {
		return nil, err
	}
	return pp, nil
}

func (u *ProductPriceUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}