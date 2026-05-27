package row

import "time"

type PriceListRow struct {
	ID           string     `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	CurrencyCode string     `json:"currency_code"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsActive     bool       `json:"is_active"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
