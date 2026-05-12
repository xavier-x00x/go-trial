package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ProductPriceHandler struct {
	uc *usecase.ProductPriceUsecase
	v  *validator.Validator
}

func NewProductPriceHandler(uc *usecase.ProductPriceUsecase, v *validator.Validator) *ProductPriceHandler {
	return &ProductPriceHandler{uc: uc, v: v}
}

func (h *ProductPriceHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProductPriceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Create(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Product price created successfully", result)
}

func (h *ProductPriceHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product price retrieved successfully", result)
}

func (h *ProductPriceHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.uc.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product prices retrieved successfully", result)
}

func (h *ProductPriceHandler) GetByProductID(c *fiber.Ctx) error {
	productID := c.Params("productId")
	result, err := h.uc.GetByProductID(c.UserContext(), productID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product prices retrieved successfully", result)
}

func (h *ProductPriceHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateProductPriceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Update(c.UserContext(), req, id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product price updated successfully", result)
}

func (h *ProductPriceHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product price deleted successfully", nil)
}