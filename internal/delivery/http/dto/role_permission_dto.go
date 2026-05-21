package dto

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Permissions []string `json:"permissions,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type CreatePermissionRequest struct {
	Path string `json:"path" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdatePermissionRequest struct {
	Path *string `json:"path,omitempty" validate:"omitempty"`
	Name *string `json:"name,omitempty" validate:"omitempty"`
}

type SyncPermissionItem struct {
	Path string `json:"path" validate:"required"`
	Name string `json:"name" validate:"required"`
}