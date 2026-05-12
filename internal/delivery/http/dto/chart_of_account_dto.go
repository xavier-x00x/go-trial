package dto

type CreateChartOfAccountRequest struct {
	AccountCode   string `json:"account_code" validate:"required,max=20"`
	Name         string `json:"name" validate:"required,max=100"`
	AccountType  string `json:"account_type" validate:"required,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	NormalBalance string `json:"normal_balance" validate:"required,oneof=DEBIT CREDIT"`
}

type UpdateChartOfAccountRequest struct {
	AccountCode   *string `json:"account_code,omitempty" validate:"omitempty,max=20"`
	Name         *string `json:"name,omitempty" validate:"omitempty,max=100"`
	AccountType  *string `json:"account_type,omitempty" validate:"omitempty,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	NormalBalance *string `json:"normal_balance,omitempty" validate:"omitempty,oneof=DEBIT CREDIT"`
	IsActive     *bool   `json:"is_active,omitempty"`
}