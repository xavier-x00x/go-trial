package route

import (
	"go-trial/internal/registry"

	"github.com/gofiber/fiber/v2"
)

func setupOperationalRoutes(r fiber.Router, reg *registry.Registry) {
	setupStoreRoutes(r, reg)
	setupProductRoutes(r, reg)
	setupSupplierRoutes(r, reg)
	setupWarehouseRoutes(r, reg)
	setupMasterDataRoutes(r, reg)
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

	r.Post("", check("products:create"), h.Create)
	r.Get("", check("products:view"), h.GetAll)
	r.Get("/pagination", check("products:view"), h.GetAllWithPagination)
	r.Get("/:id", check("products:view"), h.GetByID)
	r.Put("/:id", check("products:update"), h.Update)
	r.Delete("/:id", check("products:delete"), h.Delete)
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
