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

type StoreHandler struct {
	storeUseCase usecase.StoreUseCase
	validator    *validator.Validator
}

func NewStoreHandler(storeUseCase usecase.StoreUseCase, validator *validator.Validator) *StoreHandler {
	return &StoreHandler{
		storeUseCase: storeUseCase,
		validator:    validator,
	}
}

// Create handles POST /api/stores
func (h *StoreHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	storeResp, err := h.storeUseCase.Create(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrStoreCodeExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create store")
	}

	return response.Success(c, fiber.StatusCreated, "Store created successfully", storeResp)
}

// GetByID handles GET /api/stores/:id
func (h *StoreHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	storeResp, err := h.storeUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get store")
	}

	return response.Success(c, fiber.StatusOK, "Store retrieved successfully", storeResp)
}

// GetAll handles GET /api/stores
func (h *StoreHandler) GetAll(c *fiber.Ctx) error {
	stores, err := h.storeUseCase.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get stores")
	}

	return response.Success(c, fiber.StatusOK, "Stores retrieved successfully", stores)
}

// GetAllWithPagination handles GET /api/stores with pagination
func (h *StoreHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.storeUseCase.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get stores")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Stores retrieved successfully", data, meta)
}

// Update handles PUT /api/stores/:id
func (h *StoreHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	storeResp, err := h.storeUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrStoreCodeExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update store")
	}

	return response.Success(c, fiber.StatusOK, "Store updated successfully", storeResp)
}

// Delete handles DELETE /api/stores/:id
func (h *StoreHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.storeUseCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrStoreNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete store")
	}

	return response.Success(c, fiber.StatusOK, "Store deleted successfully", nil)
}
