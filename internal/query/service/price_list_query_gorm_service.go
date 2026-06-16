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

	allowedOrder := []string{"id", "code", "name", "updated_at"}
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
            p.is_active,
            p.store_id,
            s.name as store_name,
            p.updated_at
        `).
		Joins("LEFT JOIN store s ON s.id = p.store_id")

	return PaginateAndFilter[row.PriceListRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}

func (s *PriceListQueryService) GetByID(ctx context.Context, id string) (*entity.PriceList, error) {
	var pl entity.PriceList
	err := s.db.WithContext(ctx).Preload("Store").Where("id = ? AND deleted_at IS NULL", id).First(&pl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pl, nil
}

func (s *PriceListQueryService) GetAll(ctx context.Context) ([]entity.PriceList, error) {
	var pls []entity.PriceList
	err := s.db.WithContext(ctx).Preload("Store").Where("deleted_at IS NULL").Find(&pls).Error
	if err != nil {
		return nil, err
	}
	return pls, nil
}

func (s *PriceListQueryService) GetActive(ctx context.Context) ([]entity.PriceList, error) {
	var pls []entity.PriceList
	err := s.db.WithContext(ctx).Preload("Store").Where("is_active = ? AND deleted_at IS NULL", true).Find(&pls).Error
	if err != nil {
		return nil, err
	}
	return pls, nil
}
