package entity

// Store entity
type Store struct {
	BaseModel
	Code         string  `gorm:"type:varchar(20);not null" json:"code"`
	Name         string  `gorm:"type:varchar(150);not null" json:"name"`
	TaxRegNumber *string `gorm:"type:varchar(50)" json:"tax_reg_number,omitempty"` // NPWP untuk faktur pajak
	Address      *string `gorm:"type:text" json:"address,omitempty"`
	City         *string `gorm:"type:varchar(100)" json:"city,omitempty"`
	Province     *string `gorm:"type:varchar(100)" json:"province,omitempty"`
	PostalCode   *string `gorm:"type:varchar(10)" json:"postal_code,omitempty"`
	Phone        *string `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email        *string `gorm:"type:varchar(100)" json:"email,omitempty"`
	IsMain       bool    `gorm:"type:tinyint(1);not null;default:0" json:"is_main"`
	IsActive     bool    `gorm:"type:tinyint(1);not null;default:1" json:"is_active"`
}

func (Store) TableName() string {
	return "store"
}
