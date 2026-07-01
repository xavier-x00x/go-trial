package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrPlanningNotFound = errors.New("planning not found")
)

type PurchaseOrderPlanningUseCase interface {
	Calculate(ctx context.Context, storeID string, date string, forceRecal bool) (*dto.PlanningSummaryResponse, error)
	GetPending(ctx context.Context, storeID string) ([]dto.PurchaseOrderPlanningResponse, error)
	GetAll(ctx context.Context, storeID string, status string) ([]dto.PurchaseOrderPlanningResponse, error)
	GetByID(ctx context.Context, id string) (*dto.PurchaseOrderPlanningResponse, error)
	ApprovePlanning(ctx context.Context, req dto.ApprovePlanningRequest) ([]dto.PurchaseOrderPlanningResponse, error)
	IgnorePlanning(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, req dto.UpdatePlanningRequest) error
	BulkSelect(ctx context.Context, req dto.BulkSelectPlanningRequest) error
}

type PurchaseOrderPlanningConfig struct {
	PlanningRepo               repository.PurchaseOrderPlanningRepository
	StoreProductAssortmentRepo repository.StoreProductAssortmentRepository
	StoreRepo                  repository.StoreRepository
	Uow                        interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
	Config struct {
		ZScore            float64 // Service level, default 1.65 for 95%
		SalesLookbackDays int     // Days to look back for sales calculation (unused currently)
		MinLeadTime       int     // Minimum lead time in days (unused currently)
		DefaultLeadTime   int     // Default lead time if no history
	}
}

type purchaseOrderPlanningUseCaseImpl struct {
	repo                       repository.PurchaseOrderPlanningRepository
	storeProductAssortmentRepo repository.StoreProductAssortmentRepository
	storeRepo                  repository.StoreRepository
	uow                        interface {
		Begin(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
	zScore            float64
	salesLookbackDays int
	minLeadTime       int
	defaultLeadTime   int
}

func NewPurchaseOrderPlanningUseCase(cfg PurchaseOrderPlanningConfig) PurchaseOrderPlanningUseCase {
	zScore := cfg.Config.ZScore
	if zScore <= 0 {
		zScore = 1.65 // 95% service level
	}
	lookback := cfg.Config.SalesLookbackDays
	if lookback <= 0 {
		lookback = 30
	}
	minLeadTime := cfg.Config.MinLeadTime
	if minLeadTime <= 0 {
		minLeadTime = 1
	}
	defaultLeadTime := cfg.Config.DefaultLeadTime
	if defaultLeadTime <= 0 {
		defaultLeadTime = 7
	}

	return &purchaseOrderPlanningUseCaseImpl{
		repo:                       cfg.PlanningRepo,
		storeProductAssortmentRepo: cfg.StoreProductAssortmentRepo,
		storeRepo:                  cfg.StoreRepo,
		uow:                        cfg.Uow,
		zScore:                     zScore,
		salesLookbackDays:          lookback,
		minLeadTime:                minLeadTime,
		defaultLeadTime:            defaultLeadTime,
	}
}

func (u *purchaseOrderPlanningUseCaseImpl) Calculate(ctx context.Context, storeID string, date string, forceRecal bool) (*dto.PlanningSummaryResponse, error) {
	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = u.uow.Rollback(txCtx)
		}
	}()

	storeUUID, err := uuid.Parse(storeID)
	if err != nil {
		return nil, err
	}

	calcDate := time.Now()
	if date != "" {
		calcDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			calcDate = time.Now()
		}
	}

	// Delete existing planning for this date if force recal
	if forceRecal {
		u.repo.DeleteByDate(txCtx, storeID, calcDate.Format("2006-01-02"))
	}

	// Get store data for reference
	store, err := u.storeRepo.FindByID(txCtx, storeID)
	if err != nil || store == nil {
		return nil, fmt.Errorf("store not found")
	}
	storeUUID = store.ID

	// Get all planning data with single JOIN query
	planningDataList, err := u.storeProductAssortmentRepo.FindForPlanning(txCtx, storeID)
	if err != nil {
		return nil, err
	}

	// Build plannings first
	var plannings []entity.PurchaseOrderPlanning
	for _, pd := range planningDataList {
		avgDailySales := pd.AverageDailySales.InexactFloat64()
		if avgDailySales <= 0 {
			avgDailySales = 1
		}

		currentStock := pd.CurrentStock.InexactFloat64()
		maxStockQty := pd.MaxStockQty.InexactFloat64()

		leadTime := pd.DefaultLeadTimeDays
		if leadTime <= 0 {
			leadTime = u.defaultLeadTime
		}

		leadTimeDemand := avgDailySales * float64(leadTime)

		salesVariance := avgDailySales * 0.5
		dynamicSafetyStock := u.zScore * math.Sqrt(float64(leadTime)*salesVariance)
		if dynamicSafetyStock < 1 {
			dynamicSafetyStock = avgDailySales
		}

		staticSafetyStock := pd.SafetyStockQty.InexactFloat64()
		if staticSafetyStock <= 0 {
			staticSafetyStock = dynamicSafetyStock
		}

		rop := leadTimeDemand + dynamicSafetyStock
		recommendedOrderQty := calculateRecommendedOrder(currentStock, rop, avgDailySales, leadTime, dynamicSafetyStock, maxStockQty)

		if currentStock <= rop {
			planning := entity.PurchaseOrderPlanning{
				StoreID:             storeUUID,
				ProductID:           pd.ProductID,
				ProductSupplierID:   pd.ProductSupplierID,
				CurrentStock:        decimal.NewFromFloat(currentStock),
				SafetyStock:         decimal.NewFromFloat(staticSafetyStock),
				DynamicSafetyStock:  decimal.NewFromFloat(dynamicSafetyStock),
				MaxStockQty:         decimal.NewFromFloat(maxStockQty),
				ReorderPoint:        decimal.NewFromFloat(rop),
				AverageDailySales:   decimal.NewFromFloat(avgDailySales),
				LeadTimeDays:        leadTime,
				LeadTimeDemand:      decimal.NewFromFloat(leadTimeDemand),
				Status:              entity.PlanningStatusPending,
				RecommendedOrderQty: decimal.NewFromFloat(recommendedOrderQty),
				OrderQty:            decimal.NewFromFloat(recommendedOrderQty),
				CalculatedDate:      calcDate,
			}
			planning.GenerateID()
			plannings = append(plannings, planning)
		}
	}

	// Bulk insert with transaction
	if err := u.repo.CreateBatch(txCtx, plannings); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	return &dto.PlanningSummaryResponse{
		ProcessedCount:  len(plannings),
		NeedsOrderCount: len(plannings),
		CalculatedDate:  calcDate,
	}, nil
}

