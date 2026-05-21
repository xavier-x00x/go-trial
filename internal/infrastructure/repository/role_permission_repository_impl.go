package repository

import (
	"context"
	"errors"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domainRepo.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	return uow.GetTx(ctx, r.db).Create(role).Error
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*entity.Role, error) {
	var role entity.Role
	err := uow.GetTx(ctx, r.db).Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := uow.GetTx(ctx, r.db).Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	err := uow.GetTx(ctx, r.db).Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Role, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Role{})
	return PaginateAndFilter[entity.Role](r.db, baseQuery, filter)
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	return uow.GetTx(ctx, r.db).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	return uow.GetTx(ctx, r.db).Where("id = ?", id).Delete(&entity.Role{}).Error
}

func (r *roleRepository) ReplacePermissions(ctx context.Context, roleID string, permissions []entity.Permission) error {
	tx := uow.GetTx(ctx, r.db)

	if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	for _, p := range permissions {
		if err := tx.Exec(
			"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
			roleID, p.ID.String(),
		).Error; err != nil {
			return err
		}
	}
	return nil
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domainRepo.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, perm *entity.Permission) error {
	return uow.GetTx(ctx, r.db).Create(perm).Error
}

func (r *permissionRepository) FindByID(ctx context.Context, id string) (*entity.Permission, error) {
	var perm entity.Permission
	err := uow.GetTx(ctx, r.db).Where("id = ?", id).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) FindByIDs(ctx context.Context, ids []string) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := uow.GetTx(ctx, r.db).Where("id IN ?", ids).Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindByPath(ctx context.Context, path string) (*entity.Permission, error) {
	var perm entity.Permission
	err := uow.GetTx(ctx, r.db).Where("path = ?", path).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepository) FindByPaths(ctx context.Context, paths []string) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := uow.GetTx(ctx, r.db).Where("path IN ?", paths).Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindAll(ctx context.Context) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := uow.GetTx(ctx, r.db).Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.Permission, *entity.Meta, error) {
	baseQuery := r.db.Model(&entity.Permission{})
	return PaginateAndFilter[entity.Permission](r.db, baseQuery, filter)
}

func (r *permissionRepository) Update(ctx context.Context, perm *entity.Permission) error {
	return uow.GetTx(ctx, r.db).Save(perm).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	return uow.GetTx(ctx, r.db).Where("id = ?", id).Delete(&entity.Permission{}).Error
}