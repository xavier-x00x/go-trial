package registry

import (
	"go-trial/internal/config"
	"go-trial/internal/delivery/http/handler"
	"go-trial/internal/infrastructure/repository"
	"go-trial/internal/infrastructure/uow"
	"go-trial/internal/usecase"
	"go-trial/pkg/jwt"
	"go-trial/pkg/validator"

	"gorm.io/gorm"
)

type AuthRegistry struct {
	Handler              *handler.AuthHandler
	RolePermissionHandler *handler.RolePermissionHandler
	UseCase            usecase.AuthUseCase
	RolePermissionUC   usecase.RolePermissionUseCase
}

func NewAuthRegistry(db *gorm.DB, cfg *config.Config) *AuthRegistry {
	v := validator.New()
	uow := uow.NewUnitOfWork(db)

	userRepo := repository.NewUserRepository(db)
	authUseCase := usecase.NewAuthUseCase(userRepo, uow, jwtManager(cfg), cfg)

	rolePermRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	rolePermUseCase := usecase.NewRolePermissionUseCase(rolePermRepo, permissionRepo, uow)

	return &AuthRegistry{
		Handler:              handler.NewAuthHandler(authUseCase, v),
		RolePermissionHandler: handler.NewRolePermissionHandler(rolePermUseCase, v),
		UseCase:              authUseCase,
		RolePermissionUC:     rolePermUseCase,
	}
}

func jwtManager(cfg *config.Config) *jwt.JWTManager {
	return jwt.NewJWTManager(cfg.JWT.Secret)
}