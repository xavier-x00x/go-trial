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

type UOMHandler struct {
	uomUseCase usecase.UOMUseCase
	validator *validator.Validator
}

func NewUOMHandler(uomUseCase usecase.UOMUseCase, v *validator.Validator) *UOMHandler {
	return &UOMHandler{
		uomUseCase: uomUseCase,
		validator: v,
	}
}

func (h *UOMHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUOMRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.uomUseCase.Create(c.UserContext(), req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create UOM")
	}

	return response.Success(c, fiber.StatusCreated, "UOM created successfully", resp)
}

func (h *UOMHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.uomUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrUOMNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get UOM")
	}

	return response.Success(c, fiber.StatusOK, "UOM retrieved successfully", resp)
}

func (h *UOMHandler) GetAll(c *fiber.Ctx) error {
	resp, err := h.uomUseCase.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get UOMs")
	}

	return response.Success(c, fiber.StatusOK, "UOMs retrieved successfully", resp)
}

func (h *UOMHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.uomUseCase.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get UOMs")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "UOMs retrieved successfully", data, meta)
}

func (h *UOMHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateUOMRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.uomUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrUOMNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update UOM")
	}

	return response.Success(c, fiber.StatusOK, "UOM updated successfully", resp)
}