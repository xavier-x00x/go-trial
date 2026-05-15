package handler

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ExpenseVoucherHandler struct {
	uc        usecase.ExpenseVoucherUseCase
	validator *validator.Validator
}

func NewExpenseVoucherHandler(uc usecase.ExpenseVoucherUseCase, v *validator.Validator) *ExpenseVoucherHandler {
	return &ExpenseVoucherHandler{
		uc:        uc,
		validator: v,
	}
}

func (h *ExpenseVoucherHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateExpenseVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if errs := h.validator.Validate(req); errs != nil {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := c.Locals("user_id").(string)
	result, err := h.uc.Create(c.UserContext(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusCreated, "Expense voucher created successfully", result)
}

func (h *ExpenseVoucherHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateExpenseVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	if errs := h.validator.Validate(req); errs != nil {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := c.Locals("user_id").(string)
	if err := h.uc.Update(c.UserContext(), userID, id, req); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Expense voucher updated successfully", nil)
}

func (h *ExpenseVoucherHandler) Post(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	if err := h.uc.Post(c.UserContext(), userID, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Expense voucher posted successfully", nil)
}

func (h *ExpenseVoucherHandler) Cancel(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	if err := h.uc.Cancel(c.UserContext(), userID, id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Expense voucher cancelled successfully", nil)
}

func (h *ExpenseVoucherHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Expense voucher retrieved successfully", result)
}

func (h *ExpenseVoucherHandler) GetAll(c *fiber.Ctx) error {
	req := new(dto.MetaRequest)
	if err := c.QueryParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	data, meta, err := h.uc.GetAllWithPagination(c.UserContext(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Expense vouchers retrieved successfully", data, meta)
}
