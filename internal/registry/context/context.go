package context

import (
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}

func NewContext(db *gorm.DB) *Context {
	return &Context{DB: db}
}