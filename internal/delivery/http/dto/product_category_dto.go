package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductCategory DTOs
type CreateCategoryRequest struct {
	ParentID         *uuid.UUID      `json:"parent_id,omitempty" validate:"omitempty"`
	Name             string          `json:"name" validate:"required,min=1,max=100"`
	Slug             string          `json:"slug" validate:"required,min=1,max=120"`
	Code             string          `json:"code" validate:"required,len=3"`
	Sequence         int             `json:"sequence"`
	DefaultMarkupPct decimal.Decimal `json:"default_markup_pct"`
}

type UpdateCategoryRequest struct {
	ParentID         *uuid.UUID       `json:"parent_id,omitempty" validate:"omitempty"`
	Name             *string          `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Slug             *string          `json:"slug,omitempty" validate:"omitempty,min=1,max=120"`
	Code             *string          `json:"code,omitempty" validate:"omitempty,len=3"`
	Sequence         *int             `json:"sequence,omitempty"`
	DefaultMarkupPct *decimal.Decimal `json:"default_markup_pct"`
}

type CategoryResponse struct {
	ID               string            `json:"id"`
	ParentID         *uuid.UUID        `json:"parent_id,omitempty"`
	Parent           *CategoryResponse `json:"parent,omitempty"`
	Name             string            `json:"name"`
	Slug             string            `json:"slug"`
	Code             string            `json:"code"`
	Sequence         int               `json:"sequence"`
	DefaultMarkupPct decimal.Decimal   `json:"default_markup_pct"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}
