package service

import (
	"context"

	"gorm.io/gorm"

	"go-trial/internal/domain/entity"
	"go-trial/internal/query/params"
	"go-trial/internal/query/row"
)

type MasterDataProposalQueryService struct {
	db *gorm.DB
}

func NewMasterDataProposalQueryService(db *gorm.DB) *MasterDataProposalQueryService {
	return &MasterDataProposalQueryService{db: db}
}

func (s *MasterDataProposalQueryService) getBaseSelectAndJoins(ctx context.Context) *gorm.DB {
	return s.db.WithContext(ctx).
		Table("master_data_proposals m").
		Select(`
			m.id,
			m.reference_number,
			m.entity_type,
			m.action_type,
			m.total_items,
			m.status,
			m.proposed_by_id,
			u1.name as proposed_by_name,
			m.reviewed_by_id,
			u2.name as reviewed_by_name,
			m.reviewed_at,
			m.reason,
			m.created_at
		`).
		Joins("LEFT JOIN users u1 ON u1.id = m.proposed_by_id").
		Joins("LEFT JOIN users u2 ON u2.id = m.reviewed_by_id")
}

func (s *MasterDataProposalQueryService) GetListPagination(
	ctx context.Context,
	param *params.MetaRequest,
) ([]row.MasterDataProposalListRow, *entity.Meta, error) {

	allowedOrder := []string{"id", "reference_number", "created_at", "updated_at", "status", "entity_type"}
	searchColumns := []string{"m.reference_number", "m.reason", "u1.name", "u2.name"}

	if param.Conditions == nil {
		param.Conditions = map[string]interface{}{}
	}
	param.Conditions["m.deleted_at"] = nil

	baseQuery := s.getBaseSelectAndJoins(ctx)

	return PaginateAndFilter[row.MasterDataProposalListRow](s.db, baseQuery, param, allowedOrder, searchColumns)
}

func (s *MasterDataProposalQueryService) GetAll(ctx context.Context) ([]row.MasterDataProposalListRow, error) {
	var rows []row.MasterDataProposalListRow
	err := s.getBaseSelectAndJoins(ctx).
		Where("m.deleted_at IS NULL").
		Order("m.created_at DESC").
		Find(&rows).Error
	return rows, err
}

func (s *MasterDataProposalQueryService) GetPending(ctx context.Context) ([]row.MasterDataProposalListRow, error) {
	var rows []row.MasterDataProposalListRow
	err := s.getBaseSelectAndJoins(ctx).
		Where("m.deleted_at IS NULL").
		Where("m.status = ?", entity.ProposalStatusPending).
		Order("m.created_at DESC").
		Find(&rows).Error
	return rows, err
}

func (s *MasterDataProposalQueryService) GetByEntityType(ctx context.Context, entityType, status string) ([]row.MasterDataProposalListRow, error) {
	var rows []row.MasterDataProposalListRow
	query := s.getBaseSelectAndJoins(ctx).
		Where("m.deleted_at IS NULL").
		Where("m.entity_type = ?", entityType)

	if status != "" {
		query = query.Where("m.status = ?", status)
	}

	err := query.Order("m.created_at DESC").Find(&rows).Error
	return rows, err
}

func (s *MasterDataProposalQueryService) GetByID(ctx context.Context, id string) (*row.MasterDataProposalDetailRow, error) {
	var detail row.MasterDataProposalDetailRow

	err := s.db.WithContext(ctx).
		Table("master_data_proposals m").
		Select(`
			m.id,
			m.reference_number,
			m.entity_type,
			m.action_type,
			m.total_items,
			m.status,
			m.proposed_by_id,
			u1.name as proposed_by_name,
			m.reviewed_by_id,
			u2.name as reviewed_by_name,
			m.reviewed_at,
			m.reason,
			m.review_notes,
			m.created_at,
			m.updated_at
		`).
		Joins("LEFT JOIN users u1 ON u1.id = m.proposed_by_id").
		Joins("LEFT JOIN users u2 ON u2.id = m.reviewed_by_id").
		Where("m.id = ? AND m.deleted_at IS NULL", id).
		First(&detail).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or an error based on your preference
		}
		return nil, err
	}

	// Fetch items
	var items []row.MasterDataProposalItemRow
	err = s.db.WithContext(ctx).
		Table("master_data_proposal_items").
		Where("proposal_id = ? AND deleted_at IS NULL", id).
		Order("seq_no ASC").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	detail.Items = items
	return &detail, nil
}

func (s *MasterDataProposalQueryService) GetByGroup(ctx context.Context, groupID string) ([]row.MasterDataProposalDetailRow, error) {
	// If group_id exists in the table, we query it. 
	// We'll fetch headers and then items for each.
	var proposals []row.MasterDataProposalDetailRow
	err := s.db.WithContext(ctx).
		Table("master_data_proposals m").
		Select(`
			m.id,
			m.reference_number,
			m.entity_type,
			m.action_type,
			m.total_items,
			m.status,
			m.proposed_by_id,
			u1.name as proposed_by_name,
			m.reviewed_by_id,
			u2.name as reviewed_by_name,
			m.reviewed_at,
			m.reason,
			m.review_notes,
			m.created_at,
			m.updated_at
		`).
		Joins("LEFT JOIN users u1 ON u1.id = m.proposed_by_id").
		Joins("LEFT JOIN users u2 ON u2.id = m.reviewed_by_id").
		Where("m.group_id = ? AND m.deleted_at IS NULL", groupID).
		Find(&proposals).Error

	if err != nil {
		return nil, err
	}

	for i := range proposals {
		var items []row.MasterDataProposalItemRow
		err = s.db.WithContext(ctx).
			Table("master_data_proposal_items").
			Where("proposal_id = ? AND deleted_at IS NULL", proposals[i].ID).
			Order("seq_no ASC").
			Find(&items).Error
		if err == nil {
			proposals[i].Items = items
		}
	}

	return proposals, nil
}
