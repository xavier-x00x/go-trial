package usecase

import (
	"context"
	"errors"
	"log"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
)

var (
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleNameExists = errors.New("role name already exists")
	ErrPermNotFound   = errors.New("permission not found")
	ErrPermPathExists = errors.New("permission path already exists")
)

type RolePermissionUseCase interface {
	CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error)
	GetRole(ctx context.Context, id string) (*dto.RoleResponse, error)
	GetRoleByName(ctx context.Context, name string) (*dto.RoleResponse, error)
	ListRoles(ctx context.Context) ([]dto.RoleResponse, error)
	ListRolesWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.RoleResponse, *entity.Meta, error)
	UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id string) error

	CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error)
	GetPermission(ctx context.Context, id string) (*dto.PermissionResponse, error)
	GetPermissionByPath(ctx context.Context, path string) (*dto.PermissionResponse, error)
	ListPermissions(ctx context.Context) ([]dto.PermissionResponse, error)
	ListPermissionsWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PermissionResponse, *entity.Meta, error)
	UpdatePermission(ctx context.Context, id string, req dto.UpdatePermissionRequest) (*dto.PermissionResponse, error)
	DeletePermission(ctx context.Context, id string) error
	SyncPermissions(ctx context.Context, items []dto.SyncPermissionItem) ([]dto.PermissionResponse, error)
}

type rolePermissionUseCase struct {
	rolePermRepo   repository.RoleRepository
	permissionRepo repository.PermissionRepository
	uow            uow.UnitOfWork
}

func NewRolePermissionUseCase(
	rolePermRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	uow uow.UnitOfWork,
) RolePermissionUseCase {
	return &rolePermissionUseCase{
		rolePermRepo:   rolePermRepo,
		permissionRepo: permissionRepo,
		uow:            uow,
	}
}

func (u *rolePermissionUseCase) CreateRole(ctx context.Context, req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	existing, _ := u.rolePermRepo.FindByName(ctx, req.Name)
	if existing != nil {
		return nil, ErrRoleNameExists
	}

	id, _ := uuid.NewV7()
	role := &entity.Role{
		BaseModel:   entity.BaseModel{ID: id},
		Name:        req.Name,
		Permissions: nil,
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.rolePermRepo.Create(txCtx, role); err != nil {
		return nil, err
	}

	if len(req.Permissions) > 0 {
		perms, _ := u.permissionRepo.FindByPaths(ctx, req.Permissions)
		_ = u.rolePermRepo.ReplacePermissions(txCtx, role.ID.String(), perms)
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	loaded, _ := u.rolePermRepo.FindByID(ctx, role.ID.String())
	return toRoleResponse(loaded), nil
}

func (u *rolePermissionUseCase) GetRole(ctx context.Context, id string) (*dto.RoleResponse, error) {
	role, err := u.rolePermRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return toRoleResponse(role), nil
}

func (u *rolePermissionUseCase) GetRoleByName(ctx context.Context, name string) (*dto.RoleResponse, error) {
	role, err := u.rolePermRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return toRoleResponse(role), nil
}

func (u *rolePermissionUseCase) ListRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	roles, err := u.rolePermRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.RoleResponse, len(roles))
	for i, r := range roles {
		resp[i] = *toRoleResponse(&r)
	}
	return resp, nil
}

func (u *rolePermissionUseCase) ListRolesWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.RoleResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "created_at"}
	searchColumns := []string{"id", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.rolePermRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.RoleResponse{}
	for _, r := range data {
		resp = append(resp, *toRoleResponse(&r))
	}
	return resp, resMeta, nil
}

func (u *rolePermissionUseCase) UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := u.rolePermRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	// log print for debugging
	log.Println("Replacing permissions for role %s: %v", id, req.Permissions)

	if req.Name != nil && *req.Name != role.Name {
		existing, _ := u.rolePermRepo.FindByName(ctx, *req.Name)
		if existing != nil {
			return nil, ErrRoleNameExists
		}
		role.Name = *req.Name
	}

	if req.Permissions != nil && len(req.Permissions) > 0 {
		perms, _ := u.permissionRepo.FindByPaths(ctx, req.Permissions)
		_ = u.rolePermRepo.ReplacePermissions(txCtx, id, perms)
	}

	role.Permissions = nil
	if err := u.rolePermRepo.Update(txCtx, role); err != nil {
		return nil, err
	}

	if req.Permissions == nil {
		if err := u.uow.Commit(txCtx); err != nil {
			return nil, err
		}
	} else if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	loaded, _ := u.rolePermRepo.FindByID(ctx, id)
	return toRoleResponse(loaded), nil
}

