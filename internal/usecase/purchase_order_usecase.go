package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrPurchaseOrderNotFound = errors.New("purchase order not found")
	ErrPOInvalidStatus       = errors.New("invalid purchase order status transition")
	ErrPOItemNotFound        = errors.New("po item not found")
	ErrPONotificationFailed  = errors.New("can only resend failed notifications")
)

type PurchaseOrderUseCase interface {
	Create(ctx context.Context, userID string, req dto.CreatePurchaseOrderRequest) (*dto.PurchaseOrderDetailResponse, error)
	CreateFromPlanning(ctx context.Context, userID string, req dto.CreatePOFromPlanningRequest) ([]dto.PurchaseOrderDetailResponse, error)
	BulkCreateFromApprovedPlanning(ctx context.Context, userID string, storeID string, warehouseID string) ([]dto.PurchaseOrderDetailResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PurchaseOrderDetailResponse, error)
	GetByPONumber(ctx context.Context, poNum string) (*dto.PurchaseOrderDetailResponse, error)
	GetByStoreID(ctx context.Context, storeID string, status string) ([]dto.PurchaseOrderListResponse, error)
	GetPendingByStoreID(ctx context.Context, storeID string) ([]dto.PurchaseOrderListResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PurchaseOrderListResponse, *entity.Meta, error)
	Submit(ctx context.Context, userID string, id string) error
	Approve(ctx context.Context, userID string, id string) error
	Cancel(ctx context.Context, userID string, id string, reason string) error
	Resend(ctx context.Context, id string) error
}

