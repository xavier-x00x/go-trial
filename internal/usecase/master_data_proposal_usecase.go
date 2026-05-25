package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
)

var (
	ErrProposalNotFound   = errors.New("proposal not found")
	ErrProposalNotPending = errors.New("proposal is not pending")
	ErrProposalRejected   = errors.New("proposal already rejected")
	ErrInvalidAction      = errors.New("invalid action type")
)

type MasterDataProposalUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreateMasterDataProposalRequest) (*dto.MasterDataProposalDetailResponse, error)
	GetByID(ctx context.Context, id string) (*dto.MasterDataProposalDetailResponse, error)
	GetAll(ctx context.Context) ([]dto.MasterDataProposalListResponse, error)
	GetPending(ctx context.Context) ([]dto.MasterDataProposalListResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.MasterDataProposalListResponse, *entity.Meta, error)
	GetByEntityType(ctx context.Context, entityType, status string) ([]dto.MasterDataProposalListResponse, error)
	GetByGroup(ctx context.Context, groupID string) ([]dto.MasterDataProposalDetailResponse, error)
	Review(ctx context.Context, userID string, id string, req dto.ReviewMasterDataProposalRequest) (*dto.MasterDataProposalDetailResponse, error)
	Update(ctx context.Context, userID string, id string, req dto.UpdateMasterDataProposalRequest) (*dto.MasterDataProposalDetailResponse, error)
	Execute(ctx context.Context, id string) error
	BulkLinkProductSupplier(ctx context.Context, userID string, req dto.BulkCreateProductSupplierProposalRequest) (*dto.BulkProposalResponse, error)
}

type masterDataProposalUseCaseImpl struct {
	repo                repository.MasterDataProposalRepository
	productRepo         repository.ProductRepository
	productPriceRepo    repository.ProductPriceRepository
	productUOMRepo      repository.ProductUOMConversionRepository
	supplierRepo        repository.SupplierRepository
	productSupplierRepo repository.ProductSupplierRepository
	coaRepo             repository.ChartOfAccountRepository
	taxRepo             repository.TaxRepository
	numberSequenceRepo  repository.NumberSequenceRepository
	uow                 uow.UnitOfWork
}

type MasterDataProposalUseCaseConfig struct {
	Repo                repository.MasterDataProposalRepository
	ProductRepo         repository.ProductRepository
	ProductPriceRepo    repository.ProductPriceRepository
	ProductUOMRepo      repository.ProductUOMConversionRepository
	SupplierRepo        repository.SupplierRepository
	ProductSupplierRepo repository.ProductSupplierRepository
	CoaRepo             repository.ChartOfAccountRepository
	TaxRepo             repository.TaxRepository
	NumberSequenceRepo  repository.NumberSequenceRepository
	Uow                 uow.UnitOfWork
}

func NewMasterDataProposalUseCase(cfg MasterDataProposalUseCaseConfig) MasterDataProposalUseCase {
	return &masterDataProposalUseCaseImpl{
		repo:                cfg.Repo,
		productRepo:         cfg.ProductRepo,
		productPriceRepo:    cfg.ProductPriceRepo,
		productUOMRepo:      cfg.ProductUOMRepo,
		supplierRepo:        cfg.SupplierRepo,
		productSupplierRepo: cfg.ProductSupplierRepo,
		coaRepo:             cfg.CoaRepo,
		taxRepo:             cfg.TaxRepo,
		numberSequenceRepo:  cfg.NumberSequenceRepo,
		uow:                 cfg.Uow,
	}
}

func (u *masterDataProposalUseCaseImpl) generateReferenceNumber(entityType string, date time.Time) (string, error) {
	prefix := getEntityPrefix(entityType)
	period := date.Format("0601") // YYMM

	seqNum, err := u.numberSequenceRepo.GetNextNumber(context.Background(), prefix, period)
	if err != nil {
		return "", err
	}

	seqStr := fmt.Sprintf("%05d", seqNum)
	// Format: PREFIX/YYMM/NNNNN (contoh: PRD/2409/00001)
	return fmt.Sprintf("%s/%s/%s", prefix, period, seqStr), nil
}

