package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ProductUOMConversionHandler struct {
	uc *usecase.ProductUOMConversionUsecase
	v  *validator.Validator
}

func NewProductUOMConversionHandler(uc *usecase.ProductUOMConversionUsecase, v *validator.Validator) *ProductUOMConversionHandler {
	return &ProductUOMConversionHandler{uc: uc, v: v}
}

func (h *ProductUOMConversionHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProductUOMConversionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Create(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Product UOM conversion created successfully", result)
}

func (h *ProductUOMConversionHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product UOM conversion retrieved successfully", result)
}

func (h *ProductUOMConversionHandler) GetByProductID(c *fiber.Ctx) error {
	productID := c.Params("productId")
	result, err := h.uc.GetByProductID(c.UserContext(), productID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product UOM conversions retrieved successfully", result)
}

func (h *ProductUOMConversionHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.uc.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product UOM conversions retrieved successfully", result)
}

func (h *ProductUOMConversionHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateProductUOMConversionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Update(c.UserContext(), req, id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product UOM conversion updated successfully", result)
}

func (h *ProductUOMConversionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Product UOM conversion deleted successfully", nil)
}