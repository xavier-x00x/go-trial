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

type CustomerHandler struct {
	customerUseCase usecase.CustomerUseCase
	validator      *validator.Validator
}

func NewCustomerHandler(customerUseCase usecase.CustomerUseCase, v *validator.Validator) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
		validator:      v,
	}
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.customerUseCase.Create(c.UserContext(), req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create customer")
	}

	return response.Success(c, fiber.StatusCreated, "Customer created successfully", resp)
}

func (h *CustomerHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.customerUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrCustomerNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get customer")
	}

	return response.Success(c, fiber.StatusOK, "Customer retrieved successfully", resp)
}

func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	resp, err := h.customerUseCase.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get customers")
	}

	return response.Success(c, fiber.StatusOK, "Customers retrieved successfully", resp)
}

func (h *CustomerHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.customerUseCase.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get customers")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Customers retrieved successfully", data, meta)
}

func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.customerUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrCustomerNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update customer")
	}

	return response.Success(c, fiber.StatusOK, "Customer updated successfully", resp)
}

func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.customerUseCase.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrCustomerNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete customer")
	}

	return response.Success(c, fiber.StatusOK, "Customer deleted successfully", nil)
}