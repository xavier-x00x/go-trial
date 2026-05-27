package row

type ProductListRow struct {
	ID           string  `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	CategoryName string  `json:"category_name"`
	SellPrice    float64 `json:"sell_price"`
}
