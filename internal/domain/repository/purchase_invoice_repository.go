package repository

import (
	"context"
	"go-trial/internal/domain/entity"

	"github.com/google/uuid"
)

type QueryFilter struct {
	Page         int
	Limit        int
	OrderBy      string
	OrderDir     string
	Search       string
	SearchColumns []string
	Conditions   map[string]interface{}
}

// PurchaseInvoiceRepository defines persistence operations for PurchaseInvoice.
type PurchaseInvoiceRepository interface {
	Create(ctx context.Context, inv *entity.PurchaseInvoice) error
	Update(ctx context.Context, inv *entity.PurchaseInvoice) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteItemsByPurchaseInvoiceID(ctx context.Context, invoiceID string) error

	FindByID(ctx context.Context, id string) (*entity.PurchaseInvoice, error)
	FindByInvoiceNumber(ctx context.Context, invoiceNum string) (*entity.PurchaseInvoice, error)
	FindByStoreID(ctx context.Context, storeID string, status string) ([]entity.PurchaseInvoice, error)
	FindPendingByStoreID(ctx context.Context, storeID string) ([]entity.PurchaseInvoice, error)
	FindAllWithPagination(ctx context.Context, filter *QueryFilter) ([]entity.PurchaseInvoice, *entity.Meta, error)
}
