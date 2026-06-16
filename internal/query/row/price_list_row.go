package row

import "time"

type PriceListRow struct {
	ID           string     `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	CurrencyCode string     `json:"currency_code"`
	StoreID      *string    `json:"store_id,omitempty"`
	StoreName    *string    `json:"store_name,omitempty"`
	IsActive     bool       `json:"is_active"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
