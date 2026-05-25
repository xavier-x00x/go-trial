package handler

import (
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/jwt"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PurchaseOrderHandler struct {
	uc usecase.PurchaseOrderUseCase
	v  *validator.Validator
}

func NewPurchaseOrderHandler(uc usecase.PurchaseOrderUseCase, v *validator.Validator) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{uc: uc, v: v}
}

func getUserIDFromPOContext(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*jwt.Claims)
	if !ok {
		return ""
	}
	return claims.UserID
}

func (h *PurchaseOrderHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePurchaseOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromPOContext(c)
	result, err := h.uc.Create(c.UserContext(), userID, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Purchase order created successfully", result)
}

func (h *PurchaseOrderHandler) CreateFromPlanning(c *fiber.Ctx) error {
	var req dto.CreatePOFromPlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromPOContext(c)
	result, err := h.uc.CreateFromPlanning(c.UserContext(), userID, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Purchase order(s) created from planning", result)
}

func (h *PurchaseOrderHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Purchase order retrieved successfully", result)
}

func (h *PurchaseOrderHandler) GetByStoreID(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	status := c.Query("status")
	result, err := h.uc.GetByStoreID(c.UserContext(), storeID, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Purchase orders retrieved successfully", result)
}

func (h *PurchaseOrderHandler) GetPendingByStoreID(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	result, err := h.uc.GetPendingByStoreID(c.UserContext(), storeID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Pending purchase orders retrieved successfully", result)
}

func (h *PurchaseOrderHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdatePurchaseOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromPOContext(c)
	result, err := h.uc.Update(c.UserContext(), userID, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Purchase order updated successfully", result)
}

func (h *PurchaseOrderHandler) Submit(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := getUserIDFromPOContext(c)
	if err := h.uc.Submit(c.UserContext(), userID, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Purchase order submitted successfully", nil)
}

func (h *PurchaseOrderHandler) Approve(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := getUserIDFromPOContext(c)
	if err := h.uc.Approve(c.UserContext(), userID, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Purchase order approved successfully", nil)
}

func (h *PurchaseOrderHandler) Cancel(c *fiber.Ctx) error {
	var req dto.CancelPORequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userID := getUserIDFromPOContext(c)
	if err := h.uc.Cancel(c.UserContext(), userID, req.ID.String(), req.Reason); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Purchase order cancelled successfully", nil)
}

func (h *PurchaseOrderHandler) BulkCreateFromApprovedPlanning(c *fiber.Ctx) error {
	var req dto.BulkCreatePOFromPlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromPOContext(c)
	result, err := h.uc.BulkCreateFromApprovedPlanning(c.UserContext(), userID, req.StoreID.String(), req.WarehouseID.String())
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Purchase orders created from approved planning", result)
}

func (h *PurchaseOrderHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "created_at"),
		OrderDir:    c.Query("order_dir", "desc"),
		Conditions:  parseConditions(c),
	}

	data, meta, err := h.uc.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Purchase orders retrieved successfully", data, meta)
}

func parseConditions(c *fiber.Ctx) map[string]interface{} {
	conditions := make(map[string]interface{})
	if storeID := c.Query("store_id"); storeID != "" {
		conditions["store_id"] = storeID
	}
	if status := c.Query("status"); status != "" {
		conditions["status"] = status
	}
	return conditions
}