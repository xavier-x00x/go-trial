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
	ErrStoreNotFound = errors.New("store not found")
)

// StoreUseCase defines the business logic for store management.
type StoreUseCase interface {
	Create(ctx context.Context, req dto.CreateStoreRequest) (*dto.StoreResponse, error)
	GetByID(ctx context.Context, id string) (*dto.StoreResponse, error)
	GetAll(ctx context.Context) ([]dto.StoreResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.StoreResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateStoreRequest) (*dto.StoreResponse, error)
	Delete(ctx context.Context, id string) error
}

type storeUseCase struct {
	storeRepo repository.StoreRepository
	uow       uow.UnitOfWork
}

func NewStoreUseCase(
	storeRepo repository.StoreRepository,
	uow uow.UnitOfWork,
) StoreUseCase {
	return &storeUseCase{
		storeRepo: storeRepo,
		uow:       uow,
	}
}

// Create creates a new store within a UoW transaction.
func (u *storeUseCase) Create(ctx context.Context, req dto.CreateStoreRequest) (*dto.StoreResponse, error) {
	var fe FieldErrors

	existing, err := u.storeRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("code", "kode toko sudah digunakan")
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	store := &entity.Store{
		BaseModel:  entity.BaseModel{ID: id},
		Code:       req.Code,
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Phone:      req.Phone,
		Email:      req.Email,
		IsMain:     req.IsMain,
		IsActive:   true,
	}

	// Start Unit of Work
	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx) //nolint:errcheck

	if err := u.storeRepo.Create(txCtx, store); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}
	// End Unit of Work

	resp := toStoreResponse(store)
	return &resp, nil
}

// GetByID returns a single store by ID.
func (u *storeUseCase) GetByID(ctx context.Context, id string) (*dto.StoreResponse, error) {
	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	resp := toStoreResponse(store)
	return &resp, nil
}

// GetAll returns all stores.
func (u *storeUseCase) GetAll(ctx context.Context) ([]dto.StoreResponse, error) {
	stores, err := u.storeRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.StoreResponse
	for _, s := range stores {
		resp = append(resp, toStoreResponse(&s))
	}
	return resp, nil
}

func (u *storeUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.StoreResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "code", "name", "updated_at"}
	searchColumns := []string{"id", "code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["is_active"] = 1

	data, resMeta, err := u.storeRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.StoreResponse{}
	for _, s := range data {
		resp = append(resp, toStoreResponse(&s))
	}

	return resp, resMeta, nil
}

// applyIfSet assigns *src to *dst if src is not nil.
func applyIfSet[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}

// Update updates an existing store within a UoW transaction.
func (u *storeUseCase) Update(ctx context.Context, id string, req dto.UpdateStoreRequest) (*dto.StoreResponse, error) {
	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, ErrStoreNotFound
	}

	if req.Code != nil && *req.Code != store.Code {
		var fe FieldErrors
		existing, err := u.storeRepo.FindByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			fe.Add("code", "kode toko sudah digunakan")
		}
		if len(fe.Errors) > 0 {
			return nil, &fe
		}
		store.Code = *req.Code
	}

	// Apply partial updates — value fields (nil = don't change)
	applyIfSet(&store.Name, req.Name)
	applyIfSet(&store.IsMain, req.IsMain)
	applyIfSet(&store.IsActive, req.IsActive)

	// Nullable fields — pointer to pointer, assign directly
	store.TaxRegNumber = req.NPWP
	store.Address = req.Address
	store.City = req.City
	store.Province = req.Province
	store.PostalCode = req.PostalCode
	store.Phone = req.Phone
	store.Email = req.Email

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx) //nolint:errcheck

	if err := u.storeRepo.Update(txCtx, store); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toStoreResponse(store)
	return &resp, nil
}

// Delete soft-deletes a store.
func (u *storeUseCase) Delete(ctx context.Context, id string) error {
	store, err := u.storeRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if store == nil {
		return ErrStoreNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx) //nolint:errcheck

	if err := u.storeRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toStoreResponse(store *entity.Store) dto.StoreResponse {
	return dto.StoreResponse{
		ID:         store.ID.String(),
		Code:       store.Code,
		Name:       store.Name,
		Address:    store.Address,
		City:       store.City,
		Province:   store.Province,
		PostalCode: store.PostalCode,
		Phone:      store.Phone,
		Email:      store.Email,
		IsMain:     store.IsMain,
		IsActive:   store.IsActive,
		CreatedAt:  store.CreatedAt,
		UpdatedAt:  store.UpdatedAt,
	}
}
