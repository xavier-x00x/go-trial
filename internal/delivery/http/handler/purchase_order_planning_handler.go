package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/query/service"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"
	"go-trial/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PurchaseOrderPlanningHandler struct {
	uc usecase.PurchaseOrderPlanningUseCase
	qs *service.PurchaseOrderPlanningQueryService
	v  *validator.Validator
}

func NewPurchaseOrderPlanningHandler(uc usecase.PurchaseOrderPlanningUseCase, qs *service.PurchaseOrderPlanningQueryService, v *validator.Validator) *PurchaseOrderPlanningHandler {
	return &PurchaseOrderPlanningHandler{uc: uc, qs: qs, v: v}
}

func (h *PurchaseOrderPlanningHandler) Calculate(c *fiber.Ctx) error {
	var req dto.CalculatePlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	
	result, err := h.uc.Calculate(c.UserContext(), req.StoreID.String(), req.Date, req.ForceRecal)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Calculation completed", result)
}



func (h *PurchaseOrderPlanningHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdatePlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.uc.Update(c.UserContext(), id, req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Planning updated successfully", nil)
}

func (h *PurchaseOrderPlanningHandler) BulkSelect(c *fiber.Ctx) error {
	var req dto.BulkSelectPlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.v.Validate(&req); len(errs) > 0 {
		return response.Error(c, fiber.StatusBadRequest, errs[0].Message)
	}

	if err := h.uc.BulkSelect(c.UserContext(), req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Plannings selection updated", nil)
}

func (h *PurchaseOrderPlanningHandler) GetPending(c *fiber.Ctx) error {
	storeID := c.Query("store_id")
	search := c.Query("search")
	if storeID == "" {
		return response.Error(c, fiber.StatusBadRequest, "store_id is required")
	}
	
	result, err := h.qs.GetPending(c.UserContext(), storeID, search)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Pending planning retrieved", result)
}

func (h *PurchaseOrderPlanningHandler) GetAll(c *fiber.Ctx) error {
	storeID := c.Query("store_id")
	status := c.Query("status")
	
	result, err := h.qs.GetAll(c.UserContext(), storeID, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Planning list retrieved", result)
}

func (h *PurchaseOrderPlanningHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	
	result, err := h.qs.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if result == nil {
		return response.Error(c, fiber.StatusNotFound, "Planning not found")
	}
	
	return response.Success(c, fiber.StatusOK, "Planning detail retrieved", result)
}

func (h *PurchaseOrderPlanningHandler) Approve(c *fiber.Ctx) error {
	var req dto.ApprovePlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	var userIDStr string
	if claims, ok := c.Locals("claims").(*jwt.Claims); ok {
		userIDStr = claims.UserID
	} else if uid, ok := c.Locals("user_id").(string); ok {
		userIDStr = uid
	}

	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			req.ProcessedByID = parsed
		}
	}
	
	if errs := h.v.Validate(&req); len(errs) > 0 {
		return response.Error(c, fiber.StatusBadRequest, errs[0].Message)
	}
	
	result, err := h.uc.ApprovePlanning(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Planning approved", result)
}

func (h *PurchaseOrderPlanningHandler) Ignore(c *fiber.Ctx) error {
	var req struct {
		IDs []string `json:"ids" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}
	
	if errs := h.v.Validate(&req); len(errs) > 0 {
		return response.Error(c, fiber.StatusBadRequest, errs[0].Message)
	}
	
	if err := h.uc.IgnorePlanning(c.UserContext(), req.IDs); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Planning ignored", nil)
}