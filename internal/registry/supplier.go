package registry

import (
	"go-trial/internal/config"
	"go-trial/internal/delivery/http/handler"
	"go-trial/internal/infrastructure/repository"
	"go-trial/internal/infrastructure/uow"
	"go-trial/internal/query/service"
	"go-trial/internal/usecase"
	"go-trial/pkg/validator"

	"gorm.io/gorm"
)

type SupplierRegistry struct {
	Handler         *handler.SupplierHandler
	CategoryHandler *handler.SupplierCategoryHandler
}

func NewSupplierRegistry(db *gorm.DB, cfg *config.Config) *SupplierRegistry {
	uow := uow.NewUnitOfWork(db)
	v := validator.New()

	supplierRepo := repository.NewSupplierRepository(db)
	coaRepo := repository.NewChartOfAccountRepository(db)
	supplierCategoryRepo := repository.NewSupplierCategoryRepository(db)

	supplierUseCase := usecase.NewSupplierUseCase(supplierRepo, coaRepo, uow)
	supplierCategoryUseCase := usecase.NewSupplierCategoryUseCase(supplierCategoryRepo, uow)
	supplierCategoryQueryService := service.NewSupplierCategoryQueryService(db)

	return &SupplierRegistry{
		Handler:         handler.NewSupplierHandler(supplierUseCase, v),
		CategoryHandler: handler.NewSupplierCategoryHandler(supplierCategoryUseCase, supplierCategoryQueryService, v),
	}
}