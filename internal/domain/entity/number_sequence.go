package entity

import (
	"gorm.io/gorm"
)

type NumberSequence struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Prefix     string `gorm:"size:10;not null;uniqueIndex:uq_prefix_period"`
	Period     string `gorm:"size:4;not null;uniqueIndex:uq_prefix_period"` // format: YYMM
	LastNumber int    `gorm:"default:0"`
}

func (n *NumberSequence) TableName() string {
	return "number_sequences"
}

func (n *NumberSequence) BeforeCreate(tx *gorm.DB) error {
	return nil
}