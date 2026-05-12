package registry

import (
	"go-trial/internal/config"
	"go-trial/internal/delivery/http/handler"
	"go-trial/internal/infrastructure/repository"
	"go-trial/internal/infrastructure/uow"
	"go-trial/internal/usecase"
	"go-trial/pkg/validator"

	"gorm.io/gorm"
)

type SupplierRegistry struct {
	Handler *handler.SupplierHandler
}

func NewSupplierRegistry(db *gorm.DB, cfg *config.Config) *SupplierRegistry {
	uow := uow.NewUnitOfWork(db)
	v := validator.New()

	supplierRepo := repository.NewSupplierRepository(db)
	coaRepo := repository.NewChartOfAccountRepository(db)

	supplierUseCase := usecase.NewSupplierUseCase(supplierRepo, coaRepo, uow)

	return &SupplierRegistry{
		Handler: handler.NewSupplierHandler(supplierUseCase, v),
	}
}