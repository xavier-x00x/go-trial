package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"go-trial/internal/query/row"
)

type ProductQueryService struct {
	db *gorm.DB
}

func NewProductQueryService(
	db *gorm.DB,
) *ProductQueryService {
	return &ProductQueryService{
		db: db,
	}
}

func (s *ProductQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.ProductListRow, *entity.Meta, error) {

	allowedOrder := []string{"id", "sku", "name", "updated_at"}
	searchColumns := []string{"p.code", "p.name"}

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["p.deleted_at"] = nil

	baseQuery := s.db.WithContext(ctx).
		Table("products p").
		Select(`
            p.id,
            p.sku,
			p.barcode,
            p.name,
            c.name as category_name,
			u.name as uom_name,
			p.updated_at
        `).
		Joins(`
            LEFT JOIN categories c
            ON c.id = p.category_id
        `).
		Joins(`
            LEFT JOIN uom u
            ON u.id = p.base_uom_id
        `)

	return PaginateAndFilter[row.ProductListRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}
