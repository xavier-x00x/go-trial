package handler

import (
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/query/params"
	"go-trial/internal/query/service"
	"go-trial/internal/usecase"
	"go-trial/pkg/jwt"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PriceListHandler struct {
	uc           *usecase.PriceListUsecase
	queryService *service.PriceListQueryService
	v            *validator.Validator
}

func NewPriceListHandler(uc *usecase.PriceListUsecase, queryService *service.PriceListQueryService, v *validator.Validator) *PriceListHandler {
	return &PriceListHandler{uc: uc, queryService: queryService, v: v}
}

func (h *PriceListHandler) Create(c *fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*jwt.Claims)
	if !ok || claims == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req dto.CreatePriceListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Create(c.UserContext(), claims.StoreID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Price list created successfully", result)
}

func (h *PriceListHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.queryService.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price list retrieved successfully", result)
}

func (h *PriceListHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.queryService.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price lists retrieved successfully", result)
}

func (h *PriceListHandler) GetActive(c *fiber.Ctx) error {
	result, err := h.queryService.GetActive(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Active price lists retrieved successfully", result)
}

func (h *PriceListHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdatePriceListRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.Update(c.UserContext(), req, id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price list updated successfully", result)
}

func (h *PriceListHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Delete(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Price list deleted successfully", nil)
}

func (h *PriceListHandler) GetAllWithPagination(c *fiber.Ctx) error {
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

	return response.SuccessWithMeta(c, fiber.StatusOK, "Price lists retrieved successfully", data, meta)
}