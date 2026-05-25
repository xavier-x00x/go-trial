package repository

import (
	"fmt"
	"go-trial/internal/domain/entity"
	"math"
	"strings"

	"gorm.io/gorm"
)

// PaginateAndFilter provides generic pagination, filtering, and search
// for any GORM model type. It accepts an optional baseQuery to allow
// callers to add pre-conditions (e.g. store_code filtering).
//
// Usage:
//
//	query := r.db.Model(&entity.Category{}).Where("store_code = ?", code)
//	return PaginateAndFilter[entity.Category](r.db, query, filter)
func PaginateAndFilter[T any](db *gorm.DB, baseQuery *gorm.DB, filter entity.QueryFilter) ([]T, *entity.Meta, error) {
	page := filter.Page
	limit := filter.Limit
	search := filter.Search
	orderBy := fmt.Sprintf("%s %s", filter.OrderColumn, strings.ToUpper(filter.OrderDir))

	var dataList []T
	var total, totalFiltered int64

	// Clone the base query for filtered results
	query := baseQuery.Session(&gorm.Session{})

	// Apply conditions
	if filter.Conditions != nil {
		for key, val := range filter.Conditions {
			if val == nil {
				query = query.Where(fmt.Sprintf("%s IS NULL", key))
			} else {
				query = query.Where(fmt.Sprintf("%s = ?", key), val)
			}
		}
	}

	// Apply search
	if search != "" {
		var conditions []string
		var values []interface{}

		for _, column := range filter.SearchColumn {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
			values = append(values, "%"+search+"%")
		}

		query = query.Where("("+strings.Join(conditions, " OR ")+")", values...)
	}

	// Count filtered results
	query.Count(&totalFiltered)

	// Fetch paginated data
	query.Offset((page - 1) * limit).Limit(limit).Order(orderBy).Find(&dataList)

	if len(dataList) == 0 {
		dataList = []T{}
	}

	// Count total (without search, only base conditions)
	totalQuery := baseQuery.Session(&gorm.Session{})
	if filter.Conditions != nil {
		for key, val := range filter.Conditions {
			if val == nil {
				totalQuery = totalQuery.Where(fmt.Sprintf("%s IS NULL", key))
			} else {
				totalQuery = totalQuery.Where(fmt.Sprintf("%s = ?", key), val)
			}
		}
	}
	totalQuery.Count(&total)

	meta := &entity.Meta{
		Page:          page,
		Limit:         limit,
		Total:         int(total),
		TotalFiltered: int(totalFiltered),
		LastPage:      int(math.Ceil(float64(totalFiltered) / float64(limit))),
		Draw:          len(dataList),
	}

	return dataList, meta, nil
}
