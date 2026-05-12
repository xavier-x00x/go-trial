package usecase

import (
	"context"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
)

type TaxUsecase struct {
	repo repository.TaxRepository
}

func NewTaxUsecase(repo repository.TaxRepository) *TaxUsecase {
	return &TaxUsecase{repo: repo}
}

func (u *TaxUsecase) Create(ctx context.Context, req dto.CreateTaxRequest) (*entity.Tax, error) {
	tax := &entity.Tax{
		Name:           req.Name,
		RatePercentage: req.GetRatePercentage(),
		TaxAccountID:   req.TaxAccountID,
	}
	if err := tax.GenerateID(); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, tax); err != nil {
		return nil, err
	}
	return tax, nil
}

func (u *TaxUsecase) GetByID(ctx context.Context, id string) (*entity.Tax, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *TaxUsecase) GetAll(ctx context.Context) ([]entity.Tax, error) {
	return u.repo.FindAll(ctx)
}

func (u *TaxUsecase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]entity.Tax, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "created_at"}
	searchColumns := []string{"id", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	return u.repo.FindAllWithPagination(ctx, filter)
}

func (u *TaxUsecase) Update(ctx context.Context, req dto.UpdateTaxRequest, id string) (*entity.Tax, error) {
	tax, err := u.repo.FindByID(ctx, id)
	if err != nil || tax == nil {
		return nil, err
	}
	if req.Name != nil {
		tax.Name = *req.Name
	}
	if req.RatePercentage != nil {
		tax.RatePercentage = req.GetRatePercentage()
	}
	if req.TaxAccountID != nil {
		tax.TaxAccountID = req.TaxAccountID
	}
	if err := u.repo.Update(ctx, tax); err != nil {
		return nil, err
	}
	return tax, nil
}

func (u *TaxUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}