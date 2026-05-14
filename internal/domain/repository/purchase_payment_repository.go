package repository

import (
	"context"
	"go-trial/internal/domain/entity"

	"github.com/google/uuid"
)

type PurchasePaymentRepository interface {
	Create(ctx context.Context, pp *entity.PurchasePayment) error
	Update(ctx context.Context, pp *entity.PurchasePayment) error
	Delete(ctx context.Context, id uuid.UUID) error

	FindByID(ctx context.Context, id string) (*entity.PurchasePayment, error)
	FindByPaymentNumber(ctx context.Context, paymentNum string) (*entity.PurchasePayment, error)
	FindBySupplierID(ctx context.Context, supplierID string) ([]entity.PurchasePayment, error)
	FindAllWithPagination(ctx context.Context, filter *QueryFilter) ([]entity.PurchasePayment, *entity.Meta, error)
}