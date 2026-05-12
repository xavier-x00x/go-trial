package registry

import (
	"go-trial/internal/config"
	jwtPkg "go-trial/pkg/jwt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Registry struct {
	Auth        *AuthRegistry
	Store       *StoreRegistry
	Product     *ProductRegistry
	Supplier    *SupplierRegistry
	Finance     *FinanceRegistry
	Warehouse  *WarehouseRegistry
	MasterData *MasterDataRegistry
	JWTManager *jwtPkg.JWTManager
}

func New(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *Registry {
	return &Registry{
		Auth:       NewAuthRegistry(db, cfg),
		Store:      NewStoreRegistry(db, cfg),
		Product:    NewProductRegistry(db, cfg),
		Supplier:   NewSupplierRegistry(db, cfg),
		Finance:    NewFinanceRegistry(db, cfg),
		Warehouse:  NewWarehouseRegistry(db, cfg),
		MasterData: NewMasterDataRegistry(db, rdb, cfg),
		JWTManager: jwtPkg.NewJWTManager(cfg.JWT.Secret),
	}
}