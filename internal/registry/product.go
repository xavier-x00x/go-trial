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

type ProductRegistry struct {
	Handler              *handler.ProductHandler
	CategoryHandler      *handler.ProductCategoryHandler
	UOMHandler           *handler.UOMHandler
	ProductPriceHandler  *handler.ProductPriceHandler
	UOMConversionHandler *handler.ProductUOMConversionHandler
}

func NewProductRegistry(db *gorm.DB, cfg *config.Config) *ProductRegistry {
	uow := uow.NewUnitOfWork(db)
	v := validator.New()

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewProductCategoryRepository(db)
	uomRepo := repository.NewUOMRepository(db)
	productPriceRepo := repository.NewProductPriceRepository(db)
	productUOMRepo := repository.NewProductUOMConversionRepository(db)

	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, uomRepo, uow)
	categoryUseCase := usecase.NewProductCategoryUseCase(categoryRepo, uow)
	uomUseCase := usecase.NewUOMUseCase(uomRepo, uow)
	productPriceUseCase := usecase.NewProductPriceUsecase(productPriceRepo)
	productUOMUseCase := usecase.NewProductUOMConversionUsecase(productUOMRepo)

	//Query
	productQueryService := service.NewProductQueryService(db)
	productPriceQueryService := service.NewProductPriceQueryService(db)

	return &ProductRegistry{
		Handler:              handler.NewProductHandler(productUseCase, productQueryService, v),
		CategoryHandler:      handler.NewProductCategoryHandler(categoryUseCase, v),
		UOMHandler:           handler.NewUOMHandler(uomUseCase, v),
		ProductPriceHandler:  handler.NewProductPriceHandler(productPriceUseCase, productPriceQueryService, v),
		UOMConversionHandler: handler.NewProductUOMConversionHandler(productUOMUseCase, v),
	}
}
