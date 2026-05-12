package handler

import (
	"errors"
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type SupplierHandler struct {
	supplierUseCase usecase.SupplierUseCase
	validator     *validator.Validator
}

func NewSupplierHandler(supplierUseCase usecase.SupplierUseCase, v *validator.Validator) *SupplierHandler {
	return &SupplierHandler{
		supplierUseCase: supplierUseCase,
		validator:     v,
	}
}

func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.supplierUseCase.Create(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrSupplierCodeExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create supplier")
	}

	return response.Success(c, fiber.StatusCreated, "Supplier created successfully", resp)
}

func (h *SupplierHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.supplierUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrSupplierNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get supplier")
	}

	return response.Success(c, fiber.StatusOK, "Supplier retrieved successfully", resp)
}

func (h *SupplierHandler) GetAll(c *fiber.Ctx) error {
	resp, err := h.supplierUseCase.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get suppliers")
	}

	return response.Success(c, fiber.StatusOK, "Suppliers retrieved successfully", resp)
}

func (h *SupplierHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.supplierUseCase.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get suppliers")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Suppliers retrieved successfully", data, meta)
}

func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.supplierUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrSupplierNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrSupplierCodeExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update supplier")
	}

	return response.Success(c, fiber.StatusOK, "Supplier updated successfully", resp)
}

func (h *SupplierHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.supplierUseCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrSupplierNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete supplier")
	}

	return response.Success(c, fiber.StatusOK, "Supplier deleted successfully", nil)
}