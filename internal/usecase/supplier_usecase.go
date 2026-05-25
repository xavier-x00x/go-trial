package usecase

import (
	"context"
	"errors"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/shopspring/decimal"
)

var (
	ErrSupplierNotFound  = errors.New("supplier not found")
	ErrSupplierNotActive = errors.New("supplier is not active")
)

type SupplierUseCase interface {
	Create(ctx context.Context, req dto.CreateSupplierRequest) (*dto.SupplierResponse, error)
	GetByID(ctx context.Context, id string) (*dto.SupplierResponse, error)
	GetAll(ctx context.Context) ([]dto.SupplierResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.SupplierResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateSupplierRequest) (*dto.SupplierResponse, error)
	Delete(ctx context.Context, id string) error
}

type supplierUseCase struct {
	supplierRepo repository.SupplierRepository
	coaRepo      repository.ChartOfAccountRepository
	uow          uow.UnitOfWork
}

func NewSupplierUseCase(supplierRepo repository.SupplierRepository, coaRepo repository.ChartOfAccountRepository, uow uow.UnitOfWork) SupplierUseCase {
	return &supplierUseCase{
		supplierRepo: supplierRepo,
		coaRepo:      coaRepo,
		uow:          uow,
	}
}

func (u *supplierUseCase) Create(ctx context.Context, req dto.CreateSupplierRequest) (*dto.SupplierResponse, error) {
	var fe FieldErrors

	existing, err := u.supplierRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("code", "kode supplier sudah digunakan")
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	supplier := &entity.Supplier{}
	if err := supplier.GenerateID(); err != nil {
		return nil, err
	}
	supplier.Code = req.Code
	supplier.Name = req.Name
	supplier.ContactPerson = req.ContactPerson
	supplier.ContactPhone = req.ContactPhone
	supplier.PhoneNumber = req.PhoneNumber
	supplier.Email = req.Email
	supplier.PreferredNotificationMethod = req.PreferredNotificationMethod
	supplier.Address = req.Address
	supplier.TaxRegNumber = req.TaxRegNumber
	supplier.SupplierCategoryID = req.SupplierCategoryID
	supplier.IsPKP = req.IsPKP
	supplier.PaymentTermDays = req.PaymentTermDays
	supplier.PaymentMode = req.PaymentMode
	if req.MinOrderAmount != "" {
		d, _ := decimal.NewFromString(string(req.MinOrderAmount))
		supplier.MinOrderAmount = d
	}
	supplier.BankName = req.BankName
	supplier.BankAccount = req.BankAccount
	supplier.BankAccountName = req.BankAccountName
	supplier.APAccountID = req.APAccountID

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.supplierRepo.Create(txCtx, supplier); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toSupplierResponse(supplier, nil)
	return &resp, nil
}

func (u *supplierUseCase) GetByID(ctx context.Context, id string) (*dto.SupplierResponse, error) {
	supplier, err := u.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}

	resp := toSupplierResponse(supplier, nil)
	return &resp, nil
}

func (u *supplierUseCase) GetAll(ctx context.Context) ([]dto.SupplierResponse, error) {
	suppliers, err := u.supplierRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.SupplierResponse
	for _, s := range suppliers {
		r := toSupplierResponse(&s, nil)
		resp = append(resp, r)
	}
	return resp, nil
}

func (u *supplierUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.SupplierResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.supplierRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.SupplierResponse{}
	for _, s := range data {
		resp = append(resp, toSupplierResponse(&s, nil))
	}
	return resp, resMeta, nil
}

func (u *supplierUseCase) Update(ctx context.Context, id string, req dto.UpdateSupplierRequest) (*dto.SupplierResponse, error) {
	supplier, err := u.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}

	if req.Name != nil {
		supplier.Name = *req.Name
	}
	if req.ContactPerson != nil {
		supplier.ContactPerson = req.ContactPerson
	}
	if req.ContactPhone != nil {
		supplier.ContactPhone = req.ContactPhone
	}
	if req.PhoneNumber != nil {
		supplier.PhoneNumber = req.PhoneNumber
	}
	if req.Email != nil {
		supplier.Email = req.Email
	}
	if req.PreferredNotificationMethod != nil {
		supplier.PreferredNotificationMethod = *req.PreferredNotificationMethod
	}
	if req.Address != nil {
		supplier.Address = req.Address
	}
	if req.TaxRegNumber != nil {
		supplier.TaxRegNumber = req.TaxRegNumber
	}
	if req.SupplierCategoryID != nil {
		supplier.SupplierCategoryID = req.SupplierCategoryID
	}
	if req.IsPKP != nil {
		supplier.IsPKP = *req.IsPKP
	}
	if req.PaymentTermDays != nil {
		supplier.PaymentTermDays = *req.PaymentTermDays
	}
	if req.PaymentMode != nil {
		supplier.PaymentMode = *req.PaymentMode
	}
	if req.MinOrderAmount != nil {
		if *req.MinOrderAmount != "" {
			d, _ := decimal.NewFromString(string(*req.MinOrderAmount))
			supplier.MinOrderAmount = d
		}
	}
	if req.BankName != nil {
		supplier.BankName = req.BankName
	}
	if req.BankAccount != nil {
		supplier.BankAccount = req.BankAccount
	}
	if req.BankAccountName != nil {
		supplier.BankAccountName = req.BankAccountName
	}
	if req.IsActive != nil {
		supplier.IsActive = *req.IsActive
	}
	if req.APAccountID != nil {
		supplier.APAccountID = req.APAccountID
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.supplierRepo.Update(txCtx, supplier); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toSupplierResponse(supplier, nil)
	return &resp, nil
}

func (u *supplierUseCase) Delete(ctx context.Context, id string) error {
	supplier, err := u.supplierRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if supplier == nil {
		return ErrSupplierNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.supplierRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toSupplierResponse(s *entity.Supplier, coa *dto.ChartOfAccountResponse) dto.SupplierResponse {
	return dto.SupplierResponse{
		ID:                          s.ID,
		Code:                        s.Code,
		Name:                        s.Name,
		ContactPerson:               s.ContactPerson,
		ContactPhone:                s.ContactPhone,
		PhoneNumber:                 s.PhoneNumber,
		Email:                       s.Email,
		PreferredNotificationMethod: s.PreferredNotificationMethod,
		Address:                     s.Address,
		TaxRegNumber:                s.TaxRegNumber,
		SupplierCategoryID:          s.SupplierCategoryID,
		SupplierCategory:            nil,
		IsPKP:                       s.IsPKP,
		PaymentTermDays:             s.PaymentTermDays,
		PaymentMode:                 s.PaymentMode,
		MinOrderAmount:              s.MinOrderAmount,
		BankName:                    s.BankName,
		BankAccount:                 s.BankAccount,
		BankAccountName:             s.BankAccountName,
		IsActive:                    s.IsActive,
		APAccountID:                 s.APAccountID,
		APAccount:                   coa,
		CreatedAt:                   s.CreatedAt,
		UpdatedAt:                   s.UpdatedAt,
	}
}
