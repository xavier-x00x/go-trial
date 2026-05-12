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
	ErrCOANotFound   = errors.New("chart of account not found")
	ErrCOACodeExists = errors.New("chart of account code already exists")
)

type ChartOfAccountUseCase interface {
	Create(ctx context.Context, req dto.CreateCOARequest) (*dto.ChartOfAccountResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ChartOfAccountResponse, error)
	GetAll(ctx context.Context) ([]dto.ChartOfAccountResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.ChartOfAccountResponse, *entity.Meta, error)
	GetByType(ctx context.Context, accountType string) ([]dto.ChartOfAccountResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateCOARequest) (*dto.ChartOfAccountResponse, error)
}

type chartOfAccountUseCase struct {
	coaRepo repository.ChartOfAccountRepository
	uow    uow.UnitOfWork
}

func NewChartOfAccountUseCase(coaRepo repository.ChartOfAccountRepository, uow uow.UnitOfWork) ChartOfAccountUseCase {
	return &chartOfAccountUseCase{
		coaRepo: coaRepo,
		uow:    uow,
	}
}

func (u *chartOfAccountUseCase) Create(ctx context.Context, req dto.CreateCOARequest) (*dto.ChartOfAccountResponse, error) {
	existing, err := u.coaRepo.FindByCode(ctx, req.AccountCode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCOACodeExists
	}

	coa := &entity.ChartOfAccount{}
	if err := coa.GenerateID(); err != nil {
		return nil, err
	}
	coa.AccountCode = req.AccountCode
	coa.Name = req.Name
	coa.AccountType = req.AccountType
	coa.NormalBalance = req.NormalBalance
	coa.IsActive = true

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.coaRepo.Create(txCtx, coa); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toCOAResponse(coa)
	return &resp, nil
}

func (u *chartOfAccountUseCase) GetByID(ctx context.Context, id string) (*dto.ChartOfAccountResponse, error) {
	coa, err := u.coaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if coa == nil {
		return nil, ErrCOANotFound
	}

	resp := toCOAResponse(coa)
	return &resp, nil
}

func (u *chartOfAccountUseCase) GetAll(ctx context.Context) ([]dto.ChartOfAccountResponse, error) {
	coas, err := u.coaRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.ChartOfAccountResponse
	for _, c := range coas {
		resp = append(resp, toCOAResponse(&c))
	}
	return resp, nil
}

func (u *chartOfAccountUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.ChartOfAccountResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "account_code", "name", "account_type", "updated_at"}
	searchColumns := []string{"id", "account_code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.coaRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.ChartOfAccountResponse
	for _, c := range data {
		resp = append(resp, toCOAResponse(&c))
	}
	return resp, resMeta, nil
}

func (u *chartOfAccountUseCase) GetByType(ctx context.Context, accountType string) ([]dto.ChartOfAccountResponse, error) {
	coas, err := u.coaRepo.FindByType(ctx, accountType)
	if err != nil {
		return nil, err
	}

	var resp []dto.ChartOfAccountResponse
	for _, c := range coas {
		resp = append(resp, toCOAResponse(&c))
	}
	return resp, nil
}

func (u *chartOfAccountUseCase) Update(ctx context.Context, id string, req dto.UpdateCOARequest) (*dto.ChartOfAccountResponse, error) {
	coa, err := u.coaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if coa == nil {
		return nil, ErrCOANotFound
	}

	if req.AccountCode != nil && *req.AccountCode != coa.AccountCode {
		existing, err := u.coaRepo.FindByCode(ctx, *req.AccountCode)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrCOACodeExists
		}
		coa.AccountCode = *req.AccountCode
	}

	if req.Name != nil {
		coa.Name = *req.Name
	}
	if req.AccountType != nil {
		coa.AccountType = *req.AccountType
	}
	if req.NormalBalance != nil {
		coa.NormalBalance = *req.NormalBalance
	}
	if req.IsActive != nil {
		coa.IsActive = *req.IsActive
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.coaRepo.Update(txCtx, coa); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toCOAResponse(coa)
	return &resp, nil
}

func toCOAResponse(coa *entity.ChartOfAccount) dto.ChartOfAccountResponse {
	return dto.ChartOfAccountResponse{
		ID:            coa.ID.String(),
		AccountCode:  coa.AccountCode,
		Name:         coa.Name,
		AccountType:  coa.AccountType,
		NormalBalance: coa.NormalBalance,
		IsActive:     coa.IsActive,
		CreatedAt:    coa.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    coa.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}