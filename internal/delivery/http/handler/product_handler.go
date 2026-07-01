package handler

import (
	"errors"
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/query/params"
	"go-trial/internal/query/service"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUseCase usecase.ProductUseCase
	queryService   *service.ProductQueryService
	validator      *validator.Validator
}

func NewProductHandler(productUseCase usecase.ProductUseCase, queryService *service.ProductQueryService, v *validator.Validator) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		queryService:   queryService,
		validator:      v,
	}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.productUseCase.Create(c.UserContext(), req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create product")
	}

	return response.Success(c, fiber.StatusCreated, "Product created successfully", resp)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.queryService.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get product")
	}
	if resp == nil {
		return response.Error(c, fiber.StatusNotFound, "Product not found")
	}

	return response.Success(c, fiber.StatusOK, "Product retrieved successfully", resp)
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	resp, err := h.queryService.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return response.Success(c, fiber.StatusOK, "Products retrieved successfully", resp)
}

func (h *ProductHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &params.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.queryService.GetListPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Products retrieved successfully", data, meta)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.productUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update product")
	}

	return response.Success(c, fiber.StatusOK, "Product updated successfully", resp)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.productUseCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrProductNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete product")
	}

	return response.Success(c, fiber.StatusOK, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetProductSuppliers(c *fiber.Ctx) error {
	id := c.Params("id")
	resp, err := h.queryService.GetProductSuppliers(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get product suppliers")
	}
	return response.Success(c, fiber.StatusOK, "Product suppliers retrieved successfully", resp)
}

func (h *ProductHandler) GetProductsBySupplier(c *fiber.Ctx) error {
	supplierID := c.Query("supplier_id")
	resp, err := h.queryService.GetProductsBySupplier(c.UserContext(), supplierID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products by supplier")
	}
	return response.Success(c, fiber.StatusOK, "Products retrieved successfully", resp)
}

func (h *ProductHandler) GetProductsBySupplierWithPagination(c *fiber.Ctx) error {
	supplierID := c.Query("supplier_id")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &params.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", ""),
		OrderDir:    c.Query("order_dir", ""),
	}

	data, meta, err := h.queryService.GetProductsBySupplierWithPagination(c.UserContext(), supplierID, metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get products by supplier")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Products retrieved successfully", data, meta)
}
