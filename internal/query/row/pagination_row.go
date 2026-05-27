package row

type MetaResponse struct {
	Page          int `json:"page"`
	Limit         int `json:"limit"`
	Total         int `json:"total"`
	TotalFiltered int `json:"total_filtered"`
	LastPage      int `json:"last_page"`
	Draw          int `json:"draw"`
}

type PaginationResult struct {
	Meta MetaResponse `json:"meta"`
	Data interface{}  `json:"data"`
}
