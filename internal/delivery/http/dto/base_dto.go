package dto

type MetaRequest struct {
	Page        int                    `json:"page"`
	Limit       int                    `json:"limit"`
	Search      string                `json:"search"`
	OrderColumn string               `json:"order_column"`
	OrderDir    string              `json:"order_dir"`
	Conditions map[string]interface{} `json:"conditions,omitempty"`
}

type MetaResponse struct {
	Page          int `json:"page"`
	Limit         int `json:"limit"`
	Total         int `json:"total"`
	TotalFiltered int `json:"total_filtered"`
	LastPage      int `json:"last_page"`
	Draw          int `json:"draw"`
}
