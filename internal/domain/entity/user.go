package entity

import "time"

type User struct {
	ID           string     `gorm:"type:char(36);primaryKey" json:"id"`
	StoreID      *string    `json:"store_id,omitempty" gorm:"type:char(36);index;comment:FK → store.id (null = semua toko)"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Username     string     `gorm:"type:varchar(100);uniqueIndex:idx_username_deleted;not null" json:"username"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone        *string    `gorm:"type:varchar(20);uniqueIndex;comment:Nomor telepon;null OK" json:"phone,omitempty"`
	Password     string     `gorm:"type:varchar(255);not null" json:"-"`
	Role         string     `gorm:"type:varchar(50);not null;default:''" json:"role"`
	AuthProvider string     `gorm:"type:varchar(50);not null;default:'local'" json:"auth_provider"`
	GoogleID     *string    `gorm:"type:varchar(255);uniqueIndex;null" json:"google_id,omitempty"`
	PIN          *string    `gorm:"type:varchar(20);index" json:"pin,omitempty"`
	AvatarURL    *string    `gorm:"type:varchar(255);null" json:"avatar_url,omitempty"`
	IsActive     *bool      `gorm:"type:tinyint(1);not null;default:true" json:"is_active"`
	LastLoginAt  *time.Time `gorm:"type:datetime;comment:Terakhir login" json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"uniqueIndex:idx_username_deleted" json:"-"`
	Store        *Store     `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}
