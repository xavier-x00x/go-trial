package database

import (
	"fmt"
	"log"
	"time"

	"go-trial/internal/config"
	"go-trial/internal/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQL(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := autoMigrateWithCleanup(db); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	log.Println("Database connected and migrated")
	return db
}

func autoMigrateWithCleanup(db *gorm.DB) error {
	cleanupColumns := map[string][]string{
		"price_lists": {"start_date", "end_date"},
		"product_prices": {"discount_pct"},
		"users": {"phone"},
		"purchase_orders": {
			"subtotal", "discount_amount", "tax_amount", "freight_amount",
			"other_cost_amount", "grand_total", "is_tax_inclusive",
		},
		"purchase_order_items": {
			"discount_1_pct", "discount_2_pct", "discount_3_pct", "discount_amount",
			"total_discount_amount", "tax_pct", "tax_amount", "net_unit_price",
			"landed_cost_amount",
		},
	}

	for table, cols := range cleanupColumns {
		for _, col := range cols {
			if db.Migrator().HasColumn(table, col) {
				if err := db.Exec("ALTER TABLE " + table + " DROP COLUMN " + col).Error; err != nil {
					log.Printf("Note: could not drop column %s.%s: %v", table, col, err)
				} else {
					log.Printf("Cleaned up old column %s.%s", table, col)
				}
			}
		}
	}

	return db.AutoMigrate(
		&entity.User{},
		&entity.Store{},
		&entity.UOM{},
		&entity.ProductCategory{},
		&entity.Product{},
		&entity.SupplierCategory{},
		&entity.Supplier{},
		&entity.Customer{},
		&entity.ChartOfAccount{},
		&entity.Warehouse{},
		&entity.InventoryStock{},
		&entity.MonthlyInventoryStock{},
		&entity.Tax{},
		&entity.PaymentMethod{},
		&entity.PriceList{},
		&entity.ProductPrice{},
		&entity.MonthlyAPBalance{},
		&entity.ProductSupplier{},
		&entity.ProductUOMConversion{},
		&entity.StoreProductAssortment{},
		&entity.MasterDataProposal{},
		&entity.MasterDataProposalItem{},
		&entity.Role{},
		&entity.Permission{},
		&entity.NumberSequence{},
		&entity.PurchaseOrderPlanning{},
		&entity.PurchaseOrder{},
		&entity.PurchaseOrderItem{},
		&entity.GoodsReceipt{},
		&entity.GoodsReceiptItem{},
	)
}
