package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PriceListHandler struct {
	uc *usecase.PriceListUsecase
	v  *validator.Validator
}

func NewPriceListHandler(uc *usecase.PriceListUsecase, v *validator.Validator) *PriceListHandler {
	return &PriceListHandler{uc: uc, v: v}
}

func (h *PriceListHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePriceListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Create(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Price list created successfully", result)
}

func (h *PriceListHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price list retrieved successfully", result)
}

func (h *PriceListHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.uc.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price lists retrieved successfully", result)
}

func (h *PriceListHandler) GetActive(c *fiber.Ctx) error {
	result, err := h.uc.GetActive(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Active price lists retrieved successfully", result)
}

func (h *PriceListHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdatePriceListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Update(c.UserContext(), req, id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price list updated successfully", result)
}

func (h *PriceListHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price list deleted successfully", nil)
}