func (u *rolePermissionUseCase) DeleteRole(ctx context.Context, id string) error {
	role, err := u.rolePermRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}
	return u.rolePermRepo.Delete(ctx, id)
}

func (u *rolePermissionUseCase) CreatePermission(ctx context.Context, req dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	existing, _ := u.permissionRepo.FindByPath(ctx, req.Path)
	if existing != nil {
		return nil, ErrPermPathExists
	}

	id, _ := uuid.NewV7()
	perm := &entity.Permission{
		BaseModel: entity.BaseModel{ID: id},
		Path:      req.Path,
		Name:      req.Name,
	}

	if err := u.permissionRepo.Create(ctx, perm); err != nil {
		return nil, err
	}

	return toPermissionResponse(perm), nil
}

func (u *rolePermissionUseCase) GetPermission(ctx context.Context, id string) (*dto.PermissionResponse, error) {
	perm, err := u.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, ErrPermNotFound
	}
	return toPermissionResponse(perm), nil
}

func (u *rolePermissionUseCase) GetPermissionByPath(ctx context.Context, path string) (*dto.PermissionResponse, error) {
	perm, err := u.permissionRepo.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, ErrPermNotFound
	}
	return toPermissionResponse(perm), nil
}

func (u *rolePermissionUseCase) ListPermissions(ctx context.Context) ([]dto.PermissionResponse, error) {
	perms, err := u.permissionRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.PermissionResponse, len(perms))
	for i, p := range perms {
		resp[i] = *toPermissionResponse(&p)
	}
	return resp, nil
}

func (u *rolePermissionUseCase) ListPermissionsWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.PermissionResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "path", "created_at"}
	searchColumns := []string{"id", "name", "path"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)

	data, resMeta, err := u.permissionRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.PermissionResponse{}
	for _, p := range data {
		resp = append(resp, *toPermissionResponse(&p))
	}
	return resp, resMeta, nil
}

func (u *rolePermissionUseCase) UpdatePermission(ctx context.Context, id string, req dto.UpdatePermissionRequest) (*dto.PermissionResponse, error) {
	perm, err := u.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if perm == nil {
		return nil, ErrPermNotFound
	}

	if req.Path != nil {
		existing, _ := u.permissionRepo.FindByPath(ctx, *req.Path)
		if existing != nil {
			return nil, ErrPermPathExists
		}
		perm.Path = *req.Path
	}
	if req.Name != nil {
		perm.Name = *req.Name
	}

	if err := u.permissionRepo.Update(ctx, perm); err != nil {
		return nil, err
	}

	return toPermissionResponse(perm), nil
}

func (u *rolePermissionUseCase) SyncPermissions(ctx context.Context, items []dto.SyncPermissionItem) ([]dto.PermissionResponse, error) {
	paths := make([]string, len(items))
	for i, p := range items {
		paths[i] = p.Path
	}

	existing, err := u.permissionRepo.FindByPaths(ctx, paths)
	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool, len(existing))
	for _, p := range existing {
		existingMap[p.Path] = true
	}

	var created []dto.PermissionResponse

	for _, item := range items {
		if existingMap[item.Path] {
			continue
		}

		id, _ := uuid.NewV7()
		perm := &entity.Permission{
			BaseModel: entity.BaseModel{ID: id},
			Path:      item.Path,
			Name:      item.Name,
		}

		if err := u.permissionRepo.Create(ctx, perm); err != nil {
			return nil, err
		}

		created = append(created, *toPermissionResponse(perm))
	}

	return created, nil
}

func (u *rolePermissionUseCase) DeletePermission(ctx context.Context, id string) error {
	perm, err := u.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if perm == nil {
		return ErrPermNotFound
	}
	return u.permissionRepo.Delete(ctx, id)
}

func toRoleResponse(role *entity.Role) *dto.RoleResponse {
	if role == nil {
		return nil
	}

	perms := make([]string, len(role.Permissions))
	for i, p := range role.Permissions {
		perms[i] = p.Path
	}

	return &dto.RoleResponse{
		ID:         role.ID.String(),
		Name:       role.Name,
		Permission: perms,
		CreatedAt:  role.CreatedAt,
		UpdatedAt:  role.UpdatedAt,
	}
}

func toPermissionResponse(perm *entity.Permission) *dto.PermissionResponse {
	if perm == nil {
		return nil
	}
	return &dto.PermissionResponse{
		ID:   perm.ID.String(),
		Path: perm.Path,
		Name: perm.Name,
	}
}
