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

type SupplierCategoryHandler struct {
	categoryUseCase usecase.SupplierCategoryUseCase
	queryService    *service.SupplierCategoryQueryService
	validator       *validator.Validator
}

func NewSupplierCategoryHandler(
	categoryUseCase usecase.SupplierCategoryUseCase,
	queryService *service.SupplierCategoryQueryService,
	v *validator.Validator,
) *SupplierCategoryHandler {
	return &SupplierCategoryHandler{
		categoryUseCase: categoryUseCase,
		queryService:    queryService,
		validator:       v,
	}
}

func (h *SupplierCategoryHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSupplierCategoryRequest
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
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create supplier category")
	}

	return response.Success(c, fiber.StatusCreated, "Supplier category created successfully", resp)
}

func (h *SupplierCategoryHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	cat, err := h.queryService.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get supplier category")
	}
	if cat == nil {
		return response.Error(c, fiber.StatusNotFound, "supplier category not found")
	}

	resp := dto.SupplierCategoryResponse{
		ID:          cat.ID.String(),
		Name:        cat.Name,
		Description: cat.Description,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}

	return response.Success(c, fiber.StatusOK, "Supplier category retrieved successfully", resp)
}

func (h *SupplierCategoryHandler) GetAll(c *fiber.Ctx) error {
	cats, err := h.queryService.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get supplier categories")
	}

	resp := make([]dto.SupplierCategoryResponse, len(cats))
	for i, cat := range cats {
		resp[i] = dto.SupplierCategoryResponse{
			ID:          cat.ID.String(),
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	return response.Success(c, fiber.StatusOK, "Supplier categories retrieved successfully", resp)
}

func (h *SupplierCategoryHandler) GetAllWithPagination(c *fiber.Ctx) error {
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
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Supplier categories retrieved successfully", data, meta)
}

func (h *SupplierCategoryHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateSupplierCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.categoryUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrSupplierCategoryNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update supplier category")
	}

	return response.Success(c, fiber.StatusOK, "Supplier category updated successfully", resp)
}

func (h *SupplierCategoryHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.categoryUseCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrSupplierCategoryNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete supplier category")
	}

	return response.Success(c, fiber.StatusOK, "Supplier category deleted successfully", nil)
}
