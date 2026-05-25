package usecase

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"

	"gorm.io/gorm/utils"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrors struct {
	Errors []FieldError `json:"errors"`
}

func (e *FieldErrors) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Message
	}
	return "Validation failed"
}

func (e *FieldErrors) Add(field, msg string) {
	e.Errors = append(e.Errors, FieldError{Field: field, Message: msg})
}

func BuildQueryFilter(meta *dto.MetaRequest, allowedOrder []string, searchColumns []string) entity.QueryFilter {
	direction := []string{"asc", "desc"}

	if !utils.Contains(allowedOrder, meta.OrderColumn) || !utils.Contains(direction, meta.OrderDir) {
		meta.OrderColumn = "id"
		meta.OrderDir = "asc"
	}

	conditions := map[string]interface{}{}
	for key, val := range meta.Conditions {
		// Skip empty string values agar tidak jadi WHERE kolom = ''
		if strVal, ok := val.(string); ok && strVal == "" {
			continue
		}
		conditions[key] = val
	}

	return entity.QueryFilter{
		Page:         meta.Page,
		Limit:        meta.Limit,
		Search:       meta.Search,
		OrderColumn:  meta.OrderColumn,
		OrderDir:     meta.OrderDir,
		SearchColumn: searchColumns,
		Conditions:   conditions,
	}
}
