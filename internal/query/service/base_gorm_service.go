package service

import (
	"fmt"
	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

func PaginateAndFilter[T any](db *gorm.DB, baseQuery *gorm.DB, meta *params.MetaRequest, allowedOrder []string, searchColumns []string) ([]T, *entity.Meta, error) {
	direction := []string{"asc", "desc"}

	// fungsi : validasi meta request
	if !utils.Contains(allowedOrder, meta.OrderColumn) || !utils.Contains(direction, meta.OrderDir) {
		meta.OrderColumn = allowedOrder[0] // Default to first allowed column if invalid
		meta.OrderDir = "asc"
	}

	page := meta.Page
	if page < 1 {
		page = 1
	}

	limit := meta.Limit
	if limit < 1 {
		limit = 10 // default limit
	}

	search := meta.Search
	orderBy := fmt.Sprintf("%s %s", meta.OrderColumn, strings.ToUpper(meta.OrderDir))

	var dataList []T
	var total, totalFiltered int64

	// Clone the base query for filtered results
	query := baseQuery.Session(&gorm.Session{})

	// Apply conditions
	if meta.Conditions != nil {
		for key, val := range meta.Conditions {
			if strVal, ok := val.(string); ok && strVal == "" {
				continue
			}
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

		for _, column := range searchColumns {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
			values = append(values, "%"+search+"%")
		}

		if len(conditions) > 0 {
			query = query.Where("("+strings.Join(conditions, " OR ")+")", values...)
		}
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
	if meta.Conditions != nil {
		for key, val := range meta.Conditions {
			if strVal, ok := val.(string); ok && strVal == "" {
				continue
			}
			if val == nil {
				totalQuery = totalQuery.Where(fmt.Sprintf("%s IS NULL", key))
			} else {
				totalQuery = totalQuery.Where(fmt.Sprintf("%s = ?", key), val)
			}
		}
	}
	totalQuery.Count(&total)

	metaResp := &entity.Meta{
		Page:          page,
		Limit:         limit,
		Total:         int(total),
		TotalFiltered: int(totalFiltered),
		LastPage:      int(math.Ceil(float64(totalFiltered) / float64(limit))),
		Draw:          len(dataList),
	}

	return dataList, metaResp, nil
}
