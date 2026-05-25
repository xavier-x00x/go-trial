package handler

import (
	"net/http"
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type PurchasePaymentHandler struct {
	uc usecase.PurchasePaymentUseCase
}

func NewPurchasePaymentHandler(uc usecase.PurchasePaymentUseCase) *PurchasePaymentHandler {
	return &PurchasePaymentHandler{uc: uc}
}

func (h *PurchasePaymentHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req dto.CreatePurchasePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.uc.Create(c.Context(), userID.(string), req)
	if err != nil {
		if handleFieldErrors(c, err) {
			return nil
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(resp)
}

func (h *PurchasePaymentHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	resp, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (h *PurchasePaymentHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search")
	orderBy := c.Query("order_by", "created_at")
	orderDir := c.Query("order_dir", "desc")

	meta := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      search,
		OrderColumn: orderBy,
		OrderDir:    orderDir,
	}

	resp, resMeta, err := h.uc.GetAllWithPagination(c.Context(), meta)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": resp,
		"meta": resMeta,
	})
}

func (h *PurchasePaymentHandler) Post(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	var req dto.PostPurchasePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.uc.Post(c.Context(), userID.(string), id, req); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "posted"})
}

func (h *PurchasePaymentHandler) Void(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	var req dto.VoidPurchasePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.uc.Void(c.Context(), userID.(string), id, req.Reason); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "voided"})
}