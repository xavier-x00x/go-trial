package route

import (
	"go-trial/internal/infrastructure/middleware"
	"go-trial/internal/registry"
	jwtPkg "go-trial/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, reg *registry.Registry, jwtManager *jwtPkg.JWTManager) {
	// Public routes (no auth required)
	api := app.Group("/api")
	api.Post("/auth/register", reg.Auth.Handler.Register)
	api.Post("/auth/login", reg.Auth.Handler.Login)
	api.Post("/auth/refresh", reg.Auth.Handler.RefreshToken)
	api.Post("/auth/logout", reg.Auth.Handler.Logout)

	// Protected routes (auth required)
	protected := api.Group("", middleware.AuthMiddleware(jwtManager))
	protected.Get("/auth/me", reg.Auth.Handler.Me)

	setupRoles(protected, reg)
	setupPermissions(protected, reg)
	setupStoreRoutes(protected, reg)
	setupProductRoutes(protected, reg)
	setupSupplierRoutes(protected, reg)
	setupFinanceRoutes(protected, reg)
	setupWarehouseRoutes(protected, reg)
	setupMasterDataRoutes(protected, reg)
}

func setupStoreRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.Store.Handler
	r.Post("/stores", h.Create)
	r.Get("/stores", h.GetAll)
	r.Get("/stores/pagination", h.GetAllWithPagination)
	r.Get("/stores/:id", h.GetByID)
	r.Put("/stores/:id", h.Update)
	r.Delete("/stores/:id", h.Delete)
}

func setupProductRoutes(r fiber.Router, reg *registry.Registry) {
	p := reg.Product
	setupProducts(r, p)
	setupCategories(r, p)
	setupUOMs(r, p)
	setupProductPrices(r, p)
	setupProductUOMs(r, p)
}

func setupProducts(r fiber.Router, p *registry.ProductRegistry) {
	h := p.Handler
	r = r.Group("/products")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupCategories(r fiber.Router, p *registry.ProductRegistry) {
	h := p.CategoryHandler
	r = r.Group("/categories")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupUOMs(r fiber.Router, p *registry.ProductRegistry) {
	h := p.UOMHandler
	r = r.Group("/uoms")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
}

func setupProductPrices(r fiber.Router, p *registry.ProductRegistry) {
	h := p.ProductPriceHandler
	r = r.Group("/product-prices")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Get("/product/:productId", h.GetByProductID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupProductUOMs(r fiber.Router, p *registry.ProductRegistry) {
	h := p.UOMConversionHandler
	r = r.Group("/product-uoms")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Get("/product/:productId", h.GetByProductID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupSupplierRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.Supplier.Handler
	r.Post("/suppliers", h.Create)
	r.Get("/suppliers", h.GetAll)
	r.Get("/suppliers/pagination", h.GetAllWithPagination)
	r.Get("/suppliers/:id", h.GetByID)
	r.Put("/suppliers/:id", h.Update)
	r.Delete("/suppliers/:id", h.Delete)
}

func setupFinanceRoutes(r fiber.Router, reg *registry.Registry) {
	f := reg.Finance
	setupAccounts(r, f)
	setupCustomers(r, f)
	setupPaymentMethods(r, f)
	setupPriceLists(r, f)
	setupTaxes(r, f)
	setupPurchaseInvoices(r, f)
	setupPurchasePayments(r, f)
}

func setupPurchaseInvoices(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.PurchaseInvoiceHandler
	r = r.Group("/purchase-invoices")
	r.Post("", h.Create)
	r.Get("", h.List)
	r.Get("/:id", h.GetByID)
	r.Get("/invoice-number/:invoice_number", h.GetByInvoiceNumber)
	r.Get("/store/:storeId", h.GetByStoreID)
	r.Put("/:id", h.Update)
	r.Post("/:id/submit", h.Submit)
	r.Post("/:id/approve", h.Approve)
	r.Post("/:id/verify", h.Verify)
	r.Post("/:id/post", h.Post)
	r.Post("/:id/cancel", h.Cancel)
}

func setupAccounts(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.Handler
	r = r.Group("/accounts")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Get("/type/:type", h.GetByType)
	r.Put("/:id", h.Update)
}

func setupCustomers(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.CustomerHandler
	r = r.Group("/customers")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupPaymentMethods(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.PaymentMethodHandler
	r = r.Group("/payment-methods")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupPriceLists(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.PriceListHandler
	r = r.Group("/price-lists")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/active", h.GetActive)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupTaxes(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.TaxHandler
	r = r.Group("/taxes")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func setupPurchasePayments(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.PurchasePaymentHandler
	r = r.Group("/purchase-payments")
	r.Post("", h.Create)
	r.Get("/pagination", h.List)
	r.Get("/:id", h.GetByID)
	r.Post("/:id/post", h.Post)
	r.Post("/:id/void", h.Void)
}

func setupWarehouseRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.Warehouse.Handler
	r.Post("/warehouses", h.Create)
	r.Get("/warehouses", h.GetAll)
	r.Get("/warehouses/pagination", h.GetAllWithPagination)
	r.Get("/warehouses/:id", h.GetByID)
	r.Put("/warehouses/:id", h.Update)
	r.Delete("/warehouses/:id", h.Delete)
}

func setupMasterDataRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.MasterData.ProposalHandler
	r.Get("/master-data", h.GetAll)
	r.Post("/master-data", h.Create)
	r.Get("/master-data/pagination", h.GetAllWithPagination)
	r.Get("/master-data/pending", h.GetPending)
	r.Get("/master-data/entity/:entityType", h.GetByEntityType)
	r.Get("/master-data/group/:groupId", h.GetByGroup)
	r.Get("/master-data/:id", h.GetByID)
	r.Put("/master-data/:id", h.Update)
	r.Post("/master-data/:id/review", h.Review)
	r.Post("/master-data/:id/execute", h.Execute)
	r.Post("/master-data/bulk/product-supplier", h.BulkLinkProductSupplier)

	setupPlanningRoutes(r, reg)

	setupPurchaseOrderRoutes(r, reg)

	setupGoodsReceiptRoutes(r, reg)
}

func setupPurchaseOrderRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.MasterData.PurchaseOrderHandler
	r = r.Group("/purchase-orders")
	r.Post("", h.Create)
	r.Post("/from-planning", h.CreateFromPlanning)
	r.Post("/bulk-from-approved", h.BulkCreateFromApprovedPlanning)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/store/:storeId", h.GetByStoreID)
	r.Get("/store/:storeId/pending", h.GetPendingByStoreID)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Post("/:id/submit", h.Submit)
	r.Post("/:id/approve", h.Approve)
	r.Post("/cancel", h.Cancel)
}

func setupGoodsReceiptRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.MasterData.GoodsReceiptHandler
	r = r.Group("/goods-receipts")
	r.Post("", h.Create)
	r.Post("/with-invoice", h.CreateWithInvoice)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Get("/po/:poId", h.GetByPurchaseOrderID)
	r.Get("/warehouse/:warehouseId", h.GetByWarehouseID)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Post("/:id/confirm", h.Confirm)
	r.Post("/:id/cancel", h.Cancel)
}

func setupPlanningRoutes(r fiber.Router, reg *registry.Registry) {
	h := reg.MasterData.PlanningHandler
	r = r.Group("/planning")
	r.Post("/calculate", h.Calculate)
	r.Get("/pending", h.GetPending)
	r.Get("", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Post("/approve", h.Approve)
	r.Post("/ignore", h.Ignore)
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