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

type ProductCategoryHandler struct {
	categoryUseCase usecase.ProductCategoryUseCase
	validator       *validator.Validator
}

func NewProductCategoryHandler(categoryUseCase usecase.ProductCategoryUseCase, v *validator.Validator) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		categoryUseCase: categoryUseCase,
		validator:       v,
	}
}

func (h *ProductCategoryHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.categoryUseCase.Create(c.UserContext(), req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create category")
	}

	return response.Success(c, fiber.StatusCreated, "Category created successfully", resp)
}

func (h *ProductCategoryHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.categoryUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get category")
	}

	return response.Success(c, fiber.StatusOK, "Category retrieved successfully", resp)
}

func (h *ProductCategoryHandler) GetAll(c *fiber.Ctx) error {
	resp, err := h.categoryUseCase.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get categories")
	}

	return response.Success(c, fiber.StatusOK, "Categories retrieved successfully", resp)
}

func (h *ProductCategoryHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.categoryUseCase.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get categories")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Categories retrieved successfully", data, meta)
}

func (h *ProductCategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.categoryUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update category")
	}

	return response.Success(c, fiber.StatusOK, "Category updated successfully", resp)
}

func (h *ProductCategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.categoryUseCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete category")
	}

	return response.Success(c, fiber.StatusOK, "Category deleted successfully", nil)
}