package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"go-trial/internal/query/row"
)

type PriceListQueryService struct {
	db *gorm.DB
}

func NewPriceListQueryService(
	db *gorm.DB,
) *PriceListQueryService {
	return &PriceListQueryService{
		db: db,
	}
}

func (s *PriceListQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.PriceListRow, *entity.Meta, error) {

	allowedOrder := []string{"id", "code", "name", "start_date", "updated_at"}
	searchColumns := []string{"p.code", "p.name"}

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["p.deleted_at"] = nil

	baseQuery := s.db.WithContext(ctx).
		Table("price_lists p").
		Select(`
            p.id,
            p.code,
            p.name,
            p.currency_code,
            p.start_date,
            p.end_date,
            p.is_active,
            p.updated_at
        `)

	return PaginateAndFilter[row.PriceListRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}
