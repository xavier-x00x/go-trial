package usecase

import (
	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"

	"gorm.io/gorm/utils"
)

func BuildQueryFilter(meta *dto.MetaRequest, allowedOrder []string, searchColumns []string) entity.QueryFilter {
	direction := []string{"asc", "desc"}

	if !utils.Contains(allowedOrder, meta.OrderColumn) || !utils.Contains(direction, meta.OrderDir) {
		meta.OrderColumn = "id"
		meta.OrderDir = "asc"
	}

	return entity.QueryFilter{
		Page:         meta.Page,
		Limit:        meta.Limit,
		Search:       meta.Search,
		OrderColumn:  meta.OrderColumn,
		OrderDir:     meta.OrderDir,
		SearchColumn: searchColumns,
		Conditions:   map[string]interface{}{},
	}
}