type PurchaseOrderConfig struct {
	Repo                repository.PurchaseOrderRepository
	PlanningRepo        repository.PurchaseOrderPlanningRepository
	ProductSupplierRepo repository.ProductSupplierRepository
	ProductRepo         repository.ProductRepository
	SupplierRepo        repository.SupplierRepository
	StoreRepo           repository.StoreRepository
	WarehouseRepo       repository.WarehouseRepository
	UserRepo            repository.UserRepository
	NumberSequenceRepo  repository.NumberSequenceRepository
	NotificationService NotificationService
	Uow                 interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

type purchaseOrderUseCaseImpl struct {
	repo                repository.PurchaseOrderRepository
	planningRepo        repository.PurchaseOrderPlanningRepository
	productSupplierRepo repository.ProductSupplierRepository
	productRepo         repository.ProductRepository
	supplierRepo        repository.SupplierRepository
	storeRepo           repository.StoreRepository
	warehouseRepo       repository.WarehouseRepository
	userRepo            repository.UserRepository
	numberSequenceRepo  repository.NumberSequenceRepository
	notificationService NotificationService
	uow                 interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
}

func NewPurchaseOrderUseCase(cfg PurchaseOrderConfig) PurchaseOrderUseCase {
	return &purchaseOrderUseCaseImpl{
		repo:                cfg.Repo,
		planningRepo:        cfg.PlanningRepo,
		productSupplierRepo: cfg.ProductSupplierRepo,
		productRepo:         cfg.ProductRepo,
		supplierRepo:        cfg.SupplierRepo,
		storeRepo:           cfg.StoreRepo,
		warehouseRepo:       cfg.WarehouseRepo,
		userRepo:            cfg.UserRepo,
		numberSequenceRepo:  cfg.NumberSequenceRepo,
		notificationService: cfg.NotificationService,
		uow:                 cfg.Uow,
	}
}

func (u *purchaseOrderUseCaseImpl) generatePONumber(date time.Time) (string, error) {
	prefix := "PO"
	period := date.Format("0601") // YYMM

	seqNum, err := u.numberSequenceRepo.GetNextNumber(context.Background(), prefix, period)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%05d", prefix, period, seqNum), nil
}

func (u *purchaseOrderUseCaseImpl) Create(ctx context.Context, userID string, req dto.CreatePurchaseOrderRequest) (*dto.PurchaseOrderDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	if len(req.Items) == 0 {
		return nil, errors.New("items is required")
	}

	supplier, err := u.supplierRepo.FindByID(ctx, req.SupplierID.String())
	if err != nil || supplier == nil {
		return nil, errors.New("supplier not found")
	}
	if err != nil || supplier == nil {
		return nil, errors.New("supplier not found")
	}

	store, err := u.storeRepo.FindByID(ctx, req.StoreID.String())
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	warehouse, err := u.warehouseRepo.FindByID(ctx, req.WarehouseID.String())
	if err != nil || warehouse == nil {
		return nil, errors.New("warehouse not found")
	}

	creator, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || creator == nil {
		return nil, errors.New("creator user not found")
	}

	orderDate := req.OrderDate
	if orderDate.IsZero() {
		orderDate = time.Now()
	}

	poNum, err := u.generatePONumber(orderDate)
	if err != nil {
		return nil, err
	}

	totalAmount := decimal.Zero

	items := make([]entity.PurchaseOrderItem, len(req.Items))
	for i, item := range req.Items {
		qty := decimal.NewFromFloat(item.QtyOrdered)
		unitPrice := decimal.NewFromFloat(item.UnitPrice)
		lineSubtotal := qty.Mul(unitPrice)

		totalAmount = totalAmount.Add(lineSubtotal)

		items[i] = entity.PurchaseOrderItem{
			SeqNo:             i + 1,
			ProductID:         item.ProductID,
			UOMID:             item.UOMID,
			QtyOrdered:        qty,
			UnitPrice:         unitPrice,
			Subtotal:          lineSubtotal,
			ProductSupplierID: item.ProductSupplierID,
			PlanningID:        item.PlanningID,
			Notes:             item.Notes,
		}
	}

	po := &entity.PurchaseOrder{
		PONumber:         poNum,
		ReferenceNo:      req.ReferenceNo,
		SupplierID:       req.SupplierID,
		StoreID:          req.StoreID,
		WarehouseID:      req.WarehouseID,
		OrderDate:        orderDate,
		ExpectedDelivery: req.ExpectedDelivery,
		PaymentTermDays:  req.PaymentTermDays,
		PaymentMode:      req.PaymentMode,
		TotalAmount:      totalAmount,
		Status:           entity.POStatusDraft,
		CreatedByID:      userUUID,
		Notes:            req.Notes,
		SupplierNotes:    req.SupplierNotes,
		Items:            items,
		SupplierName:     supplier.Name,
		SupplierCode:     supplier.Code,
		SupplierAddress:  supplier.Address,
		StoreCode:        store.Code,
		StoreName:        store.Name,
		StoreAddress:     store.Address,
		WarehouseName:    warehouse.Name,
		CreatedByName:    creator.Name,
	}

	if err := po.GenerateID(); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].PurchaseOrderID = po.ID
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Create(txCtx, po); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return toPODetailResponse(po), nil
}

