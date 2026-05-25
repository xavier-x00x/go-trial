package route

import (
	"go-trial/internal/registry"

	"github.com/gofiber/fiber/v2"
)

func setupFinanceRoutes(r fiber.Router, reg *registry.Registry) {
	f := reg.Finance
	setupAccounts(r, f)
	setupCustomers(r, f)
	setupPaymentMethods(r, f)
	setupPriceLists(r, f)
	setupTaxes(r, f)
	setupPurchaseInvoices(r, f)
	setupPurchasePayments(r, f)
	setupPurchaseReturns(r, f)
	setupExpenseVouchers(r, f)
}

func setupAccounts(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.Handler
	r = r.Group("/accounts")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/pagination", h.GetAllWithPagination)
	r.Post("/import", h.Import)
	r.Get("/import/template", h.DownloadTemplate)
	r.Get("/tree", h.GetTree)
	r.Get("/:id", h.GetByID)
	r.Get("/type/:type", h.GetByType)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
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

func setupPurchasePayments(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.PurchasePaymentHandler
	r = r.Group("/purchase-payments")
	r.Post("", h.Create)
	r.Get("/pagination", h.List)
	r.Get("/:id", h.GetByID)
	r.Post("/:id/post", h.Post)
	r.Post("/:id/void", h.Void)
}

func setupPurchaseReturns(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.PurchaseReturnHandler
	r = r.Group("/purchase-returns")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Post("/:id/post", h.Post)
}

func setupExpenseVouchers(r fiber.Router, f *registry.FinanceRegistry) {
	h := f.ExpenseVoucherHandler
	r = r.Group("/expense-vouchers")
	r.Post("", h.Create)
	r.Get("", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Post("/:id/post", h.Post)
	r.Post("/:id/cancel", h.Cancel)
}
