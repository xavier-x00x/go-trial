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

type MasterDataProposalHandler struct {
	uc           usecase.MasterDataProposalUseCase
	queryService *service.MasterDataProposalQueryService
	v            *validator.Validator
}

func NewMasterDataProposalHandler(uc usecase.MasterDataProposalUseCase, qs *service.MasterDataProposalQueryService, v *validator.Validator) *MasterDataProposalHandler {
	return &MasterDataProposalHandler{uc: uc, queryService: qs, v: v}
}

func getUserIDFromContext(c *fiber.Ctx) string {
	claims, ok := c.Locals("claims").(*jwt.Claims)
	if !ok {
		return ""
	}
	return claims.UserID
}

func (h *MasterDataProposalHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMasterDataProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}
	
	userID := getUserIDFromContext(c)
	result, err := h.uc.Create(c.UserContext(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Proposal created successfully", result)
}

func (h *MasterDataProposalHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.queryService.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if result == nil {
		return response.Error(c, fiber.StatusNotFound, "Proposal not found")
	}
	return response.Success(c, fiber.StatusOK, "Proposal retrieved successfully", result)
}

func (h *MasterDataProposalHandler) GetPending(c *fiber.Ctx) error {
	result, err := h.queryService.GetPending(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Pending proposals retrieved successfully", result)
}

func (h *MasterDataProposalHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.queryService.GetAll(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Proposals retrieved successfully", result)
}

func (h *MasterDataProposalHandler) GetAllWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &params.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "created_at"),
		OrderDir:    c.Query("order_dir", "desc"),
		Conditions: map[string]interface{}{
			"m.status":      c.Query("status", ""),
			"m.entity_type": c.Query("entity_type", ""),
		},
	}

	data, meta, err := h.queryService.GetListPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Proposals retrieved successfully", data, meta)
}

func (h *MasterDataProposalHandler) GetByEntityType(c *fiber.Ctx) error {
	entityType := c.Params("entityType")
	status := c.Query("status")
	result, err := h.queryService.GetByEntityType(c.UserContext(), entityType, status)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Proposals retrieved successfully", result)
}

func (h *MasterDataProposalHandler) Review(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.ReviewMasterDataProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	
	userID := getUserIDFromContext(c)
	result, err := h.uc.Review(c.UserContext(), userID, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Proposal reviewed successfully", result)
}

func (h *MasterDataProposalHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateMasterDataProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	userID := getUserIDFromContext(c)
	result, err := h.uc.Update(c.UserContext(), userID, id, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Proposal updated successfully", result)
}

func (h *MasterDataProposalHandler) Execute(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.Execute(c.UserContext(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Proposal executed successfully", nil)
}

func (h *MasterDataProposalHandler) BulkLinkProductSupplier(c *fiber.Ctx) error {
	var req dto.BulkCreateProductSupplierProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	
	userID := getUserIDFromContext(c)
	result, err := h.uc.BulkLinkProductSupplier(c.UserContext(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Bulk product-supplier linking proposals created successfully", result)
}

func (h *MasterDataProposalHandler) GetByGroup(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	result, err := h.queryService.GetByGroup(c.UserContext(), groupID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Proposals by group retrieved successfully", result)
}