func calculateRecommendedOrder(currentStock float64, rop float64, avgDailySales float64, leadTime int, safetyStock float64, maxStockQty float64) float64 {
	if currentStock > rop {
		return 0
	}

	// Base target = (Lead Time × Daily Sales × 2) + Safety Stock (approximately 2 weeks + safety)
	targetStock := (avgDailySales * float64(leadTime) * 2) + safetyStock

	// If max_stock_qty is set and less than target, use max_stock_qty
	if maxStockQty > 0 && maxStockQty < targetStock {
		targetStock = maxStockQty
	}

	recommended := targetStock - currentStock
	if recommended < 0 {
		recommended = 0
	}

	// Round to nearest whole number
	return math.Round(recommended)
}

func (u *purchaseOrderPlanningUseCaseImpl) GetPending(ctx context.Context, storeID string) ([]dto.PurchaseOrderPlanningResponse, error) {
	plans, err := u.repo.FindPendingByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return toPlanningResponses(plans), nil
}

func (u *purchaseOrderPlanningUseCaseImpl) GetAll(ctx context.Context, storeID string, status string) ([]dto.PurchaseOrderPlanningResponse, error) {
	plans, err := u.repo.FindByStoreID(ctx, storeID, status)
	if err != nil {
		return nil, err
	}

	return toPlanningResponses(plans), nil
}

func (u *purchaseOrderPlanningUseCaseImpl) GetByID(ctx context.Context, id string) (*dto.PurchaseOrderPlanningResponse, error) {
	plan, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := toPlanningResponse(plan)
	return &resp, nil
}

