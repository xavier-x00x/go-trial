package repository

import (
	"context"
	"fmt"

	"go-trial/internal/domain/entity"
	domainRepo "go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"gorm.io/gorm"
)

type masterDataProposalRepository struct {
	db *gorm.DB
}

func NewMasterDataProposalRepository(db *gorm.DB) domainRepo.MasterDataProposalRepository {
	return &masterDataProposalRepository{db: db}
}

func (r *masterDataProposalRepository) Create(ctx context.Context, p *entity.MasterDataProposal) error {
	return uow.GetTx(ctx, r.db).Create(p).Error
}

func (r *masterDataProposalRepository) FindByID(ctx context.Context, id string) (*entity.MasterDataProposal, error) {
	var p entity.MasterDataProposal
	err := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Preload("Items").Where("id = ?", id).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find proposal: %w", err)
	}
	return &p, nil
}

func (r *masterDataProposalRepository) FindByReferenceNumber(ctx context.Context, refNum string) (*entity.MasterDataProposal, error) {
	var p entity.MasterDataProposal
	err := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Preload("Items").Where("reference_number = ?", refNum).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find proposal by reference number: %w", err)
	}
	return &p, nil
}

func (r *masterDataProposalRepository) FindByGroupID(ctx context.Context, groupID string) ([]entity.MasterDataProposal, error) {
	var proposals []entity.MasterDataProposal
	err := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Preload("Items").Where("group_id = ?", groupID).Find(&proposals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find proposals by group: %w", err)
	}
	return proposals, nil
}

func (r *masterDataProposalRepository) FindByEntityType(ctx context.Context, entityType string, status string) ([]entity.MasterDataProposal, error) {
	var proposals []entity.MasterDataProposal
	query := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Where("entity_type = ?", entityType)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&proposals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find proposals: %w", err)
	}
	return proposals, nil
}

func (r *masterDataProposalRepository) FindByProposedByID(ctx context.Context, userID string) ([]entity.MasterDataProposal, error) {
	var proposals []entity.MasterDataProposal
	err := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Where("proposed_by_id = ?", userID).Find(&proposals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find proposals: %w", err)
	}
	return proposals, nil
}

func (r *masterDataProposalRepository) FindPending(ctx context.Context) ([]entity.MasterDataProposal, error) {
	var proposals []entity.MasterDataProposal
	err := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Where("status = ?", entity.ProposalStatusPending).Find(&proposals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find pending proposals: %w", err)
	}
	return proposals, nil
}

func (r *masterDataProposalRepository) FindAll(ctx context.Context) ([]entity.MasterDataProposal, error) {
	var proposals []entity.MasterDataProposal
	err := r.db.WithContext(ctx).Preload("ProposedBy").Preload("ReviewedBy").Find(&proposals).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find proposals: %w", err)
	}
	return proposals, nil
}

func (r *masterDataProposalRepository) FindAllWithPaginationGrouped(ctx context.Context, filter entity.QueryFilter) ([]entity.MasterDataProposal, *entity.Meta, error) {
	page := filter.Page
	limit := filter.Limit

	// Get distinct group_ids with pagination
	var groupIDs []string
	var total, totalFiltered int64

	// Count total distinct groups
	r.db.Model(&entity.MasterDataProposal{}).Distinct("group_id").Count(&total)

	// Apply status filter if present
	query := r.db.Model(&entity.MasterDataProposal{})
	if filter.Conditions != nil {
		if status, ok := filter.Conditions["status"].(string); ok && status != "" {
			query = query.Where("status = ?", status)
		}
		if entityType, ok := filter.Conditions["entity_type"].(string); ok && entityType != "" {
			query = query.Where("entity_type = ?", entityType)
		}
	}
	query.Distinct("group_id").Count(&totalFiltered)

	// Get group IDs with pagination
	subQuery := r.db.Model(&entity.MasterDataProposal{}).
		Select("DISTINCT group_id").
		Order("MAX(created_at) DESC").
		Offset((page - 1) * limit).
		Limit(limit)

	if filter.Conditions != nil {
		if status, ok := filter.Conditions["status"].(string); ok && status != "" {
			subQuery = subQuery.Where("status = ?", status)
		}
		if entityType, ok := filter.Conditions["entity_type"].(string); ok && entityType != "" {
			subQuery = subQuery.Where("entity_type = ?", entityType)
		}
	}

	if err := subQuery.Pluck("group_id", &groupIDs); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch group IDs: %w", err)
	}

	// If no groups, return empty
	if len(groupIDs) == 0 {
		return []entity.MasterDataProposal{}, &entity.Meta{
			Page:    page,
			Limit:   limit,
			Total:   int(total),
		}, nil
	}

	// Get all proposals for these groups, maintaining order
	var proposals []entity.MasterDataProposal
	finalQuery := r.db.WithContext(ctx).
		Preload("ProposedBy").
		Preload("ReviewedBy").
		Where("group_id IN ?", groupIDs).
		Order("created_at DESC")

	if err := finalQuery.Find(&proposals).Error; err != nil {
		return nil, nil, fmt.Errorf("failed to fetch proposals: %w", err)
	}

	meta := &entity.Meta{
		Page:    page,
		Limit:   limit,
		Total:   int(total),
	}

	return proposals, meta, nil
}

func (r *masterDataProposalRepository) Update(ctx context.Context, p *entity.MasterDataProposal) error {
	return uow.GetTx(ctx, r.db).Save(p).Error
}