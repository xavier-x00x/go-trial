package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/query/params"
	"go-trial/internal/query/service"
	"go-trial/internal/usecase"
	"go-trial/pkg/jwt"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (h *MasterDataProposalHandler) validatePayloads(entityType, actionType string, items []dto.CreateMasterDataProposalItemInput) []validator.ValidationError {
	var allErrs []validator.ValidationError

	for i, item := range items {
		var req interface{}

		switch entityType {
		case "PRODUCT":
			if actionType == "CREATE" {
				req = &dto.CreateProductRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateProductRequest{}
			}
		case "PRODUCT_PRICE":
			if actionType == "CREATE" {
				req = &dto.CreateProductPriceRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateProductPriceRequest{}
			}
		case "PRODUCT_UOM_CONVERSION":
			if actionType == "CREATE" {
				req = &dto.CreateProductUOMConversionRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateProductUOMConversionRequest{}
			}
		case "SUPPLIER":
			if actionType == "CREATE" {
				req = &dto.CreateSupplierRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateSupplierRequest{}
			}
		case "PRODUCT_SUPPLIER":
			if actionType == "CREATE" {
				req = &dto.CreateProductSupplierRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateProductSupplierRequest{}
			}
		case "CHART_OF_ACCOUNT":
			if actionType == "CREATE" {
				req = &dto.CreateChartOfAccountRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateChartOfAccountRequest{}
			}
		case "TAX":
			if actionType == "CREATE" {
				req = &dto.CreateTaxRequest{}
			} else if actionType == "UPDATE" {
				req = &dto.UpdateTaxRequest{}
			}
		}

		if req != nil {
			if err := json.Unmarshal([]byte(item.PayloadJSON), req); err != nil {
				allErrs = append(allErrs, validator.ValidationError{
					Field:   fmt.Sprintf("items[%d].payload_json", i),
					Message: "Invalid JSON structure",
				})
				continue
			}

			errs := h.v.Validate(req)
			if len(errs) > 0 {
				for _, e := range errs {
					allErrs = append(allErrs, validator.ValidationError{
						Field:   fmt.Sprintf("items[%d].%s", i, e.Field),
						Message: e.Message,
					})
				}
			}
		}
	}

	return allErrs
}

func (h *MasterDataProposalHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMasterDataProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.v.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}
	
	if payloadErrs := h.validatePayloads(req.EntityType, req.ActionType, req.Items); len(payloadErrs) > 0 {
		return response.ValidationError(c, "Payload validation failed", payloadErrs)
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

	proposal, err := h.queryService.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	if proposal == nil {
		return response.Error(c, fiber.StatusNotFound, "Proposal not found")
	}

	if payloadErrs := h.validatePayloads(proposal.EntityType, proposal.ActionType, req.Items); len(payloadErrs) > 0 {
		return response.ValidationError(c, "Payload validation failed", payloadErrs)
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

func (h *MasterDataProposalHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := getUserIDFromContext(c)

	err := h.uc.Delete(c.UserContext(), userID, id)
	if err != nil {
		if errors.Is(err, usecase.ErrProposalNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrProposalNotPending) {
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Proposal deleted successfully", nil)
}

func (h *MasterDataProposalHandler) GeneratePricesFromGR(c *fiber.Ctx) error {
	userIDStr := getUserIDFromContext(c)
	if userIDStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	count, err := h.uc.GenerateProductPricesFromTodayGR(c.UserContext(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.StatusOK, fmt.Sprintf("Successfully generated %d price proposal documents from today's Goods Receipts", count), map[string]interface{}{
		"proposals_generated": count,
	})
}