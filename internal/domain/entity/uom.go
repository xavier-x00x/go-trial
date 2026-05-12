package entity

// UOM (Unit of Measure) merepresentasikan satuan ukuran dasar
type UOM struct {
	BaseModel
	Code string `gorm:"type:varchar(10);not null" json:"code"`
	Name string `gorm:"type:varchar(50);not null" json:"name"`
}

func (UOM) TableName() string {
	return "uom"
}