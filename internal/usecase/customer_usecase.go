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
	ErrCustomerNotFound = errors.New("customer not found")
)

type CustomerUseCase interface {
	Create(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error)
	GetByID(ctx context.Context, id string) (*dto.CustomerResponse, error)
	GetAll(ctx context.Context) ([]dto.CustomerResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.CustomerResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error)
	Delete(ctx context.Context, id string) error
}

type customerUseCase struct {
	customerRepo repository.CustomerRepository
	coaRepo      repository.ChartOfAccountRepository
	uow         uow.UnitOfWork
}

func NewCustomerUseCase(customerRepo repository.CustomerRepository, coaRepo repository.ChartOfAccountRepository, uow uow.UnitOfWork) CustomerUseCase {
	return &customerUseCase{
		customerRepo: customerRepo,
		coaRepo:      coaRepo,
		uow:         uow,
	}
}

func (u *customerUseCase) Create(ctx context.Context, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	var fe FieldErrors

	existing, err := u.customerRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("code", "kode pelanggan sudah digunakan")
	}

	if req.PhoneNumber != nil {
		existing, err := u.customerRepo.FindByPhone(ctx, *req.PhoneNumber)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			fe.Add("phone_number", "nomor telepon sudah digunakan")
		}
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	pointBalance := req.PointBalance
	if pointBalance.IsZero() {
		pointBalance = req.PointBalance
	}

	customer := &entity.Customer{}
	if err := customer.GenerateID(); err != nil {
		return nil, err
	}
	customer.Code = req.Code
	customer.Name = req.Name
	customer.PhoneNumber = req.PhoneNumber
	customer.Email = req.Email
	customer.Address = req.Address
	customer.PointBalance = pointBalance
	customer.CreditLimit = req.CreditLimit
	customer.ARAccountID = req.ARAccountID

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.customerRepo.Create(txCtx, customer); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toCustomerResponse(customer, nil)
	return &resp, nil
}

func (u *customerUseCase) GetByID(ctx context.Context, id string) (*dto.CustomerResponse, error) {
	customer, err := u.customerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	resp := toCustomerResponse(customer, nil)
	return &resp, nil
}

func (u *customerUseCase) GetAll(ctx context.Context) ([]dto.CustomerResponse, error) {
	customers, err := u.customerRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.CustomerResponse
	for _, c := range customers {
		r := toCustomerResponse(&c, nil)
		resp = append(resp, r)
	}
	return resp, nil
}

func (u *customerUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.CustomerResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.customerRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.CustomerResponse
	for _, c := range data {
		resp = append(resp, toCustomerResponse(&c, nil))
	}
	return resp, resMeta, nil
}

func (u *customerUseCase) Update(ctx context.Context, id string, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	customer, err := u.customerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	if req.Code != nil && *req.Code != customer.Code {
		var fe FieldErrors
		existing, err := u.customerRepo.FindByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			fe.Add("code", "kode pelanggan sudah digunakan")
		}
		if len(fe.Errors) > 0 {
			return nil, &fe
		}
		customer.Code = *req.Code
	}

	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		customer.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		customer.Email = req.Email
	}
	if req.Address != nil {
		customer.Address = req.Address
	}
	if req.IsActive != nil {
		customer.IsActive = *req.IsActive
	}
	if req.PointBalance != nil {
		customer.PointBalance = *req.PointBalance
	}
	if req.CreditLimit != nil {
		customer.CreditLimit = *req.CreditLimit
	}
	if req.ARAccountID != nil {
		customer.ARAccountID = req.ARAccountID
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.customerRepo.Update(txCtx, customer); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toCustomerResponse(customer, nil)
	return &resp, nil
}

func (u *customerUseCase) Delete(ctx context.Context, id string) error {
	customer, err := u.customerRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrCustomerNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.customerRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toCustomerResponse(c *entity.Customer, coa *dto.ChartOfAccountResponse) dto.CustomerResponse {
	return dto.CustomerResponse{
		ID:           c.ID,
		Code:        c.Code,
		Name:        c.Name,
		PhoneNumber: c.PhoneNumber,
		Email:       c.Email,
		Address:     c.Address,
		IsActive:    c.IsActive,
		PointBalance: c.PointBalance,
		CreditLimit: c.CreditLimit,
		ARAccountID: c.ARAccountID,
		ARAccount:   coa,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}