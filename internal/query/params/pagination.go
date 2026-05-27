package params

type MetaRequest struct {
	Page        int                    `json:"page"`
	Limit       int                    `json:"limit"`
	Search      string                 `json:"search"`
	OrderColumn string                 `json:"order_column"`
	OrderDir    string                 `json:"order_dir"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
}
