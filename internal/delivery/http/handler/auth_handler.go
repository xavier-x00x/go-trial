package handler

import (
	"errors"
	"strconv"
	"time"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/usecase"
	jwtPkg "go-trial/pkg/jwt"
	"go-trial/pkg/response"
	"go-trial/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	validator   *validator.Validator
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator,
	}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	authResp, refreshToken, err := h.authUseCase.Register(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		if errors.Is(err, usecase.ErrUsernameAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to register user")
	}

	setRefreshTokenCookie(c, refreshToken)

	return response.Success(c, fiber.StatusCreated, "User registered successfully", authResp)
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	authResp, refreshToken, err := h.authUseCase.Login(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return response.Error(c, fiber.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, usecase.ErrAccountInactive) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		if errors.Is(err, usecase.ErrRoleNotAssigned) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to login")
	}

	setRefreshTokenCookie(c, refreshToken)

	return response.Success(c, fiber.StatusOK, "Login successful", authResp)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Error(c, fiber.StatusUnauthorized, "Refresh token not found")
	}

	resp, err := h.authUseCase.RefreshToken(c.UserContext(), refreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRefreshToken) {
			clearRefreshTokenCookie(c)
			return response.Error(c, fiber.StatusUnauthorized, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to refresh token")
	}

	return response.Success(c, fiber.StatusOK, "Token refreshed successfully", resp)
}

// Me handles GET /api/auth/me (protected)
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*jwtPkg.Claims)

	userResp, err := h.authUseCase.GetMe(c.UserContext(), claims.UserID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get user")
	}

	return response.Success(c, fiber.StatusOK, "User retrieved successfully", userResp)
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	clearRefreshTokenCookie(c)
	return response.Success(c, fiber.StatusOK, "Logged out successfully", nil)
}

// GetAllUsers handles GET /api/users (admin only)
func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
	result, err := h.authUseCase.GetAllUsers(c.UserContext())
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "Users retrieved successfully", result)
}

// GetAllUsersWithPagination handles GET /api/users/pagination (admin only)
func (h *AuthHandler) GetAllUsersWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	metaRequest := &dto.MetaRequest{
		Page:        page,
		Limit:       limit,
		Search:      c.Query("search", ""),
		OrderColumn: c.Query("order_column", "id"),
		OrderDir:    c.Query("order_dir", "asc"),
	}

	data, meta, err := h.authUseCase.GetAllUsersWithPagination(c.UserContext(), metaRequest)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Users retrieved successfully", data, meta)
}

// GetUserByID handles GET /api/users/:id (admin only)
func (h *AuthHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	resp, err := h.authUseCase.GetMe(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User retrieved successfully", resp)
}

// UpdateUser handles PUT /api/users/:id (admin only)
func (h *AuthHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateUserByAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	result, err := h.authUseCase.UpdateUser(c.UserContext(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, usecase.ErrEmailAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		if errors.Is(err, usecase.ErrUsernameAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User updated successfully", result)
}

// DeleteUser handles DELETE /api/users/:id (admin only)
func (h *AuthHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.authUseCase.DeleteUser(c.UserContext(), id); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return response.Error(c, fiber.StatusNotFound, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, fiber.StatusOK, "User deleted successfully", nil)
}

// GoogleRedirect handles GET /api/auth/google
func (h *AuthHandler) GoogleRedirect(c *fiber.Ctx) error {
	url := h.authUseCase.GetGoogleLoginURL()
	return c.Redirect(url)
}

// GoogleCallback handles GET /api/auth/google/callback
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return response.Error(c, fiber.StatusBadRequest, "Code not found")
	}

	authResp, refreshToken, err := h.authUseCase.GoogleLogin(c.UserContext(), code)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotAssigned) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	setRefreshTokenCookie(c, refreshToken)

	return response.Success(c, fiber.StatusOK, "Login successful via Google", authResp)
}

// GoogleTokenLogin handles POST /api/auth/google/token
func (h *AuthHandler) GoogleTokenLogin(c *fiber.Ctx) error {
	var req dto.GoogleTokenLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if errs := h.validator.Validate(req); len(errs) > 0 {
		return response.ValidationError(c, "Validation failed", errs)
	}

	authResp, refreshToken, err := h.authUseCase.GoogleLoginWithToken(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotAssigned) {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	setRefreshTokenCookie(c, refreshToken)

	return response.Success(c, fiber.StatusOK, "Login successful via Google Token", authResp)
}

// setRefreshTokenCookie sets the refresh token as an HTTPOnly cookie.
func setRefreshTokenCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  time.Now().Add(jwtPkg.RefreshTokenDuration),
		HTTPOnly: true,
		Secure:   false, // set to true in production with HTTPS
		SameSite: "Lax",
		Path:     "/",
	})
}

// clearRefreshTokenCookie removes the refresh token cookie.
func clearRefreshTokenCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})
}