func (u *purchaseOrderPlanningUseCaseImpl) ApprovePlanning(ctx context.Context, req dto.ApprovePlanningRequest) ([]dto.PurchaseOrderPlanningResponse, error) {
	var results []dto.PurchaseOrderPlanningResponse

	for i, productID := range req.ProductIDs {
		plan, err := u.repo.FindByID(ctx, productID.String())
		if err != nil {
			continue
		}

		now := time.Now()
		plan.Status = entity.PlanningStatusApproved
		plan.ProcessedDate = &now
		plan.ProcessedByID = &req.ProcessedByID

		if i < len(req.OrderQuantities) && req.OrderQuantities[i] > 0 {
			plan.OrderQty = decimal.NewFromFloat(req.OrderQuantities[i])
		}

		if i < len(req.ProductSupplierIDs) && req.ProductSupplierIDs[i] != "" {
			if psID, err := uuid.Parse(req.ProductSupplierIDs[i]); err == nil && plan.ProductSupplierID != psID {
				plan.ProductSupplierID = psID
				plan.IsManualSupplier = true
			}
		}

		if err := u.repo.Update(ctx, plan); err != nil {
			continue
		}

		results = append(results, toPlanningResponse(plan))
	}

	return results, nil
}

func (u *purchaseOrderPlanningUseCaseImpl) IgnorePlanning(ctx context.Context, ids []string) error {
	for _, id := range ids {
		plan, err := u.repo.FindByID(ctx, id)
		if err != nil {
			continue
		}

		now := time.Now()
		plan.Status = entity.PlanningStatusIgnored
		plan.ProcessedDate = &now

		u.repo.Update(ctx, plan)
	}

	return nil
}

func (u *purchaseOrderPlanningUseCaseImpl) Update(ctx context.Context, id string, req dto.UpdatePlanningRequest) error {
	plan, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if plan == nil {
		return errors.New("planning not found")
	}

	if req.OrderQty != nil {
		plan.OrderQty = decimal.NewFromFloat(*req.OrderQty)
	}

	if req.ProductSupplierID != nil && *req.ProductSupplierID != "" {
		if psID, err := uuid.Parse(*req.ProductSupplierID); err == nil && plan.ProductSupplierID != psID {
			plan.ProductSupplierID = psID
			plan.IsManualSupplier = true
		}
	}

	if req.IsSelected != nil {
		plan.IsSelected = *req.IsSelected
	}

	return u.repo.Update(ctx, plan)
}

func (u *purchaseOrderPlanningUseCaseImpl) BulkSelect(ctx context.Context, req dto.BulkSelectPlanningRequest) error {
	for _, id := range req.IDs {
		plan, err := u.repo.FindByID(ctx, id.String())
		if err == nil && plan != nil {
			plan.IsSelected = req.IsSelected
			u.repo.Update(ctx, plan)
		}
	}
	return nil
}

func toPlanningResponses(plans []entity.PurchaseOrderPlanning) []dto.PurchaseOrderPlanningResponse {
	responses := make([]dto.PurchaseOrderPlanningResponse, len(plans))
	for i, plan := range plans {
		responses[i] = toPlanningResponse(&plan)
	}
	return responses
}

func toPlanningResponse(plan *entity.PurchaseOrderPlanning) dto.PurchaseOrderPlanningResponse {
	return dto.PurchaseOrderPlanningResponse{
		ID:                  plan.ID,
		StoreID:             plan.StoreID,
		ProductID:           plan.ProductID,
		ProductSupplierID:   plan.ProductSupplierID,
		CurrentStock:        plan.CurrentStock.InexactFloat64(),
		SafetyStock:         plan.SafetyStock.InexactFloat64(),
		DynamicSafetyStock:  plan.DynamicSafetyStock.InexactFloat64(),
		ReorderPoint:        plan.ReorderPoint.InexactFloat64(),
		AverageDailySales:   plan.AverageDailySales.InexactFloat64(),
		LeadTimeDays:        plan.LeadTimeDays,
		LeadTimeDemand:      plan.LeadTimeDemand.InexactFloat64(),
		Status:              plan.Status,
		RecommendedOrderQty: plan.RecommendedOrderQty.InexactFloat64(),
		OrderQty:            plan.OrderQty.InexactFloat64(),
		IsManualSupplier:    plan.IsManualSupplier,
		IsSelected:          plan.IsSelected,
		CalculatedDate:      plan.CalculatedDate,
		ProcessedDate:       plan.ProcessedDate,
		ProcessedByID:       plan.ProcessedByID,
	}
}
