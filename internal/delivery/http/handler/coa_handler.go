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

type COAHandler struct {
	coaUseCase usecase.ChartOfAccountUseCase
	validator *validator.Validator
}

func NewCOAHandler(coaUseCase usecase.ChartOfAccountUseCase, v *validator.Validator) *COAHandler {
	return &COAHandler{
		coaUseCase: coaUseCase,
		validator: v,
	}
}

func (h *COAHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCOARequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.coaUseCase.Create(c.UserContext(), req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to create chart of account")
	}

	return response.Success(c, fiber.StatusCreated, "Chart of account created successfully", resp)
}

func (h *COAHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	resp, err := h.coaUseCase.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrCOANotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get chart of account")
	}

	return response.Success(c, fiber.StatusOK, "Chart of account retrieved successfully", resp)
}

func (h *COAHandler) GetAll(c *fiber.Ctx) error {
	resp, err := h.coaUseCase.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get chart of accounts")
	}

	return response.Success(c, fiber.StatusOK, "Chart of accounts retrieved successfully", resp)
}

func (h *COAHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.coaUseCase.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get chart of accounts")
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Chart of accounts retrieved successfully", data, meta)
}

func (h *COAHandler) GetByType(c *fiber.Ctx) error {
	accountType := c.Params("type")

	resp, err := h.coaUseCase.GetByType(c.UserContext(), accountType)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get chart of accounts by type")
	}

	return response.Success(c, fiber.StatusOK, "Chart of accounts retrieved successfully", resp)
}

func (h *COAHandler) GetTree(c *fiber.Ctx) error {
	resp, err := h.coaUseCase.GetTree(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get chart of accounts tree")
	}

	return response.Success(c, fiber.StatusOK, "Chart of accounts tree retrieved successfully", resp)
}

func (h *COAHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.coaUseCase.Delete(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrCOANotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrCOAHasChildren) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to delete chart of account")
	}

	return response.Success(c, fiber.StatusOK, "Chart of account deleted successfully", nil)
}

func (h *COAHandler) Import(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "file is required")
	}

	f, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to open file")
	}
	defer f.Close()

	result, err := h.coaUseCase.Import(c.UserContext(), f, file.Filename)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Import completed", result)
}

func (h *COAHandler) DownloadTemplate(c *fiber.Ctx) error {
	data, err := h.coaUseCase.GenerateTemplate(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to generate template")
	}

	c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="account_import_template.xlsx"`)
	return c.Send(data)
}

func (h *COAHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateCOARequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	resp, err := h.coaUseCase.Update(c.UserContext(), id, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		if errors.Is(err, usecase.ErrCOANotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to update chart of account")
	}

	return response.Success(c, fiber.StatusOK, "Chart of account updated successfully", resp)
}