func getEntityPrefix(entityType string) string {
	switch entityType {
	case "PRODUCT":
		return "PRD"
	case "PRODUCT_PRICE":
		return "PRC"
	case "PRODUCT_UOM_CONVERSION":
		return "PUC"
	case "SUPPLIER":
		return "SUP"
	case "PRODUCT_SUPPLIER":
		return "PSP"
	case "CHART_OF_ACCOUNT":
		return "COA"
	case "TAX":
		return "TAX"
	default:
		return entityType[:min(3, len(entityType))]
	}
}

func (u *masterDataProposalUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreateMasterDataProposalRequest) (*dto.MasterDataProposalDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	refNum, err := u.generateReferenceNumber(req.EntityType, time.Now())
	if err != nil {
		return nil, err
	}

	items := make([]entity.MasterDataProposalItem, len(req.Items))
	for i, item := range req.Items {
		payloadJSON, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		items[i] = entity.MasterDataProposalItem{
			SeqNo:       i + 1,
			EntityID:    item.EntityID,
			PayloadJSON: string(payloadJSON),
		}
		if err := items[i].GenerateID(); err != nil {
			return nil, err
		}
	}

	proposal := &entity.MasterDataProposal{
		ReferenceNumber: refNum,
		EntityType:      req.EntityType,
		ActionType:      req.ActionType,
		TotalItems:      len(items),
		Reason:          req.Reason,
		Status:          entity.ProposalStatusPending,
		ProposedByID:    userUUID,
		Items:           items,
	}

	if err := proposal.GenerateID(); err != nil {
		return nil, err
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.repo.Create(txCtx, proposal); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toMasterDataProposalDetailResponse(proposal), nil
}

func (u *masterDataProposalUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.MasterDataProposalDetailResponse, error) {
	proposal, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if proposal == nil {
		return nil, ErrProposalNotFound
	}

	resp := toMasterDataProposalDetailResponse(proposal)
	return resp, nil
}

func (u *masterDataProposalUseCaseImpl) GetPending(ctx context.Context) ([]dto.MasterDataProposalListResponse, error) {
	proposals, err := u.repo.FindPending(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.MasterDataProposalListResponse
	for _, p := range proposals {
		resp = append(resp, *toMasterDataProposalListResponse(&p))
	}
	return resp, nil
}

func (u *masterDataProposalUseCaseImpl) GetAll(ctx context.Context) ([]dto.MasterDataProposalListResponse, error) {
	proposals, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.MasterDataProposalListResponse
	for _, p := range proposals {
		resp = append(resp, *toMasterDataProposalListResponse(&p))
	}
	return resp, nil
}

func (u *masterDataProposalUseCaseImpl) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.MasterDataProposalListResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "reference_number", "entity_type", "status", "created_at"}
	searchColumns := []string{"id", "reference_number"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.MasterDataProposalListResponse{}
	for _, p := range data {
		resp = append(resp, *toMasterDataProposalListResponse(&p))
	}
	return resp, resMeta, nil
}

func (u *masterDataProposalUseCaseImpl) GetByEntityType(ctx context.Context, entityType, status string) ([]dto.MasterDataProposalListResponse, error) {
	proposals, err := u.repo.FindByEntityType(ctx, entityType, status)
	if err != nil {
		return nil, err
	}

	var resp []dto.MasterDataProposalListResponse
	for _, p := range proposals {
		resp = append(resp, *toMasterDataProposalListResponse(&p))
	}
	return resp, nil
}

func (u *masterDataProposalUseCaseImpl) Review(ctx context.Context, userID string, id string, req dto.ReviewMasterDataProposalRequest) (*dto.MasterDataProposalDetailResponse, error) {
	proposal, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if proposal == nil {
		return nil, ErrProposalNotFound
	}

	if proposal.Status != entity.ProposalStatusPending {
		return nil, ErrProposalNotPending
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	if req.Action == "APPROVE" {
		proposal.Status = entity.ProposalStatusApproved
	} else {
		proposal.Status = entity.ProposalStatusRejected
	}
	proposal.ReviewedByID = &userUUID
	now := time.Now()
	proposal.ReviewedAt = &now
	proposal.ReviewNotes = req.ReviewNotes

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.repo.Update(txCtx, proposal); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toMasterDataProposalDetailResponse(proposal), nil
}

func (u *masterDataProposalUseCaseImpl) Update(ctx context.Context, userID string, id string, req dto.UpdateMasterDataProposalRequest) (*dto.MasterDataProposalDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	proposal, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if proposal == nil {
		return nil, ErrProposalNotFound
	}

	if proposal.Status != entity.ProposalStatusPending {
		return nil, ErrProposalNotPending
	}

	// Validate ownership
	if proposal.ProposedByID != userUUID {
		return nil, errors.New("unauthorized: only the proposer can update the proposal")
	}

	items := make([]entity.MasterDataProposalItem, len(req.Items))
	for i, item := range req.Items {
		payloadJSON, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		items[i] = entity.MasterDataProposalItem{
			ProposalID:  proposal.ID,
			SeqNo:       i + 1,
			EntityID:    item.EntityID,
			PayloadJSON: string(payloadJSON),
		}
		if err := items[i].GenerateID(); err != nil {
			return nil, err
		}
	}

	proposal.Reason = req.Reason
	proposal.TotalItems = len(items)
	proposal.Items = items

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	// Delete old items
	if err := u.repo.DeleteItemsByProposalID(txCtx, id); err != nil {
		return nil, err
	}

	// Save header and new items
	if err := u.repo.Update(txCtx, proposal); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toMasterDataProposalDetailResponse(proposal), nil
}

func (u *masterDataProposalUseCaseImpl) Execute(ctx context.Context, id string) error {
	proposal, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if proposal == nil {
		return ErrProposalNotFound
	}

	if proposal.Status != entity.ProposalStatusApproved {
		return ErrProposalNotPending
	}

	switch proposal.EntityType {
	case entity.ProposalEntityProduct:
		return u.executeProduct(ctx, proposal)
	case entity.ProposalEntityProductPrice:
		return u.executeProductPrice(ctx, proposal)
	case entity.ProposalEntityProductUOM:
		return u.executeProductUOMConversion(ctx, proposal)
	case entity.ProposalEntitySupplier:
		return u.executeSupplier(ctx, proposal)
	case entity.ProposalEntityProductSupplier:
		return u.executeProductSupplier(ctx, proposal)
	case entity.ProposalEntityChartOfAccount:
		return u.executeChartOfAccount(ctx, proposal)
	case entity.ProposalEntityTax:
		return u.executeTax(ctx, proposal)
	default:
		return ErrInvalidAction
	}
}

func (u *masterDataProposalUseCaseImpl) executeProduct(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateProductRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateProduct(ctx, u.productRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateProductRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateProduct(ctx, u.productRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := executeDeleteProduct(ctx, u.productRepo, item.EntityID.String(), u.uow); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *masterDataProposalUseCaseImpl) executeProductPrice(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateProductPriceRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateProductPrice(ctx, u.productPriceRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateProductPriceRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateProductPrice(ctx, u.productPriceRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := u.productPriceRepo.Delete(ctx, item.EntityID.String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *masterDataProposalUseCaseImpl) executeProductUOMConversion(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateProductUOMConversionRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateProductUOMConversion(ctx, u.productUOMRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateProductUOMConversionRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateProductUOMConversion(ctx, u.productUOMRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := u.productUOMRepo.Delete(ctx, item.EntityID.String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *masterDataProposalUseCaseImpl) executeSupplier(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateSupplierRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateSupplier(ctx, u.supplierRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateSupplierRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateSupplier(ctx, u.supplierRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := u.supplierRepo.Delete(ctx, item.EntityID.String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *masterDataProposalUseCaseImpl) executeProductSupplier(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateProductSupplierRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateProductSupplier(ctx, u.productSupplierRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateProductSupplierRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateProductSupplier(ctx, u.productSupplierRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := u.productSupplierRepo.Delete(ctx, item.EntityID.String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *masterDataProposalUseCaseImpl) executeChartOfAccount(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateChartOfAccountRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateChartOfAccount(ctx, u.coaRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateChartOfAccountRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateChartOfAccount(ctx, u.coaRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := u.coaRepo.Delete(ctx, item.EntityID.String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *masterDataProposalUseCaseImpl) executeTax(ctx context.Context, p *entity.MasterDataProposal) error {
	for _, item := range p.Items {
		switch p.ActionType {
		case entity.ProposalActionCreate:
			var req dto.CreateTaxRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeCreateTax(ctx, u.taxRepo, &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionUpdate:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			var req dto.UpdateTaxRequest
			if err := json.Unmarshal([]byte(item.PayloadJSON), &req); err != nil {
				return err
			}
			if err := executeUpdateTax(ctx, u.taxRepo, item.EntityID.String(), &req, u.uow); err != nil {
				return err
			}
		case entity.ProposalActionDelete:
			if item.EntityID == nil {
				return ErrProposalNotFound
			}
			if err := u.taxRepo.Delete(ctx, item.EntityID.String()); err != nil {
				return err
			}
		}
	}
	return nil
}

func executeCreateProduct(ctx context.Context, repo repository.ProductRepository, req *dto.CreateProductRequest, uow uow.UnitOfWork) error {
	product := &entity.Product{
		SKU:         req.SKU,
		Barcode:     req.Barcode,
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		BaseUOMID:   req.BaseUOMID,
		IsStockable: req.IsStockable,
		Length:      req.Length,
		Width:       req.Width,
		Height:      req.Height,
		Weight:      req.Weight,
		IsStackable: req.IsStackable,
	}
	if err := product.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, product); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateProduct(ctx context.Context, repo repository.ProductRepository, id string, req *dto.UpdateProductRequest, uow uow.UnitOfWork) error {
	product, err := repo.FindByID(ctx, id)
	if err != nil || product == nil {
		return err
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Barcode != nil {
		product.Barcode = req.Barcode
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.BaseUOMID != nil {
		product.BaseUOMID = *req.BaseUOMID
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, product); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeDeleteProduct(ctx context.Context, repo repository.ProductRepository, id string, uow uow.UnitOfWork) error {
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Delete(txCtx, id); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeCreateProductPrice(ctx context.Context, repo repository.ProductPriceRepository, req *dto.CreateProductPriceRequest, uow uow.UnitOfWork) error {
	pp := &entity.ProductPrice{
		PriceListID: req.PriceListID,
		ProductID:   req.ProductID,
		UOMID:       req.UOMID,
		MarkupPct:   req.MarkupPct,
		SellPrice:   req.SellPrice,
	}
	if err := pp.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, pp); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateProductPrice(ctx context.Context, repo repository.ProductPriceRepository, id string, req *dto.UpdateProductPriceRequest, uow uow.UnitOfWork) error {
	pp, err := repo.FindByID(ctx, id)
	if err != nil || pp == nil {
		return err
	}
	if req.SellPrice != nil {
		pp.SellPrice = *req.SellPrice
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, pp); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeCreateProductUOMConversion(ctx context.Context, repo repository.ProductUOMConversionRepository, req *dto.CreateProductUOMConversionRequest, uow uow.UnitOfWork) error {
	puc := &entity.ProductUOMConversion{
		ProductID:      req.ProductID,
		UOMID:          req.UOMID,
		ConversionRate: req.ConversionRate,
		Barcode:        req.Barcode,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		Weight:         req.Weight,
		IsStackable:    req.IsStackable,
		MaxStackLayer:  req.MaxStackLayer,
	}
	if err := puc.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, puc); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateProductUOMConversion(ctx context.Context, repo repository.ProductUOMConversionRepository, id string, req *dto.UpdateProductUOMConversionRequest, uow uow.UnitOfWork) error {
	puc, err := repo.FindByID(ctx, id)
	if err != nil || puc == nil {
		return err
	}
	if req.ConversionRate != nil {
		puc.ConversionRate = *req.ConversionRate
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, puc); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeCreateSupplier(ctx context.Context, repo repository.SupplierRepository, req *dto.CreateSupplierRequest, uow uow.UnitOfWork) error {
	supplier := &entity.Supplier{
		Code:          req.Code,
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		Address:       req.Address,
	}
	if err := supplier.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, supplier); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateSupplier(ctx context.Context, repo repository.SupplierRepository, id string, req *dto.UpdateSupplierRequest, uow uow.UnitOfWork) error {
	supplier, err := repo.FindByID(ctx, id)
	if err != nil || supplier == nil {
		return err
	}
	if req.Name != nil {
		supplier.Name = *req.Name
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, supplier); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeCreateProductSupplier(ctx context.Context, repo repository.ProductSupplierRepository, req *dto.CreateProductSupplierRequest, uow uow.UnitOfWork) error {
	ps := &entity.ProductSupplier{
		ProductID:           req.ProductID,
		SupplierID:          req.SupplierID,
		StoreID:             req.StoreID,
		SupplierSKU:         req.SupplierSKU,
		IsPrimary:           req.IsPrimary,
		IsConsignment:       req.IsConsignment,
		IsReturnable:        req.IsReturnable,
		DefaultLeadTimeDays: req.DefaultLeadTimeDays,
		OfferedPrice:        req.OfferedPrice,
		MinOrderQty:         req.MinOrderQty,
	}
	if err := ps.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, ps); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateProductSupplier(ctx context.Context, repo repository.ProductSupplierRepository, id string, req *dto.UpdateProductSupplierRequest, uow uow.UnitOfWork) error {
	ps, err := repo.FindByID(ctx, id)
	if err != nil || ps == nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, ps); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeCreateChartOfAccount(ctx context.Context, repo repository.ChartOfAccountRepository, req *dto.CreateChartOfAccountRequest, uow uow.UnitOfWork) error {
	coa := &entity.ChartOfAccount{
		AccountCode: req.AccountCode,
		Name:        req.Name,
		AccountType: req.AccountType,
	}
	if err := coa.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, coa); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateChartOfAccount(ctx context.Context, repo repository.ChartOfAccountRepository, id string, req *dto.UpdateChartOfAccountRequest, uow uow.UnitOfWork) error {
	coa, err := repo.FindByID(ctx, id)
	if err != nil || coa == nil {
		return err
	}
	if req.Name != nil {
		coa.Name = *req.Name
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, coa); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeCreateTax(ctx context.Context, repo repository.TaxRepository, req *dto.CreateTaxRequest, uow uow.UnitOfWork) error {
	tax := &entity.Tax{
		Name:           req.Name,
		RatePercentage: req.RatePercentage,
		TaxAccountID:   req.TaxAccountID,
	}
	if err := tax.GenerateID(); err != nil {
		return err
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Create(txCtx, tax); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func executeUpdateTax(ctx context.Context, repo repository.TaxRepository, id string, req *dto.UpdateTaxRequest, uow uow.UnitOfWork) error {
	tax, err := repo.FindByID(ctx, id)
	if err != nil || tax == nil {
		return err
	}
	if req.Name != nil {
		tax.Name = *req.Name
	}
	txCtx, _ := uow.Begin(ctx)
	defer uow.Rollback(txCtx)
	if err := repo.Update(txCtx, tax); err != nil {
		return err
	}
	return uow.Commit(txCtx)
}

func toMasterDataProposalListResponse(p *entity.MasterDataProposal) *dto.MasterDataProposalListResponse {
	resp := &dto.MasterDataProposalListResponse{
		ID:              p.ID,
		ReferenceNumber: p.ReferenceNumber,
		EntityType:      p.EntityType,
		ActionType:      p.ActionType,
		TotalItems:      p.TotalItems,
		Status:          p.Status,
		ProposedByID:    p.ProposedByID,
		ReviewedByID:    p.ReviewedByID,
		ReviewedAt:      p.ReviewedAt,
		Reason:          p.Reason,
		CreatedAt:       p.CreatedAt,
	}
	if p.ProposedBy.Name != "" {
		resp.ProposedByName = p.ProposedBy.Name
	}
	if p.ReviewedBy != nil && p.ReviewedBy.Name != "" {
		name := p.ReviewedBy.Name
		resp.ReviewedByName = &name
	}
	return resp
}

func toMasterDataProposalDetailResponse(p *entity.MasterDataProposal) *dto.MasterDataProposalDetailResponse {
	items := make([]dto.MasterDataProposalItemResponse, len(p.Items))
	for i, item := range p.Items {
		items[i] = dto.MasterDataProposalItemResponse{
			ID:           item.ID,
			SeqNo:        item.SeqNo,
			EntityID:     item.EntityID,
			PayloadJSON:  item.PayloadJSON,
			SnapshotJSON: item.SnapshotJSON,
		}
	}

	return &dto.MasterDataProposalDetailResponse{
		ID:              p.ID,
		ReferenceNumber: p.ReferenceNumber,
		EntityType:      p.EntityType,
		ActionType:      p.ActionType,
		TotalItems:      p.TotalItems,
		Status:          p.Status,
		ProposedByID:    p.ProposedByID,
		ReviewedByID:    p.ReviewedByID,
		ReviewedAt:      p.ReviewedAt,
		Reason:          p.Reason,
		ReviewNotes:     p.ReviewNotes,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		Items:           items,
	}
}

func (u *masterDataProposalUseCaseImpl) GetByGroup(ctx context.Context, groupID string) ([]dto.MasterDataProposalDetailResponse, error) {
	proposals, err := u.repo.FindByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	var resp []dto.MasterDataProposalDetailResponse
	for _, p := range proposals {
		resp = append(resp, *toMasterDataProposalDetailResponse(&p))
	}
	return resp, nil
}

func (u *masterDataProposalUseCaseImpl) BulkLinkProductSupplier(ctx context.Context, userID string, req dto.BulkCreateProductSupplierProposalRequest) (*dto.BulkProposalResponse, error) {
	var refNumbers []string
	var proposalResponses []dto.MasterDataProposalDetailResponse
	successCount := 0
	failedCount := 0

	refNum, err := u.generateReferenceNumber(entity.ProposalEntityProductSupplier, time.Now())
	if err != nil {
		return nil, err
	}

	userUUID, _ := uuid.Parse(userID)

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	for i, item := range req.Items {
		linkDTO := dto.ProductSupplierLinkDetailDTO{
			ProductID:           item.ProductID,
			StoreID:             item.StoreID,
			SupplierSKU:         item.SupplierSKU,
			IsPrimary:           item.IsPrimary,
			IsConsignment:       item.IsConsignment,
			IsReturnable:        item.IsReturnable,
			DefaultLeadTimeDays: item.DefaultLeadTimeDays,
			OfferedPrice:        item.OfferedPrice,
			MinOrderQty:         item.MinOrderQty,
		}
		payloadJSON, err := json.Marshal(linkDTO)
		if err != nil {
			failedCount++
			continue
		}

		proposalItem := entity.MasterDataProposalItem{
			SeqNo:       i + 1,
			PayloadJSON: string(payloadJSON),
		}
		if err := proposalItem.GenerateID(); err != nil {
			failedCount++
			continue
		}

		proposal := &entity.MasterDataProposal{
			ReferenceNumber: refNum,
			EntityType:      entity.ProposalEntityProductSupplier,
			ActionType:      entity.ProposalActionCreate,
			TotalItems:      len(req.Items),
			Reason:          req.Reason,
			Status:          entity.ProposalStatusPending,
			ProposedByID:    userUUID,
			Items:           []entity.MasterDataProposalItem{proposalItem},
		}

		if err := proposal.GenerateID(); err != nil {
			failedCount++
			continue
		}

		proposalItem.ProposalID = proposal.ID
		if err := u.repo.Create(txCtx, proposal); err != nil {
			failedCount++
			continue
		}

		refNumbers = append(refNumbers, proposal.ReferenceNumber)
		proposalResponses = append(proposalResponses, *toMasterDataProposalDetailResponse(proposal))
		successCount++
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return &dto.BulkProposalResponse{
		ReferenceNumbers: refNumbers,
		TotalCount:       len(req.Items),
		SuccessCount:     successCount,
		FailedCount:      failedCount,
		Proposals:        proposalResponses,
	}, nil
}
