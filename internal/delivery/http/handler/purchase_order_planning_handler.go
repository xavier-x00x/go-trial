package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PurchaseOrderPlanningHandler struct {
	uc usecase.PurchaseOrderPlanningUseCase
	v  *validator.Validator
}

func NewPurchaseOrderPlanningHandler(uc usecase.PurchaseOrderPlanningUseCase, v *validator.Validator) *PurchaseOrderPlanningHandler {
	return &PurchaseOrderPlanningHandler{uc: uc, v: v}
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

func (h *PurchaseOrderPlanningHandler) GetPending(c *fiber.Ctx) error {
	storeID := c.Query("store_id")
	if storeID == "" {
		return response.Error(c, fiber.StatusBadRequest, "store_id is required")
	}
	
	result, err := h.uc.GetPending(c.UserContext(), storeID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Pending planning retrieved", result)
}

func (h *PurchaseOrderPlanningHandler) GetAll(c *fiber.Ctx) error {
	storeID := c.Query("store_id")
	status := c.Query("status")
	
	result, err := h.uc.GetAll(c.UserContext(), storeID, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Planning list retrieved", result)
}

func (h *PurchaseOrderPlanningHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Planning detail retrieved", result)
}

func (h *PurchaseOrderPlanningHandler) Approve(c *fiber.Ctx) error {
	var req dto.ApprovePlanningRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
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
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	
	if err := h.uc.IgnorePlanning(c.UserContext(), req.IDs); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	
	return response.Success(c, fiber.StatusOK, "Planning ignored", nil)
}