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

type FinanceRegistry struct {
	Handler                *handler.COAHandler
	CustomerHandler        *handler.CustomerHandler
	PaymentMethodHandler   *handler.PaymentMethodHandler
	PriceListHandler       *handler.PriceListHandler
	TaxHandler             *handler.TaxHandler
	PurchaseInvoiceHandler *handler.PurchaseInvoiceHandler
	PurchasePaymentHandler *handler.PurchasePaymentHandler
	PurchaseReturnHandler  *handler.PurchaseReturnHandler
	ExpenseVoucherHandler  *handler.ExpenseVoucherHandler
}

func NewFinanceRegistry(db *gorm.DB, cfg *config.Config) *FinanceRegistry {
	uow := uow.NewUnitOfWork(db)
	v := validator.New()

	coaRepo := repository.NewChartOfAccountRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	priceListRepo := repository.NewPriceListRepository(db)
	taxRepo := repository.NewTaxRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	userRepo := repository.NewUserRepository(db)
	numberSequenceRepo := repository.NewNumberSequenceRepository(db)
	monthlyAPBalanceRepo := repository.NewMonthlyAPBalanceRepository(db)
	purchaseInvoiceRepo := repository.NewPurchaseInvoiceRepository(db)
	purchasePaymentRepo := repository.NewPurchasePaymentRepository(db)
	purchaseReturnRepo := repository.NewPurchaseReturnRepository(db)
	inventoryStockRepo := repository.NewInventoryStockRepository(db)
	expenseVoucherRepo := repository.NewExpenseVoucherRepository(db)

	coaUseCase := usecase.NewChartOfAccountUseCase(coaRepo, uow)
	customerUseCase := usecase.NewCustomerUseCase(customerRepo, coaRepo, uow)
	paymentMethodUseCase := usecase.NewPaymentMethodUsecase(paymentMethodRepo)
	priceListUseCase := usecase.NewPriceListUsecase(priceListRepo)
	taxUseCase := usecase.NewTaxUsecase(taxRepo)
	purchaseInvoiceUseCase := usecase.NewPurchaseInvoiceUseCase(usecase.PurchaseInvoiceConfig{
		Repo:                 purchaseInvoiceRepo,
		PurchaseOrderRepo:    repository.NewPurchaseOrderRepository(db),
		SupplierRepo:         supplierRepo,
		StoreRepo:            repository.NewStoreRepository(db),
		WarehouseRepo:        repository.NewWarehouseRepository(db),
		UserRepo:             userRepo,
		NumberSequenceRepo:   numberSequenceRepo,
		ChartOfAccountRepo:   coaRepo,
		MonthlyAPBalanceRepo: monthlyAPBalanceRepo,
		DB:                   db,
		Uow:                  uow,
	})
	purchasePaymentUseCase := usecase.NewPurchasePaymentUseCase(usecase.PurchasePaymentConfig{
		Repo:                    purchasePaymentRepo,
		PurchaseInvoiceRepo:    purchaseInvoiceRepo,
		PurchaseReturnRepo:     purchaseReturnRepo,
		SupplierRepo:            supplierRepo,
		UserRepo:                userRepo,
		ChartOfAccountRepo:      coaRepo,
		MonthlyAPBalanceRepo:    monthlyAPBalanceRepo,
		NumberSequenceRepo:     numberSequenceRepo,
		DB:                      db,
		Uow:                     uow,
	})
	purchaseReturnUseCase := usecase.NewPurchaseReturnUseCase(usecase.PurchaseReturnConfig{
		PRRepo:               purchaseReturnRepo,
		PIRepo:               purchaseInvoiceRepo,
		UserRepo:             userRepo,
		StoreRepo:            repository.NewStoreRepository(db),
		InventoryStockRepo:   inventoryStockRepo,
		MonthlyAPBalanceRepo: monthlyAPBalanceRepo,
		ChartOfAccountRepo:   coaRepo,
		NumberSequenceRepo:   numberSequenceRepo,
		DB:                   db,
		Uow:                  uow,
	})
	expenseVoucherUseCase := usecase.NewExpenseVoucherUseCase(
		expenseVoucherRepo,
		coaRepo,
		userRepo,
		numberSequenceRepo,
		uow,
		db,
	)

	coaQueryService := service.NewCOAQueryService(db)

	return &FinanceRegistry{
		Handler:                handler.NewCOAHandler(coaUseCase, coaQueryService, v),
		CustomerHandler:       handler.NewCustomerHandler(customerUseCase, v),
		PaymentMethodHandler:   handler.NewPaymentMethodHandler(paymentMethodUseCase, v),
		PriceListHandler:      handler.NewPriceListHandler(priceListUseCase, v),
		TaxHandler:            handler.NewTaxHandler(taxUseCase, v),
		PurchaseInvoiceHandler: handler.NewPurchaseInvoiceHandler(purchaseInvoiceUseCase),
		PurchasePaymentHandler: handler.NewPurchasePaymentHandler(purchasePaymentUseCase),
		PurchaseReturnHandler:  handler.NewPurchaseReturnHandler(purchaseReturnUseCase),
		ExpenseVoucherHandler:  handler.NewExpenseVoucherHandler(expenseVoucherUseCase, v),
	}
}