package registry

import (
	"go-trial/internal/config"
	"go-trial/internal/delivery/http/handler"
	"go-trial/internal/infrastructure/repository"
	"go-trial/internal/usecase"
	"go-trial/pkg/validator"

	"gorm.io/gorm"
)

type WarehouseRegistry struct {
	Handler *handler.WarehouseHandler
}

func NewWarehouseRegistry(db *gorm.DB, cfg *config.Config) *WarehouseRegistry {
	v := validator.New()

	warehouseRepo := repository.NewWarehouseRepository(db)
	warehouseUseCase := usecase.NewWarehouseUsecase(warehouseRepo)

	return &WarehouseRegistry{
		Handler: handler.NewWarehouseHandler(warehouseUseCase, v),
	}
}