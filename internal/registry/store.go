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

type StoreRegistry struct {
	Handler  *handler.StoreHandler
	Usecase  usecase.StoreUseCase
}

func NewStoreRegistry(db *gorm.DB, cfg *config.Config) *StoreRegistry {
	uow := uow.NewUnitOfWork(db)
	v := validator.New()

	storeRepo := repository.NewStoreRepository(db)
	storeUseCase := usecase.NewStoreUseCase(storeRepo, uow)

	return &StoreRegistry{
		Handler: handler.NewStoreHandler(storeUseCase, v),
		Usecase:  storeUseCase,
	}
}