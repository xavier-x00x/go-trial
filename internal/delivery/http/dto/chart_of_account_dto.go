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

type AccountImportRow struct {
	RowNumber     int    `json:"row_number"`
	AccountCode   string `json:"account_code"`
	Name          string `json:"name"`
	AccountType   string `json:"account_type"`
	NormalBalance string `json:"normal_balance"`
}

type AccountImportResult struct {
	TotalRows   int            `json:"total_rows"`
	SuccessRows int            `json:"success_rows"`
	ErrorRows   int            `json:"error_rows"`
	Errors      []ImportRowError `json:"errors,omitempty"`
}

type ImportRowError struct {
	RowNumber int    `json:"row_number"`
	Message  string `json:"message"`
}