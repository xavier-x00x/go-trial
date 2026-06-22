package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/shopspring/decimal"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
var codeRegex = regexp.MustCompile(`^[A-Z]{3}$`)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type ProductCategoryUseCase interface {
	Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error)
	GetAll(ctx context.Context) ([]dto.CategoryResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.CategoryResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id string) error
	GetNextSKU(ctx context.Context, id string) (string, error)
}

type productCategoryUseCase struct {
	categoryRepo repository.ProductCategoryRepository
	uow          uow.UnitOfWork
}

func NewProductCategoryUseCase(categoryRepo repository.ProductCategoryRepository, uow uow.UnitOfWork) ProductCategoryUseCase {
	return &productCategoryUseCase{
		categoryRepo: categoryRepo,
		uow:          uow,
	}
}

func (u *productCategoryUseCase) Create(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	var fe FieldErrors
	if !slugRegex.MatchString(req.Slug) {
		fe.Add("slug", "slug hanya boleh mengandung huruf kecil, angka, dan tanda hubung (-)")
	}
	if req.DefaultMarkupPct.LessThan(decimal.Zero) || req.DefaultMarkupPct.GreaterThan(decimal.NewFromFloat(99999.99)) {
		fe.Add("default_markup_pct", "default_markup_pct harus antara 0 dan 99999.99")
	}

	if req.Code != "" {
		req.Code = strings.ToUpper(req.Code)
		if !codeRegex.MatchString(req.Code) {
			fe.Add("code", "code harus terdiri dari 3 huruf besar")
		}
	} else {
		fe.Add("code", "code harus diisi")
	}

	existing, err := u.categoryRepo.FindBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("slug", "slug sudah digunakan")
		// return nil, ErrCategorySlugExists
	}

	if req.Code != "" {
		existingCode, err := u.categoryRepo.FindByCode(ctx, req.Code)
		if err != nil {
			return nil, err
		}
		if existingCode != nil {
			fe.Add("code", "code sudah digunakan")
		}
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	cat := &entity.ProductCategory{}
	if err := cat.GenerateID(); err != nil {
		return nil, err
	}
	cat.ParentID = req.ParentID
	cat.Name = req.Name
	cat.Slug = req.Slug
	cat.Code = req.Code
	cat.Sequence = req.Sequence
	cat.DefaultMarkupPct = req.DefaultMarkupPct

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

func (u *productCategoryUseCase) GetNextSKU(ctx context.Context, id string) (string, error) {
	cat, err := u.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if cat == nil {
		return "", ErrCategoryNotFound
	}

	nextSeq := cat.Sequence + 1
	// format: KODE-5DIGIT
	return fmt.Sprintf("%s-%05d", cat.Code, nextSeq), nil
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

	var fe FieldErrors

	if req.Slug != nil {
		if !slugRegex.MatchString(*req.Slug) {
			fe.Add("slug", "slug hanya boleh mengandung huruf kecil, angka, dan tanda hubung (-)")
		}
		if req.DefaultMarkupPct != nil {
			if req.DefaultMarkupPct.LessThan(decimal.Zero) || req.DefaultMarkupPct.GreaterThan(decimal.NewFromFloat(99999.99)) {
				fe.Add("default_markup_pct", "default_markup_pct harus antara 0 dan 99999.99")
			}
		}
		if *req.Slug != cat.Slug {
			existing, err := u.categoryRepo.FindBySlug(ctx, *req.Slug)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				fe.Add("slug", "slug sudah digunakan")
			}
			cat.Slug = *req.Slug
		}
	}
	if req.Code != nil {
		codeStr := strings.ToUpper(*req.Code)
		if !codeRegex.MatchString(codeStr) {
			fe.Add("code", "code harus terdiri dari 3 huruf besar")
		} else if codeStr != cat.Code {
			existingCode, err := u.categoryRepo.FindByCode(ctx, codeStr)
			if err != nil {
				return nil, err
			}
			if existingCode != nil {
				fe.Add("code", "code sudah digunakan")
			} else {
				cat.Code = codeStr
			}
		}
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.ParentID != nil {
		if req.ParentID.String() == id {
			fe.Add("parent_id", "parent_id tidak boleh sama dengan ID kategori itu sendiri")
		}
		cat.ParentID = req.ParentID
	}
	if req.DefaultMarkupPct != nil {
		cat.DefaultMarkupPct = *req.DefaultMarkupPct
	}
	if req.Sequence != nil {
		cat.Sequence = *req.Sequence
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
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
	resp := dto.CategoryResponse{
		ID:               cat.ID.String(),
		ParentID:         cat.ParentID,
		Name:             cat.Name,
		Slug:             cat.Slug,
		Code:             cat.Code,
		Sequence:         cat.Sequence,
		DefaultMarkupPct: cat.DefaultMarkupPct,
		CreatedAt:        cat.CreatedAt,
		UpdatedAt:        cat.UpdatedAt,
	}
	if cat.Parent != nil {
		p := toCategoryResponse(cat.Parent)
		p.Parent = nil
		resp.Parent = &p
	}
	return resp
}

func toCategoryListResponse(cat *entity.ProductCategory) dto.CategoryResponse {
	resp := dto.CategoryResponse{
		ID:               cat.ID.String(),
		ParentID:         cat.ParentID,
		Name:             cat.Name,
		Slug:             cat.Slug,
		Code:             cat.Code,
		Sequence:         cat.Sequence,
		DefaultMarkupPct: cat.DefaultMarkupPct,
		CreatedAt:        cat.CreatedAt,
		UpdatedAt:        cat.UpdatedAt,
	}
	if cat.Parent != nil {
		p := toCategoryResponse(cat.Parent)
		p.Parent = nil
		resp.Parent = &p
	}
	return resp
}
