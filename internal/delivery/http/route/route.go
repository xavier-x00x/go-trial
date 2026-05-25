package route

import (
	"go-trial/internal/infrastructure/middleware"
	"go-trial/internal/registry"
	jwtPkg "go-trial/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

// check adalah helper tingkat paket untuk memverifikasi hak akses permission.
var check func(permission string) fiber.Handler

func Setup(app *fiber.App, reg *registry.Registry, jwtManager *jwtPkg.JWTManager) {
	check = func(permission string) fiber.Handler {
		return middleware.RequirePermission(reg.DB, permission)
	}

	api := app.Group("/api")
	api.Post("/auth/register", reg.Auth.Handler.Register)
	api.Post("/auth/login", reg.Auth.Handler.Login)
	api.Get("/auth/google", reg.Auth.Handler.GoogleRedirect)
	api.Get("/auth/google/callback", reg.Auth.Handler.GoogleCallback)
	api.Post("/auth/google/token", reg.Auth.Handler.GoogleTokenLogin)
	api.Post("/auth/google/register", reg.Auth.Handler.RegisterWithGoogle)
	api.Post("/auth/refresh", reg.Auth.Handler.RefreshToken)
	api.Post("/auth/logout", reg.Auth.Handler.Logout)

	protected := api.Group("", middleware.AuthMiddleware(jwtManager))
	protected.Get("/auth/me", reg.Auth.Handler.Me)

	setupUsers(protected, reg)
	setupRoles(protected, reg)
	setupPermissions(protected, reg)
	setupOperationalRoutes(protected, reg)
	setupFinanceRoutes(protected, reg)
}

func setupRoles(r fiber.Router, reg *registry.Registry) {
	h := reg.Auth.RolePermissionHandler
	r.Post("/roles", h.CreateRole)
	r.Get("/roles", h.ListRoles)
	r.Get("/roles/pagination", h.ListRolesWithPagination)
	r.Get("/roles/name/:name", h.GetRoleByName)
	r.Get("/roles/:id", h.GetRole)
	r.Put("/roles/:id", h.UpdateRole)
	r.Delete("/roles/:id", h.DeleteRole)
}

func setupPermissions(r fiber.Router, reg *registry.Registry) {
	h := reg.Auth.RolePermissionHandler
	r.Post("/permissions", h.CreatePermission)
	r.Post("/permissions/sync", h.SyncPermissions)
	r.Get("/permissions", h.ListPermissions)
	r.Get("/permissions/pagination", h.ListPermissionsWithPagination)
	r.Get("/permissions/:id", h.GetPermission)
	r.Put("/permissions/:id", h.UpdatePermission)
	r.Delete("/permissions/:id", h.DeletePermission)
}

func setupUsers(r fiber.Router, reg *registry.Registry) {
	h := reg.Auth.Handler
	r.Get("/users", h.GetAllUsers)
	r.Get("/users/pagination", h.GetAllUsersWithPagination)
	r.Get("/users/:id", h.GetUserByID)
	r.Put("/users/:id", h.UpdateUser)
	r.Delete("/users/:id", h.DeleteUser)
}
