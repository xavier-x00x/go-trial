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
	ErrProductNotFound = errors.New("product not found")
)

type ProductUseCase interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ProductResponse, error)
	GetAll(ctx context.Context) ([]dto.ProductResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.ProductResponse, *entity.Meta, error)
	Update(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(ctx context.Context, id string) error
}

type productUseCase struct {
	productRepo    repository.ProductRepository
	categoryRepo   repository.ProductCategoryRepository
	uomRepo        repository.UOMRepository
	uow            uow.UnitOfWork
}

func NewProductUseCase(
	productRepo repository.ProductRepository,
	categoryRepo repository.ProductCategoryRepository,
	uomRepo repository.UOMRepository,
	uow uow.UnitOfWork,
) ProductUseCase {
	return &productUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		uomRepo:     uomRepo,
		uow:         uow,
	}
}

func (u *productUseCase) Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	var fe FieldErrors

	existing, err := u.productRepo.FindBySKU(ctx, req.SKU)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("sku", "sku sudah digunakan")
	}

	if req.Barcode != nil {
		existing, err := u.productRepo.FindByBarcode(ctx, *req.Barcode)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			fe.Add("barcode", "barcode sudah digunakan")
		}
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	product := &entity.Product{}
	if err := product.GenerateID(); err != nil {
		return nil, err
	}
	product.SKU = req.SKU
	product.Barcode = req.Barcode
	product.Name = req.Name
	product.CategoryID = req.CategoryID
	product.BaseUOMID = req.BaseUOMID
	product.IsStockable = req.IsStockable
	product.Length = req.Length
	product.Width = req.Width
	product.Height = req.Height
	product.Weight = req.Weight
	product.IsStackable = req.IsStackable
	product.MaxStackLayer = req.MaxStackLayer

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.productRepo.Create(txCtx, product); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toProductResponse(product, nil, nil)
	return &resp, nil
}

func (u *productUseCase) GetByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	resp := toProductResponse(product, nil, nil)
	return &resp, nil
}

func (u *productUseCase) GetAll(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := u.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.ProductResponse
	for _, p := range products {
		r := toProductResponse(&p, nil, nil)
		resp = append(resp, r)
	}
	return resp, nil
}

func (u *productUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.ProductResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "sku", "name", "updated_at"}
	searchColumns := []string{"id", "sku", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.productRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.ProductResponse
	for _, p := range data {
		resp = append(resp, toProductResponse(&p, nil, nil))
	}
	return resp, resMeta, nil
}

func (u *productUseCase) Update(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	if req.Barcode != nil && *req.Barcode != "" {
		var fe FieldErrors
		if product.Barcode == nil || *req.Barcode != *product.Barcode {
			existing, err := u.productRepo.FindByBarcode(ctx, *req.Barcode)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				fe.Add("barcode", "barcode sudah digunakan")
			}
		}
		if len(fe.Errors) > 0 {
			return nil, &fe
		}
		product.Barcode = req.Barcode
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.BaseUOMID != nil {
		product.BaseUOMID = *req.BaseUOMID
	}
	if req.IsStockable != nil {
		product.IsStockable = *req.IsStockable
	}
	if req.Length != nil {
		product.Length = *req.Length
	}
	if req.Width != nil {
		product.Width = *req.Width
	}
	if req.Height != nil {
		product.Height = *req.Height
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.IsStackable != nil {
		product.IsStackable = *req.IsStackable
	}
	if req.MaxStackLayer != nil {
		product.MaxStackLayer = *req.MaxStackLayer
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.productRepo.Update(txCtx, product); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toProductResponse(product, nil, nil)
	return &resp, nil
}

func (u *productUseCase) Delete(ctx context.Context, id string) error {
	product, err := u.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.productRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toProductResponse(p *entity.Product, category *dto.CategoryResponse, uom *dto.UOMResponse) dto.ProductResponse {
	// Auto-map dari entity jika preloaded tapi parameter DTO nil
	if uom == nil && p.BaseUOM.ID != uuid.Nil {
		uom = &dto.UOMResponse{
			ID:   p.BaseUOM.ID.String(),
			Code: p.BaseUOM.Code,
			Name: p.BaseUOM.Name,
		}
	}
	if category == nil && p.Category.ID != uuid.Nil {
		category = &dto.CategoryResponse{
			ID:   p.Category.ID.String(),
			Name: p.Category.Name,
			Slug: p.Category.Slug,
		}
	}

	return dto.ProductResponse{
		ID:         p.ID,
		SKU:        p.SKU,
		Barcode:    p.Barcode,
		Name:       p.Name,
		CategoryID: p.CategoryID,
		Category:  category,
		BaseUOMID: p.BaseUOMID,
		BaseUOM:   uom,
		IsStockable: p.IsStockable,
		Length:    p.Length,
		Width:     p.Width,
		Height:    p.Height,
		Weight:    p.Weight,
		IsStackable: p.IsStackable,
		MaxStackLayer: p.MaxStackLayer,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}