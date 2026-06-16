package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"go-trial/internal/query/row"
)

type SupplierCategoryQueryService struct {
	db *gorm.DB
}

func NewSupplierCategoryQueryService(db *gorm.DB) *SupplierCategoryQueryService {
	return &SupplierCategoryQueryService{db: db}
}

func (s *SupplierCategoryQueryService) GetByID(ctx context.Context, id string) (*entity.SupplierCategory, error) {
	var cat entity.SupplierCategory
	err := s.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

func (s *SupplierCategoryQueryService) GetAll(ctx context.Context) ([]entity.SupplierCategory, error) {
	var cats []entity.SupplierCategory
	err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&cats).Error
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *SupplierCategoryQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.SupplierCategoryRow, *entity.Meta, error) {
	allowedOrder := []string{"id", "name", "updated_at"}
	searchColumns := []string{"sc.name", "sc.description"}

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["sc.deleted_at"] = nil

	baseQuery := s.db.WithContext(ctx).
		Table("supplier_categories sc").
		Select(`
            sc.id,
            sc.name,
            sc.description,
            sc.updated_at
        `)

	return PaginateAndFilter[row.SupplierCategoryRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}
