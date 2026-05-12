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

type WarehouseHandler struct {
	uc *usecase.WarehouseUsecase
	v  *validator.Validator
}

func NewWarehouseHandler(uc *usecase.WarehouseUsecase, v *validator.Validator) *WarehouseHandler {
	return &WarehouseHandler{uc: uc, v: v}
}

func (h *WarehouseHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateWarehouseRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}
	result, err := h.uc.Create(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Warehouse created successfully", result)
}

func (h *WarehouseHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrWarehouseNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Warehouse retrieved successfully", result)
}

func (h *WarehouseHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.uc.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Warehouses retrieved successfully", result)
}

func (h *WarehouseHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.uc.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Warehouses retrieved successfully", data, meta)
}

func (h *WarehouseHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateWarehouseRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Update(c.UserContext(), req, id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Warehouse updated successfully", result)
}

func (h *WarehouseHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Warehouse deleted successfully", nil)
}