package handler

import (
	"errors"
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	"go-trial/pkg/jwt"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type GoodsReceiptHandler struct {
	uc usecase.GoodsReceiptUseCase
	v  *validator.Validator
}

func NewGoodsReceiptHandler(uc usecase.GoodsReceiptUseCase, v *validator.Validator) *GoodsReceiptHandler {
	return &GoodsReceiptHandler{uc: uc, v: v}
}

func getUserIDFromGRContext(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*jwt.Claims)
	if !ok {
		return ""
	}
	return claims.UserID
}

func (h *GoodsReceiptHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateGoodsReceiptRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromGRContext(c)
	result, err := h.uc.Create(c.UserContext(), userID, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		if errors.Is(err, usecase.ErrOverReceiveNeedsPIN) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Goods receipt created successfully", result)
}

func (h *GoodsReceiptHandler) Update(c *fiber.Ctx) error {
	var req dto.UpdateGoodsReceiptRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	id := c.Params("id")
	userID := getUserIDFromGRContext(c)
	result, err := h.uc.Update(c.UserContext(), userID, id, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		if errors.Is(err, usecase.ErrOverReceiveNeedsPIN) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Goods receipt updated successfully", result)
}

func (h *GoodsReceiptHandler) Confirm(c *fiber.Ctx) error {
	var req dto.ConfirmGoodsReceiptRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	id := c.Params("id")
	userID := getUserIDFromGRContext(c)
	if err := h.uc.Confirm(c.UserContext(), userID, id, req); err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Goods receipt confirmed successfully", nil)
}

func (h *GoodsReceiptHandler) Cancel(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := getUserIDFromGRContext(c)
	if err := h.uc.Cancel(c.UserContext(), userID, id, ""); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Goods receipt cancelled successfully", nil)
}

func (h *GoodsReceiptHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Goods receipt retrieved successfully", result)
}

func (h *GoodsReceiptHandler) GetByPurchaseOrderID(c *fiber.Ctx) error {
	poID := c.Params("poId")
	result, err := h.uc.GetByPurchaseOrderID(c.UserContext(), poID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Goods receipts retrieved successfully", result)
}

func (h *GoodsReceiptHandler) GetByWarehouseID(c *fiber.Ctx) error {
	warehouseID := c.Params("warehouseId")
	status := c.Query("status")
	result, err := h.uc.GetByWarehouseID(c.UserContext(), warehouseID, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Goods receipts retrieved successfully", result)
}

func (h *GoodsReceiptHandler) CreateWithInvoice(c *fiber.Ctx) error {
	var req dto.CreateGoodsReceiptWithInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromGRContext(c)
	result, err := h.uc.CreateWithInvoice(c.UserContext(), userID, req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		if errors.Is(err, usecase.ErrOverReceiveNeedsPIN) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Goods receipt with invoice created successfully", result)
}

func (h *GoodsReceiptHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "created_at"),
		OrderDir:    c.Query("order_dir", "desc"),
		Conditions:  parseGRConditions(c),
	}

	data, meta, err := h.uc.GetAllWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Goods receipts retrieved successfully", data, meta)
}

func parseGRConditions(c *fiber.Ctx) map[string]interface{} {
	conditions := make(map[string]interface{})
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		conditions["warehouse_id"] = warehouseID
	}
	if status := c.Query("status"); status != "" {
		conditions["status"] = status
	}
	if poID := c.Query("po_id"); poID != "" {
		conditions["purchase_order_id"] = poID
	}
	return conditions
}