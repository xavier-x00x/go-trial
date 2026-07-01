package row

type ProductListRow struct {
	ID            string  `json:"id"`
	SKU           string  `json:"sku"`
	Barcode       *string `json:"barcode"`
	Name          string  `json:"name"`
	Variant       *string `json:"variant"`
	CategoryID    *string `json:"category_id"`
	CategoryName  *string `json:"category_name"`
	BaseUOMID     string  `json:"base_uom_id"`
	UOMName       *string `json:"uom_name"`
	IsStockable   bool    `json:"is_stockable"`
	IsStackable   bool    `json:"is_stackable"`
	IsTaxable     bool    `json:"is_taxable"`
	Length        float64 `json:"length"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	Weight        float64 `json:"weight"`
	MaxStackLayer int     `json:"max_stack_layer"`
	TaxID         *string `json:"tax_id"`
	TaxName       *string `json:"tax_name"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ProductDetailRow struct {
	ProductListRow
}

type ProductSupplierOptionRow struct {
	ID            string  `json:"id" gorm:"column:id"`
	SKU           string  `json:"sku" gorm:"column:sku"`
	Name          string  `json:"name" gorm:"column:name"`
	BaseUOMID     string  `json:"base_uom_id" gorm:"column:base_uom_id"`
	IsContracted  bool    `json:"is_contracted" gorm:"column:is_contracted"`
	OfferedPrice  float64 `json:"offered_price" gorm:"column:offered_price"`
	PurchaseUOMID *string `json:"purchase_uom_id" gorm:"column:purchase_uom_id"`
	MinOrderQty   float64 `json:"min_order_qty" gorm:"column:min_order_qty"`
}
