package handler

import (
	"net/http"
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type PurchaseInvoiceHandler struct {
	uc usecase.PurchaseInvoiceUseCase
}

func NewPurchaseInvoiceHandler(uc usecase.PurchaseInvoiceUseCase) *PurchaseInvoiceHandler {
	return &PurchaseInvoiceHandler{uc: uc}
}

func (h *PurchaseInvoiceHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req dto.CreatePurchaseInvoiceRequest
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

func (h *PurchaseInvoiceHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	var req dto.UpdatePurchaseInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.uc.Update(c.Context(), userID.(string), id, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (h *PurchaseInvoiceHandler) GetByID(c *fiber.Ctx) error {
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

func (h *PurchaseInvoiceHandler) GetByInvoiceNumber(c *fiber.Ctx) error {
	invoiceNum := c.Params("invoice_number")
	if invoiceNum == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invoice_number is required"})
	}

	resp, err := h.uc.GetByInvoiceNumber(c.Context(), invoiceNum)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (h *PurchaseInvoiceHandler) GetByStoreID(c *fiber.Ctx) error {
	storeID := c.Query("store_id")
	if storeID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "store_id is required"})
	}

	status := c.Query("status")

	resp, err := h.uc.GetByStoreID(c.Context(), storeID, status)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

func (h *PurchaseInvoiceHandler) List(c *fiber.Ctx) error {
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

func (h *PurchaseInvoiceHandler) Submit(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	if err := h.uc.Submit(c.Context(), userID.(string), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "submitted"})
}

func (h *PurchaseInvoiceHandler) Approve(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	if err := h.uc.Approve(c.Context(), userID.(string), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "approved"})
}

func (h *PurchaseInvoiceHandler) Verify(c *fiber.Ctx) error {
	return h.Approve(c)
}

func (h *PurchaseInvoiceHandler) Post(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	if err := h.uc.Post(c.Context(), userID.(string), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "posted"})
}

func (h *PurchaseInvoiceHandler) Cancel(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.uc.Cancel(c.Context(), userID.(string), id, req.Reason); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "cancelled"})
}
