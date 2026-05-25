package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PurchaseReturnHandler struct {
	uc usecase.PurchaseReturnUseCase
}

func NewPurchaseReturnHandler(uc usecase.PurchaseReturnUseCase) *PurchaseReturnHandler {
	return &PurchaseReturnHandler{uc: uc}
}

func (h *PurchaseReturnHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePurchaseReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	userID := c.Locals("user_id").(string)
	result, err := h.uc.Create(c.UserContext(), userID, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Purchase return created successfully", result)
}

func (h *PurchaseReturnHandler) Post(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	if err := h.uc.Post(c.UserContext(), userID, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Purchase return posted successfully", nil)
}

func (h *PurchaseReturnHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Purchase return retrieved successfully", result)
}

func (h *PurchaseReturnHandler) GetAll(c *fiber.Ctx) error {
	req := new(dto.MetaRequest)
	if err := c.QueryParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	data, meta, err := h.uc.GetAllWithPagination(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Purchase returns retrieved successfully", data, meta)
}
