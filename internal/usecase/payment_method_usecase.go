package usecase

import (
	"context"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
)

type PaymentMethodUsecase struct {
	repo repository.PaymentMethodRepository
}

func NewPaymentMethodUsecase(repo repository.PaymentMethodRepository) *PaymentMethodUsecase {
	return &PaymentMethodUsecase{repo: repo}
}

type PaymentMethodUsecaseInterface interface {
	Create(ctx context.Context, req dto.CreatePaymentMethodRequest) (*entity.PaymentMethod, error)
	GetAll(ctx context.Context) ([]entity.PaymentMethod, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]entity.PaymentMethod, *entity.Meta, error)
	GetByID(ctx context.Context, id string) (*entity.PaymentMethod, error)
	Update(ctx context.Context, req dto.UpdatePaymentMethodRequest, id string) (*entity.PaymentMethod, error)
	Delete(ctx context.Context, id string) error
}

func (u *PaymentMethodUsecase) Create(ctx context.Context, req dto.CreatePaymentMethodRequest) (*entity.PaymentMethod, error) {
	pm := &entity.PaymentMethod{
		Code:             req.Code,
		Name:             req.Name,
		MdrPercentage:    req.GetMdrPercentage(),
		DepositAccountID: req.DepositAccountID,
		ExpenseAccountID: req.ExpenseAccountID,
		IsActive:         true,
	}
	if err := pm.GenerateID(); err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (u *PaymentMethodUsecase) GetByID(ctx context.Context, id string) (*entity.PaymentMethod, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *PaymentMethodUsecase) GetAll(ctx context.Context) ([]entity.PaymentMethod, error) {
	return u.repo.FindAll(ctx)
}

func (u *PaymentMethodUsecase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]entity.PaymentMethod, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "created_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	return u.repo.FindAllWithPagination(ctx, filter)
}

func (u *PaymentMethodUsecase) Update(ctx context.Context, req dto.UpdatePaymentMethodRequest, id string) (*entity.PaymentMethod, error) {
	pm, err := u.repo.FindByID(ctx, id)
	if err != nil || pm == nil {
		return nil, err
	}
	if req.Name != nil {
		pm.Name = *req.Name
	}
	if req.MdrPercentage != nil {
		pm.MdrPercentage = req.GetMdrPercentage()
	}
	if req.DepositAccountID != nil {
		pm.DepositAccountID = req.DepositAccountID
	}
	if req.ExpenseAccountID != nil {
		pm.ExpenseAccountID = req.ExpenseAccountID
	}
	if req.IsActive != nil {
		pm.IsActive = *req.IsActive
	}
	if err := u.repo.Update(ctx, pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (u *PaymentMethodUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}