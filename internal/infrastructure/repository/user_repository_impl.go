package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainRepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return uow.GetTx(ctx, r.db).Create(user).Error
}

func (r *userRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.User, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.User{})
	return PaginateAndFilter[entity.User](r.db, baseQuery, filter)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := uow.GetTx(ctx, r.db).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	err := uow.GetTx(ctx, r.db).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := uow.GetTx(ctx, r.db).Find(&users).Error
	return users, err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := uow.GetTx(ctx, r.db).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByIdentity looks up a user by email or username.
func (r *userRepository) FindByIdentity(ctx context.Context, identity string) (*entity.User, error) {
	var user entity.User
	err := uow.GetTx(ctx, r.db).
		Where("email = ? OR username = ?", identity, identity).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return uow.GetTx(ctx, r.db).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	// suffix username & email with _DEL_<timestamp>
	suffix := fmt.Sprintf("_DEL_%d", time.Now().Unix())
	updates := map[string]interface{}{
		"username": gorm.Expr("CONCAT(username, ?)", suffix),
		"email":    gorm.Expr("CONCAT(email, ?)", suffix),
	}

	if err := uow.GetTx(ctx, r.db).Model(&entity.User{}).
		Where("id = ?", id).UpdateColumns(updates).Error; err != nil {
		return err
	}

	return uow.GetTx(ctx, r.db).Where("id = ?", id).Delete(&entity.User{}).Error
}

func (r *userRepository) FindByPIN(ctx context.Context, pin string) (*entity.User, error) {
	var user entity.User
	err := uow.GetTx(ctx, r.db).Where("pin = ?", pin).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	var user entity.User
	err := uow.GetTx(ctx, r.db).Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetPermissions(ctx context.Context, userID string, role string) ([]string, error) {
	var permissions []string

	if role == "programmer" || role == "administrator" {
		err := uow.GetTx(ctx, r.db).Model(&entity.Permission{}).
			Where("deleted_at IS NULL").
			Pluck("name", &permissions).Error
		if err != nil {
			return nil, err
		}
		if permissions == nil {
			permissions = []string{}
		}
		return permissions, nil
	}

	err := uow.GetTx(ctx, r.db).Table("users").
		Joins("JOIN roles ON users.role = roles.name").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.id = ?", userID).
		Where("users.deleted_at IS NULL AND roles.deleted_at IS NULL AND permissions.deleted_at IS NULL").
		Pluck("permissions.name", &permissions).Error

	if err != nil {
		return nil, err
	}

	if permissions == nil {
		permissions = []string{}
	}
	return permissions, nil
}

