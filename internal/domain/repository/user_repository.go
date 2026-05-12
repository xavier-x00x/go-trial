package repository

import (
	"context"

	"go-trial/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByIdentity(ctx context.Context, identity string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindAll(ctx context.Context) ([]entity.User, error)
	FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.User, *entity.Meta, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
