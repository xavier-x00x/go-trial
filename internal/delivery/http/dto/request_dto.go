package dto

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	StoreID  *string `json:"store_id,omitempty" validate:"omitempty,len=36"`
	Name     string  `json:"name" validate:"required,min=2,max=255"`
	Username string  `json:"username" validate:"required,min=3,max=100"`
	Email    string  `json:"email" validate:"required,email"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	Password string  `json:"password" validate:"required,min=8,max=72"`
	Role     string  `json:"role" validate:"omitempty,oneof=admin manager staff"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
	Remember *bool  `json:"remember,omitempty"`
}

// GoogleTokenLoginRequest is the payload for login via Google Token (Access or ID Token).
type GoogleTokenLoginRequest struct {
	Token     string `json:"token" validate:"required"`
	TokenType string `json:"token_type" validate:"required,oneof=access id"`
}

// CreateCOARequest is the payload for creating a chart of account.
type CreateCOARequest struct {
	AccountCode   string `json:"account_code" validate:"required,min=1,max=20"`
	Name         string  `json:"name" validate:"required,min=1,max=100"`
	AccountType  string  `json:"account_type" validate:"required,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	NormalBalance string `json:"normal_balance" validate:"required,oneof=DEBIT CREDIT"`
}

// UpdateCOARequest is the payload for updating a chart of account.
type UpdateCOARequest struct {
	AccountCode   *string `json:"account_code,omitempty" validate:"omitempty,min=1,max=20"`
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	AccountType  *string  `json:"account_type,omitempty" validate:"omitempty,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	NormalBalance *string `json:"normal_balance,omitempty" validate:"omitempty,oneof=DEBIT CREDIT"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

