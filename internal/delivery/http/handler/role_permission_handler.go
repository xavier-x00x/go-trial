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

type RolePermissionHandler struct {
	uc  usecase.RolePermissionUseCase
	v   *validator.Validator
}

func NewRolePermissionHandler(uc usecase.RolePermissionUseCase, v *validator.Validator) *RolePermissionHandler {
	return &RolePermissionHandler{uc: uc, v: v}
}

func (h *RolePermissionHandler) CreateRole(c *fiber.Ctx) error {
	var req dto.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}
	
	if req.Name == "" {
		return response.Error(c, fiber.StatusBadRequest, "name is required")
	}
	
	result, err := h.uc.CreateRole(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNameExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Role created successfully", result)
}

func (h *RolePermissionHandler) GetRole(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetRole(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role retrieved successfully", result)
}

func (h *RolePermissionHandler) GetRoleByName(c *fiber.Ctx) error {
	name := c.Params("name")
	result, err := h.uc.GetRoleByName(c.UserContext(), name)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role retrieved successfully", result)
}

func (h *RolePermissionHandler) ListRoles(c *fiber.Ctx) error {
	result, err := h.uc.ListRoles(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Roles retrieved successfully", result)
}

func (h *RolePermissionHandler) ListRolesWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.uc.ListRolesWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Roles retrieved successfully", data, meta)
}

func (h *RolePermissionHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.UpdateRole(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrRoleNameExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role updated successfully", result)
}

func (h *RolePermissionHandler) DeleteRole(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.DeleteRole(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Role deleted successfully", nil)
}

func (h *RolePermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var req dto.CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}
	
	if req.Path == "" {
		return response.Error(c, fiber.StatusBadRequest, "path is required")
	}
	if req.Name == "" {
		return response.Error(c, fiber.StatusBadRequest, "name is required")
	}
	
	result, err := h.uc.CreatePermission(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrPermPathExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusCreated, "Permission created successfully", result)
}

func (h *RolePermissionHandler) GetPermission(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.uc.GetPermission(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrPermNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Permission retrieved successfully", result)
}

func (h *RolePermissionHandler) ListPermissions(c *fiber.Ctx) error {
	result, err := h.uc.ListPermissions(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Permissions retrieved successfully", result)
}

func (h *RolePermissionHandler) ListPermissionsWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.uc.ListPermissionsWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Permissions retrieved successfully", data, meta)
}

func (h *RolePermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.uc.UpdatePermission(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrPermNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrPermPathExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Permission updated successfully", result)
}

func (h *RolePermissionHandler) DeletePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.uc.DeletePermission(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrPermNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Permission deleted successfully", nil)
}