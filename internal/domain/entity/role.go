package entity

import "github.com/google/uuid"

type Role struct {
	BaseModel
	Name        string       `gorm:"unique" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"primaryKey"`
	PermissionID uuid.UUID `gorm:"primaryKey"`
}