func (u *purchaseOrderUseCaseImpl) CreateFromPlanning(ctx context.Context, userID string, req dto.CreatePOFromPlanningRequest) ([]dto.PurchaseOrderDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	if len(req.PlanningIDs) == 0 {
		return nil, errors.New("planning_ids is required")
	}

	poNum, err := u.generatePONumber(time.Now())
	if err != nil {
		return nil, err
	}

	orderDate := req.OrderDate
	if orderDate.IsZero() {
		orderDate = time.Now()
	}

	var items []entity.PurchaseOrderItem
	var consignmentItems []entity.PurchaseOrderItem
	var totalSubtotal decimal.Decimal
	var consignmentSubtotal decimal.Decimal
	seqCount := 0
	consignmentSeqCount := 0
	supplierID := uuid.Nil

	for _, planningID := range req.PlanningIDs {
		plan, err := u.planningRepo.FindByID(ctx, planningID)
		if err != nil {
			continue
		}

		if plan.Status != entity.PlanningStatusApproved {
			continue
		}

		ps, err := u.productSupplierRepo.FindByID(ctx, plan.ProductSupplierID.String())
		if err != nil || ps == nil {
			continue
		}

		prod, err := u.productRepo.FindByID(ctx, plan.ProductID.String())
		if err != nil || prod == nil {
			continue
		}

		if supplierID == uuid.Nil {
			supplierID = ps.SupplierID
		}

		qty := plan.RecommendedOrderQty
		unitPrice := ps.OfferedPrice
		lineSubtotal := qty.Mul(unitPrice)

		if !ps.IsConsignment {
			seqCount++
			items = append(items, entity.PurchaseOrderItem{
				SeqNo:             seqCount,
				ProductID:         plan.ProductID,
				UOMID:             prod.BaseUOMID,
				QtyOrdered:        qty,
				UnitPrice:         unitPrice,
				Subtotal:          lineSubtotal,
				ProductSupplierID: &plan.ProductSupplierID,
				PlanningID:        &plan.ID,
			})
			totalSubtotal = totalSubtotal.Add(lineSubtotal)
		} else {
			consignmentSeqCount++
			consignmentItems = append(consignmentItems, entity.PurchaseOrderItem{
				SeqNo:             consignmentSeqCount,
				ProductID:         plan.ProductID,
				UOMID:             prod.BaseUOMID,
				QtyOrdered:        qty,
				UnitPrice:         unitPrice,
				Subtotal:          lineSubtotal,
				ProductSupplierID: &plan.ProductSupplierID,
				PlanningID:        &plan.ID,
			})
			consignmentSubtotal = consignmentSubtotal.Add(lineSubtotal)
		}
	}

	if len(items) == 0 && len(consignmentItems) == 0 {
		return nil, errors.New("no valid planning items found")
	}

	store, err := u.storeRepo.FindByID(ctx, req.StoreID.String())
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	warehouse, err := u.warehouseRepo.FindByID(ctx, req.WarehouseID.String())
	if err != nil || warehouse == nil {
		return nil, errors.New("warehouse not found")
	}

	supplier, err := u.supplierRepo.FindByID(ctx, supplierID.String())
	if err != nil || supplier == nil {
		return nil, errors.New("supplier not found")
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	var poResponses []dto.PurchaseOrderDetailResponse

	if len(items) > 0 {
		po := &entity.PurchaseOrder{
			PONumber:        poNum,
			SupplierID:      supplierID,
			StoreID:         req.StoreID,
			WarehouseID:     req.WarehouseID,
			OrderDate:       orderDate,
			TotalAmount:     totalSubtotal,
			Status:          entity.POStatusDraft,
			CreatedByID:     userUUID,
			Notes:           req.Notes,
			Items:           items,
			SupplierName:    supplier.Name,
			SupplierCode:    supplier.Code,
			SupplierAddress: supplier.Address,
			StoreCode:       store.Code,
			StoreName:       store.Name,
			StoreAddress:    store.Address,
			WarehouseName:   warehouse.Name,
		}

		if err := po.GenerateID(); err != nil {
			return nil, err
		}

		for i := range items {
			items[i].PurchaseOrderID = po.ID
		}

		if err := u.repo.Create(txCtx, po); err != nil {
			return nil, err
		}

		poResponses = append(poResponses, *toPODetailResponse(po))
	}

	if len(consignmentItems) > 0 {
		poNumConsignment, err := u.generatePONumber(time.Now())
		if err != nil {
			return nil, err
		}

		poConsignment := &entity.PurchaseOrder{
			PONumber:        poNumConsignment,
			SupplierID:      supplierID,
			StoreID:         req.StoreID,
			WarehouseID:     req.WarehouseID,
			OrderDate:       orderDate,
			TotalAmount:     consignmentSubtotal,
			Status:          entity.POStatusDraft,
			CreatedByID:     userUUID,
			Notes:           req.Notes,
			Items:           consignmentItems,
			SupplierName:    supplier.Name,
			SupplierCode:    supplier.Code,
			SupplierAddress: supplier.Address,
			StoreCode:       store.Code,
			StoreName:       store.Name,
			StoreAddress:    store.Address,
			WarehouseName:   warehouse.Name,
		}

		if err := poConsignment.GenerateID(); err != nil {
			return nil, err
		}

		for i := range consignmentItems {
			consignmentItems[i].PurchaseOrderID = poConsignment.ID
		}

		if err := u.repo.Create(txCtx, poConsignment); err != nil {
			return nil, err
		}

		poResponses = append(poResponses, *toPODetailResponse(poConsignment))
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return poResponses, nil
}

func (u *purchaseOrderUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.PurchaseOrderDetailResponse, error) {
	po, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if po == nil {
		return nil, ErrPurchaseOrderNotFound
	}

	return toPODetailResponse(po), nil
}

func (u *purchaseOrderUseCaseImpl) BulkCreateFromApprovedPlanning(ctx context.Context, userID string, storeID string, warehouseID string) ([]dto.PurchaseOrderDetailResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	storeUUID, err := uuid.Parse(storeID)
	if err != nil {
		return nil, err
	}

	warehouseUUID, err := uuid.Parse(warehouseID)
	if err != nil {
		return nil, err
	}

	approvedPlannings, err := u.planningRepo.FindByStoreID(ctx, storeID, entity.PlanningStatusApproved)
	if err != nil {
		return nil, err
	}

	if len(approvedPlannings) == 0 {
		return nil, errors.New("no approved planning found for this store")
	}

	type supplierGroup struct {
		SupplierID    uuid.UUID
		IsConsignment bool
		Items         []entity.PurchaseOrderItem
		Subtotal      decimal.Decimal
		TaxAmount     decimal.Decimal
	}

	groups := make(map[string]*supplierGroup)
	seqCount := make(map[string]int)

	for _, plan := range approvedPlannings {
		ps, err := u.productSupplierRepo.FindByID(ctx, plan.ProductSupplierID.String())
		if err != nil || ps == nil {
			continue
		}

		prod, err := u.productRepo.FindByID(ctx, plan.ProductID.String())
		if err != nil || prod == nil {
			continue
		}

		consignmentKey := "regular"
		if ps.IsConsignment {
			consignmentKey = "consignment"
		}

		supplierKey := ps.SupplierID.String() + ":" + consignmentKey

		if _, exists := groups[supplierKey]; !exists {
			groups[supplierKey] = &supplierGroup{
				SupplierID:    ps.SupplierID,
				IsConsignment: ps.IsConsignment,
			}
			seqCount[supplierKey] = 0
		}

		seqCount[supplierKey]++

		qty := plan.RecommendedOrderQty
		unitPrice := ps.OfferedPrice
		lineSubtotal := qty.Mul(unitPrice)

		item := entity.PurchaseOrderItem{
			SeqNo:             seqCount[supplierKey],
			ProductID:         plan.ProductID,
			UOMID:             prod.BaseUOMID,
			QtyOrdered:        qty,
			UnitPrice:         unitPrice,
			Subtotal:          lineSubtotal,
			ProductSupplierID: &plan.ProductSupplierID,
			PlanningID:        &plan.ID,
		}

		groups[supplierKey].Items = append(groups[supplierKey].Items, item)
		groups[supplierKey].Subtotal = groups[supplierKey].Subtotal.Add(lineSubtotal)
	}

	if len(groups) == 0 {
		return nil, errors.New("no valid planning items found")
	}

	store, err := u.storeRepo.FindByID(ctx, storeID)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	warehouse, err := u.warehouseRepo.FindByID(ctx, warehouseID)
	if err != nil || warehouse == nil {
		return nil, errors.New("warehouse not found")
	}

	creator, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || creator == nil {
		return nil, errors.New("creator user not found")
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	var poResponses []dto.PurchaseOrderDetailResponse

	for _, group := range groups {
		poNum, err := u.generatePONumber(time.Now())
		if err != nil {
			return nil, err
		}

		supplier, err := u.supplierRepo.FindByID(txCtx, group.SupplierID.String())
		if err != nil || supplier == nil {
			return nil, errors.New("supplier not found")
		}

		po := &entity.PurchaseOrder{
			PONumber:        poNum,
			SupplierID:      group.SupplierID,
			StoreID:         storeUUID,
			WarehouseID:     warehouseUUID,
			OrderDate:       time.Now(),
			TotalAmount:     group.Subtotal,
			Status:          entity.POStatusDraft,
			CreatedByID:     userUUID,
			SupplierName:    supplier.Name,
			SupplierCode:    supplier.Code,
			SupplierAddress: supplier.Address,
			StoreCode:       store.Code,
			StoreName:       store.Name,
			StoreAddress:    store.Address,
			WarehouseName:   warehouse.Name,
			CreatedByName:   creator.Name,
		}

		if err := po.GenerateID(); err != nil {
			return nil, err
		}

		for i := range group.Items {
			group.Items[i].PurchaseOrderID = po.ID
		}

		po.Items = group.Items

		if err := u.repo.Create(txCtx, po); err != nil {
			return nil, err
		}

		poResponses = append(poResponses, *toPODetailResponse(po))
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return poResponses, nil
}

func (u *purchaseOrderUseCaseImpl) GetByPONumber(ctx context.Context, poNum string) (*dto.PurchaseOrderDetailResponse, error) {
	po, err := u.repo.FindByPONumber(ctx, poNum)
	if err != nil {
		return nil, err
	}
	if po == nil {
		return nil, ErrPurchaseOrderNotFound
	}

	return toPODetailResponse(po), nil
}

func (u *purchaseOrderUseCaseImpl) GetByStoreID(ctx context.Context, storeID string, status string) ([]dto.PurchaseOrderListResponse, error) {
	pos, err := u.repo.FindByStoreID(ctx, storeID, status)
	if err != nil {
		return nil, err
	}

	return toPOListResponses(pos), nil
}

func (u *purchaseOrderUseCaseImpl) GetPendingByStoreID(ctx context.Context, storeID string) ([]dto.PurchaseOrderListResponse, error) {
	pos, err := u.repo.FindPendingByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return toPOListResponses(pos), nil
}

func (u *purchaseOrderUseCaseImpl) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PurchaseOrderListResponse, *entity.Meta, error) {
	allowedOrder := []string{"created_at", "updated_at", "po_number", "order_date", "status", "grand_total"}
	searchColumns := []string{"po_number", "reference_no"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	if meta.Conditions != nil {
		filter.Conditions = meta.Conditions
	}

	data, resMeta, err := u.repo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	return toPOListResponses(data), resMeta, nil
}

func (u *purchaseOrderUseCaseImpl) Submit(ctx context.Context, userID string, id string) error {
	po, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if po == nil {
		return ErrPurchaseOrderNotFound
	}

	if po.Status != entity.POStatusDraft {
		return ErrPOInvalidStatus
	}

	po.Status = entity.POStatusSubmitted

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, po); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseOrderUseCaseImpl) Approve(ctx context.Context, userID string, id string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	po, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if po == nil {
		return ErrPurchaseOrderNotFound
	}

	if po.Status != entity.POStatusSubmitted {
		return ErrPOInvalidStatus
	}

	approver, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || approver == nil {
		return errors.New("approver user not found")
	}

	po.Status = entity.POStatusApproved
	po.ApprovedByID = &userUUID
	now := time.Now()
	po.ApprovedAt = &now
	po.ApprovedByName = &approver.Name

	po.NotificationStatus = entity.NotificationStatusPending

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, po); err != nil {
		return err
	}

	if err := u.notificationService.QueueNotification(txCtx, po.ID.String()); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseOrderUseCaseImpl) Cancel(ctx context.Context, userID string, id string, reason string) error {
	po, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if po == nil {
		return ErrPurchaseOrderNotFound
	}

	if po.Status == entity.POStatusReceived || po.Status == entity.POStatusClosed {
		return ErrPOInvalidStatus
	}

	po.Status = entity.POStatusCancelled

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, po); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func (u *purchaseOrderUseCaseImpl) Resend(ctx context.Context, id string) error {
	po, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if po == nil {
		return ErrPurchaseOrderNotFound
	}

	if po.NotificationStatus == entity.NotificationStatusNone {
		return nil
	}
	if po.NotificationStatus != entity.NotificationStatusFailed {
		return errors.New("can only resend failed notifications")
	}

	po.NotificationStatus = entity.NotificationStatusPending

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	if err := u.repo.Update(txCtx, po); err != nil {
		return err
	}

	if err := u.notificationService.QueueNotification(txCtx, po.ID.String()); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toPOListResponses(pos []entity.PurchaseOrder) []dto.PurchaseOrderListResponse {
	responses := make([]dto.PurchaseOrderListResponse, len(pos))
	for i, po := range pos {
		responses[i] = dto.PurchaseOrderListResponse{
			ID:               po.ID,
			PONumber:         po.PONumber,
			ReferenceNo:      po.ReferenceNo,
			SupplierID:       po.SupplierID,
			SupplierName:     po.Supplier.Name,
			SupplierCode:     po.Supplier.Code,
			StoreID:          po.StoreID,
			StoreName:        po.Store.Name,
			WarehouseID:      po.WarehouseID,
			WarehouseName:    po.Warehouse.Name,
			OrderDate:        po.OrderDate,
			ExpectedDelivery: po.ExpectedDelivery,
			TotalAmount:      po.TotalAmount.InexactFloat64(),
			Status:           po.Status,
			ApprovedByID:     po.ApprovedByID,
			ApprovedAt:       po.ApprovedAt,
			CreatedByID:      po.CreatedByID,
			CreatedByName:    po.CreatedBy.Name,
			CreatedAt:        po.CreatedAt,
			UpdatedAt:        po.UpdatedAt,
		}
	}
	return responses
}

func toPODetailResponse(po *entity.PurchaseOrder) *dto.PurchaseOrderDetailResponse {
	items := make([]dto.PurchaseOrderItemResponse, len(po.Items))
	for i, item := range po.Items {
		items[i] = dto.PurchaseOrderItemResponse{
			ID:                item.ID,
			SeqNo:             item.SeqNo,
			ProductID:         item.ProductID,
			ProductSKU:        item.Product.SKU,
			ProductName:       item.Product.Name,
			UOMID:             item.UOMID,
			UOMName:           item.UOM.Name,
			QtyOrdered:        item.QtyOrdered.InexactFloat64(),
			QtyReceived:       item.QtyReceived.InexactFloat64(),
			UnitPrice:         item.UnitPrice.InexactFloat64(),
			Subtotal:          item.Subtotal.InexactFloat64(),
			ProductSupplierID: item.ProductSupplierID,
			PlanningID:        item.PlanningID,
			Notes:             item.Notes,
		}
	}

	return &dto.PurchaseOrderDetailResponse{
		ID:          po.ID,
		PONumber:    po.PONumber,
		ReferenceNo: po.ReferenceNo,
		Supplier: dto.SupplierResponse{
			ID:   po.Supplier.ID,
			Code: po.Supplier.Code,
			Name: po.Supplier.Name,
		},
		Store: dto.StoreResponse{
			ID:   po.StoreID.String(),
			Code: po.StoreCode,
			Name: po.StoreName,
		},
		Warehouse: dto.WarehouseResponse{
			ID:   po.WarehouseID,
			Name: po.WarehouseName,
		},
		OrderDate:               po.OrderDate,
		ExpectedDelivery:        po.ExpectedDelivery,
		PaymentTermDays:         po.PaymentTermDays,
		PaymentMode:             po.PaymentMode,
		TotalAmount:             po.TotalAmount.InexactFloat64(),
		Status:                  po.Status,
		ApprovedByID:            po.ApprovedByID,
		ApprovedAt:              po.ApprovedAt,
		CreatedByID:             po.CreatedByID,
		CreatedByName:           po.CreatedByName,
		Notes:                   po.Notes,
		SupplierNotes:           po.SupplierNotes,
		SupplierNameSnapshot:    po.SupplierName,
		SupplierCodeSnapshot:    po.SupplierCode,
		SupplierAddressSnapshot: po.SupplierAddress,
		StoreCodeSnapshot:       po.StoreCode,
		StoreNameSnapshot:       po.StoreName,
		StoreAddressSnapshot:    po.StoreAddress,
		WarehouseNameSnapshot:   po.WarehouseName,
		ApprovedByName:          po.ApprovedByName,
		CreatedAt:               po.CreatedAt,
		UpdatedAt:               po.UpdatedAt,
		Items:                   items,
	}
}
