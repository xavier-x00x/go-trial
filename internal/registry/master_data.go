package registry

import (
	"go-trial/internal/config"
	"go-trial/internal/delivery/http/handler"
	"go-trial/internal/infrastructure/repository"
	"go-trial/internal/infrastructure/uow"
	"go-trial/internal/query/service"
	"go-trial/internal/usecase"
	"go-trial/pkg/validator"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MasterDataRegistry struct {
	ProposalHandler         *handler.MasterDataProposalHandler
	PlanningHandler       *handler.PurchaseOrderPlanningHandler
	PurchaseOrderHandler *handler.PurchaseOrderHandler
	GoodsReceiptHandler  *handler.GoodsReceiptHandler
}

func NewMasterDataRegistry(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *MasterDataRegistry {
	uow := uow.NewUnitOfWork(db)
	v := validator.New()

	notificationService := usecase.NewNotificationService(
		usecase.NotificationServiceConfig{
			WAEndpoint: cfg.App.WAEndpoint,
			SMTP: usecase.SMTPConfig{
				Host:      cfg.SMTP.Host,
				Port:      cfg.SMTP.Port,
				Username:  cfg.SMTP.Username,
				Password:  cfg.SMTP.Password,
				FromName:  cfg.SMTP.FromName,
				FromEmail: cfg.SMTP.FromEmail,
			},
		},
		repository.NewPurchaseOrderRepository(db),
		repository.NewNotificationQueueRepository(rdb),
	)

	proposalUseCase := usecase.NewMasterDataProposalUseCase(usecase.MasterDataProposalUseCaseConfig{
		Repo:                repository.NewMasterDataProposalRepository(db),
		ProductRepo:         repository.NewProductRepository(db),
		ProductPriceRepo:    repository.NewProductPriceRepository(db),
		ProductUOMRepo:      repository.NewProductUOMConversionRepository(db),
		SupplierRepo:        repository.NewSupplierRepository(db),
		ProductSupplierRepo: repository.NewProductSupplierRepository(db),
		CoaRepo:             repository.NewChartOfAccountRepository(db),
		TaxRepo:             repository.NewTaxRepository(db),
		NumberSequenceRepo:  repository.NewNumberSequenceRepository(db),
		GoodsReceiptRepo:    repository.NewGoodsReceiptRepository(db),
		PriceListRepo:       repository.NewPriceListRepository(db),
		InventoryStockRepo:  repository.NewInventoryStockRepository(db),
		Uow:                 uow,
	})

	planningUseCase := usecase.NewPurchaseOrderPlanningUseCase(usecase.PurchaseOrderPlanningConfig{
		PlanningRepo:                repository.NewPurchaseOrderPlanningRepository(db),
		StoreProductAssortmentRepo: repository.NewStoreProductAssortmentRepository(db),
		StoreRepo:                  repository.NewStoreRepository(db),
		Uow:                       uow,
	})

	poUseCase := usecase.NewPurchaseOrderUseCase(usecase.PurchaseOrderConfig{
		Repo:                repository.NewPurchaseOrderRepository(db),
		PlanningRepo:        repository.NewPurchaseOrderPlanningRepository(db),
		ProductSupplierRepo: repository.NewProductSupplierRepository(db),
		ProductRepo:         repository.NewProductRepository(db),
		SupplierRepo:        repository.NewSupplierRepository(db),
		StoreRepo:           repository.NewStoreRepository(db),
		WarehouseRepo:       repository.NewWarehouseRepository(db),
		UserRepo:            repository.NewUserRepository(db),
		NumberSequenceRepo:  repository.NewNumberSequenceRepository(db),
		NotificationService: notificationService,
		Uow:                 uow,
	})

	grUseCase := usecase.NewGoodsReceiptUseCase(usecase.GoodsReceiptConfig{
		GRRepo:              repository.NewGoodsReceiptRepository(db),
		PurchaseOrderRepo:   repository.NewPurchaseOrderRepository(db),
		InventoryStockRepo: repository.NewInventoryStockRepository(db),
		ProductRepo:        repository.NewProductRepository(db),
		ProductUOMRepo:     repository.NewProductUOMConversionRepository(db),
		UserRepo:            repository.NewUserRepository(db),
		WarehouseRepo:      repository.NewWarehouseRepository(db),
		NumberSequenceRepo: repository.NewNumberSequenceRepository(db),
		Uow:                uow,
	})

	proposalQueryService := service.NewMasterDataProposalQueryService(db)
	poQueryService := service.NewPurchaseOrderQueryService(db)
	poPlanningQueryService := service.NewPurchaseOrderPlanningQueryService(db)

	return &MasterDataRegistry{
		ProposalHandler:         handler.NewMasterDataProposalHandler(proposalUseCase, proposalQueryService, v),
		PlanningHandler:       handler.NewPurchaseOrderPlanningHandler(planningUseCase, poPlanningQueryService, v),
		PurchaseOrderHandler: handler.NewPurchaseOrderHandler(poUseCase, poQueryService, v),
		GoodsReceiptHandler:  handler.NewGoodsReceiptHandler(grUseCase, v),
	}
}
