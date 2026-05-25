package dto

import "time"

// ChartOfAccountResponse is the public representation of a chart of account.
type ChartOfAccountResponse struct {
	ID            string  `json:"id"`
	AccountCode   string  `json:"account_code"`
	Name          string  `json:"name"`
	AccountType   string  `json:"account_type"`
	NormalBalance string  `json:"normal_balance"`
	IsActive      bool    `json:"is_active"`
	ParentID      *string `json:"parent_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// ChartOfAccountTreeResponse is the tree representation with children.
type ChartOfAccountTreeResponse struct {
	ID            string                       `json:"id"`
	AccountCode   string                       `json:"account_code"`
	Name          string                       `json:"name"`
	AccountType   string                       `json:"account_type"`
	NormalBalance string                       `json:"normal_balance"`
	IsActive      bool                         `json:"is_active"`
	Children      []ChartOfAccountTreeResponse `json:"children,omitempty"`
}

// AuthResponse is returned after login or token refresh.
type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// UserResponse is the public representation of a user.
type UserResponse struct {
	ID          string     `json:"id"`
	StoreID     *string    `json:"store_id,omitempty"`
	StoreName   *string    `json:"store_name,omitempty"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone,omitempty"`
	Role        string     `json:"role"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	IsActive    *bool      `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Permissions []string   `json:"permissions"`
}

// RefreshResponse is returned after refreshing the access token.
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// MessageResponse is a simple message-only response.
type MessageResponse struct {
	Message string `json:"message"`
}

type RoleResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Permission []string  `json:"permissions,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PermissionResponse struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
