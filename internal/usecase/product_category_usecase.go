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
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category slug already exists")
)

type ProductCategoryUseCase interface {
	Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error)
	GetAll(ctx context.Context) ([]dto.CategoryResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.CategoryResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id string) error
}

type productCategoryUseCase struct {
	categoryRepo repository.ProductCategoryRepository
	uow         uow.UnitOfWork
}

func NewProductCategoryUseCase(categoryRepo repository.ProductCategoryRepository, uow uow.UnitOfWork) ProductCategoryUseCase {
	return &productCategoryUseCase{
		categoryRepo: categoryRepo,
		uow:        uow,
	}
}

func (u *productCategoryUseCase) Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	existing, err := u.categoryRepo.FindBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCategorySlugExists
	}

	cat := &entity.ProductCategory{}
	if err := cat.GenerateID(); err != nil {
		return nil, err
	}
	cat.ParentID = req.ParentID
	cat.Name = req.Name
	cat.Slug = req.Slug

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

	resp := toCategoryResponse(cat)
	return &resp, nil
}

func (u *productCategoryUseCase) GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	cat, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	resp := toCategoryResponse(cat)
	return &resp, nil
}

func (u *productCategoryUseCase) GetAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	cats, err := u.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.CategoryResponse
	for _, cat := range cats {
		resp = append(resp, toCategoryResponse(&cat))
	}
	return resp, nil
}

func (u *productCategoryUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.CategoryResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "slug", "created_at"}
	searchColumns := []string{"id", "name", "slug"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.categoryRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.CategoryResponse
	for _, cat := range data {
		resp = append(resp, toCategoryResponse(&cat))
	}
	return resp, resMeta, nil
}

func (u *productCategoryUseCase) Update(ctx context.Context, id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	cat, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Slug != nil && *req.Slug != cat.Slug {
		existing, err := u.categoryRepo.FindBySlug(ctx, *req.Slug)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, ErrCategorySlugExists
		}
		cat.Slug = *req.Slug
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.ParentID != nil {
		cat.ParentID = req.ParentID
	}
	if req.DefaultMarkupPct != nil {
		cat.DefaultMarkupPct = *req.DefaultMarkupPct
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

	resp := toCategoryResponse(cat)
	return &resp, nil
}

func (u *productCategoryUseCase) Delete(ctx context.Context, id string) error {
	cat, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrCategoryNotFound
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

func toCategoryResponse(cat *entity.ProductCategory) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:                cat.ID.String(),
		ParentID:         cat.ParentID,
		Name:             cat.Name,
		Slug:             cat.Slug,
		DefaultMarkupPct:  cat.DefaultMarkupPct,
		CreatedAt:        cat.CreatedAt,
		UpdatedAt:        cat.UpdatedAt,
	}
}