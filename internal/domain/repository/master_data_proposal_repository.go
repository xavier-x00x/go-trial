package repository

import (
	"context"
	"go-trial/internal/domain/entity"
)

type MasterDataProposalRepository interface {
	Create(ctx context.Context, p *entity.MasterDataProposal) error
	FindByID(ctx context.Context, id string) (*entity.MasterDataProposal, error)
	FindByReferenceNumber(ctx context.Context, refNum string) (*entity.MasterDataProposal, error)
	FindByGroupID(ctx context.Context, groupID string) ([]entity.MasterDataProposal, error)
	FindByEntityType(ctx context.Context, entityType string, status string) ([]entity.MasterDataProposal, error)
	FindByProposedByID(ctx context.Context, userID string) ([]entity.MasterDataProposal, error)
	FindPending(ctx context.Context) ([]entity.MasterDataProposal, error)
	FindAll(ctx context.Context) ([]entity.MasterDataProposal, error)
	FindAllWithPaginationGrouped(ctx context.Context, filter entity.QueryFilter) ([]entity.MasterDataProposal, *entity.Meta, error)
	Update(ctx context.Context, p *entity.MasterDataProposal) error
}