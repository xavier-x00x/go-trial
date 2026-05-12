package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	FindByID(ctx context.Context, id string) (*entity.Role, error)
	FindByName(ctx context.Context, name string) (*entity.Role, error)
	FindAll(ctx context.Context) ([]entity.Role, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Role, *entity.Meta, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
	ReplacePermissions(ctx context.Context, roleID string, permissions []entity.Permission) error
}

type PermissionRepository interface {
	Create(ctx context.Context, perm *entity.Permission) error
	FindByID(ctx context.Context, id string) (*entity.Permission, error)
	FindByIDs(ctx context.Context, ids []string) ([]entity.Permission, error)
	FindByPath(ctx context.Context, path string) (*entity.Permission, error)
	FindAll(ctx context.Context) ([]entity.Permission, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Permission, *entity.Meta, error)
	Update(ctx context.Context, perm *entity.Permission) error
	Delete(ctx context.Context, id string) error
}