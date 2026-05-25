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

func (r *masterDataProposalRepository) FindAllWithPagination(ctx context.Context, filter entity.QueryFilter) ([]entity.MasterDataProposal, *entity.Meta, error) {
	baseQuery := r.db.WithContext(ctx).Model(&entity.MasterDataProposal{}).
		Preload("ProposedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Preload("ReviewedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		})
	return PaginateAndFilter[entity.MasterDataProposal](r.db, baseQuery, filter)
}

func (r *masterDataProposalRepository) Update(ctx context.Context, p *entity.MasterDataProposal) error {
	return uow.GetTx(ctx, r.db).Save(p).Error
}

func (r *masterDataProposalRepository) DeleteItemsByProposalID(ctx context.Context, proposalID string) error {
	return uow.GetTx(ctx, r.db).Where("proposal_id = ?", proposalID).Delete(&entity.MasterDataProposalItem{}).Error
}
