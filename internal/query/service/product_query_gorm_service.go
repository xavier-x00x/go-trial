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

func (s *ProductQueryService) getBaseSelectAndJoins(ctx context.Context) *gorm.DB {
	return s.db.WithContext(ctx).
		Table("products p").
		Select(`
			p.id,
			p.sku,
			p.barcode,
			p.name,
			p.variant,
			p.category_id,
			c.name as category_name,
			p.base_uom_id,
			u.name as uom_name,
			p.is_stockable,
			p.is_stackable,
			p.is_taxable,
			p.length,
			p.width,
			p.height,
			p.weight,
			p.max_stack_layer,
			p.tax_id,
			t.name as tax_name,
			p.created_at,
			p.updated_at
		`).
		Joins(`
			LEFT JOIN product_categories c
			ON c.id = p.category_id
		`).
		Joins(`
			LEFT JOIN uom u
			ON u.id = p.base_uom_id
		`).
		Joins(`
			LEFT JOIN taxes t
			ON t.id = p.tax_id
		`)
}

func (s *ProductQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.ProductListRow, *entity.Meta, error) {

	allowedOrder := []string{"id", "sku", "name", "variant", "updated_at"}
	searchColumns := []string{"p.sku", "p.name", "p.barcode", "p.variant"}

	// Map order column to use table alias to prevent ambiguous column error
	switch param.OrderColumn {
	case "id", "sku", "name", "variant", "updated_at":
		param.OrderColumn = "p." + param.OrderColumn
	}
	allowedOrder = append(allowedOrder, "p.id", "p.sku", "p.name", "p.variant", "p.updated_at")

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["p.deleted_at"] = nil

	baseQuery := s.getBaseSelectAndJoins(ctx)

	return PaginateAndFilter[row.ProductListRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}

func (s *ProductQueryService) GetAll(ctx context.Context) ([]row.ProductListRow, error) {
	var rows []row.ProductListRow
	err := s.getBaseSelectAndJoins(ctx).
		Where("p.deleted_at IS NULL").
		Order("p.updated_at DESC").
		Find(&rows).Error
	return rows, err
}

func (s *ProductQueryService) GetByID(ctx context.Context, id string) (*row.ProductDetailRow, error) {
	var detail row.ProductDetailRow
	err := s.getBaseSelectAndJoins(ctx).
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&detail).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}

func (s *ProductQueryService) GetProductSuppliers(ctx context.Context, productID string) ([]entity.ProductSupplier, error) {
	var rows []entity.ProductSupplier
	err := s.db.WithContext(ctx).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Preload("Supplier").
		Preload("Store").
		Find(&rows).Error
	return rows, err
}

func (s *ProductQueryService) GetProductsBySupplier(ctx context.Context, supplierID string) ([]row.ProductSupplierOptionRow, error) {
	var rows []row.ProductSupplierOptionRow
	err := s.db.WithContext(ctx).
		Table("products p").
		Select(`
			p.id,
			p.sku,
			p.name,
			p.base_uom_id,
			(CASE WHEN ps.id IS NOT NULL THEN true ELSE false END) as is_contracted,
			COALESCE(ps.offered_price, 0) as offered_price,
			ps.purchase_uom_id,
			COALESCE(ps.min_order_qty, 0) as min_order_qty
		`).
		Joins("LEFT JOIN product_suppliers ps ON ps.product_id = p.id AND ps.supplier_id = ? AND ps.deleted_at IS NULL", supplierID).
		Where("p.deleted_at IS NULL").
		Order("is_contracted DESC, p.name ASC").
		Find(&rows).Error
	return rows, err
}

func (s *ProductQueryService) GetProductsBySupplierWithPagination(
	ctx context.Context,
	supplierID string,
	param *params.MetaRequest,
) ([]row.ProductSupplierOptionRow, *entity.Meta, error) {

	allowedOrder := []string{"is_contracted", "name", "sku"}
	searchColumns := []string{"p.sku", "p.name", "p.barcode"}

	if param.OrderColumn == "" {
		param.OrderColumn = "is_contracted"
		param.OrderDir = "desc"
	}

	baseQuery := s.db.WithContext(ctx).
		Table("products p").
		Select(`
			p.id,
			p.sku,
			p.name,
			p.base_uom_id,
			(CASE WHEN ps.id IS NOT NULL THEN true ELSE false END) as is_contracted,
			COALESCE(ps.offered_price, 0) as offered_price,
			ps.purchase_uom_id,
			COALESCE(ps.min_order_qty, 0) as min_order_qty
		`).
		Joins("LEFT JOIN product_suppliers ps ON ps.product_id = p.id AND ps.supplier_id = ? AND ps.deleted_at IS NULL", supplierID).
		Where("p.deleted_at IS NULL")

	return PaginateAndFilter[row.ProductSupplierOptionRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}
