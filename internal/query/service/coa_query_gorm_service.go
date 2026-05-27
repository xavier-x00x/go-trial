package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"go-trial/internal/query/row"
)

type COAQueryService struct {
	db *gorm.DB
}

func NewCOAQueryService(
	db *gorm.DB,
) *COAQueryService {
	return &COAQueryService{
		db: db,
	}
}

func (s *COAQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.COARow, *entity.Meta, error) {

	allowedOrder := []string{"id", "account_code", "name", "updated_at"}
	searchColumns := []string{"c.account_code", "c.name"}

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["c.deleted_at"] = nil

	baseQuery := s.db.WithContext(ctx).
		Table("chart_of_accounts c").
		Select(`
            c.id,
            c.account_code as code,
            c.name,
			c.updated_at
        `)

	return PaginateAndFilter[row.COARow](s.db, baseQuery, param, allowedOrder, searchColumns)
}
